package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type ProductFromSeller struct {
	ProductId int     `json:"product_id"`
	SellerId  int     `json:"seller_id"`
	Quantity  int     `json:"quantity"  binding:"required"`
	Cost      float64 `json:"cost"  binding:"required"`
	Published bool    `json:"published"`
}

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func scanProductFromSeller(row pgx.Row) (ProductFromSeller, error) {
	var product ProductFromSeller
	err := row.Scan(&product.ProductId, &product.SellerId, &product.Quantity, &product.Cost, &product.Published)
	return product, err
}

func scanProduct(row pgx.Row) (Product, error) {
	var product Product
	err := row.Scan(&product.Id, &product.Name, &product.Description)
	return product, err
}

func ProductSellers(product_id int) ([]ProductFromSeller, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
	SELECT ps.seller_id, ps.quantity, ps.cost, ps.published
	FROM product_seller ps
		JOIN users u ON ps.seller_id = u.id
	WHERE ps.product_id = $1 and published = true`, product_id)
	if err != nil {
		return nil, err
	}
	return scanManyData(rows, scanProductFromSeller)
}

func ListProducts() ([]Product, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	return scanManyData(rows, scanProduct)
}

func SearchProduct(name string) ([]Product, error) {
	db := getConn()
	rows, err := db.Query(context.Background(), fmt.Sprintf("SELECT * FROM products WHERE name LIKE '%%%s%%'", name))
	if err != nil {
		return nil, err
	}
	return scanManyData(rows, scanProduct)
}

// Add Update
func AddProduct(product Product) (Product, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "insert into products(name, description) values ($1, $2) returning *", product.Name, product.Description)
	return scanProduct(row)
}

func UpdateProduct(product_id int, product Product) (Product, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "update product set name = $1, description = $2 where product_id = $3 returning *", product.Name, product.Description, product_id)
	return scanProduct(row)
}

// Seller part
func SellerProducts(seller_id int, show_published bool) ([]ProductFromSeller, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
        SELECT ps.product_id, ps.seller_id, ps.quantity, ps.cost, ps.published
        FROM product_seller ps
        JOIN users u ON ps.seller_id = u.id
        WHERE u.id = $1 and ps.published = $1`, seller_id, show_published)
	if err != nil {
		return nil, err
	}
	return scanManyData(rows, scanProductFromSeller)
}

func UpdateSellerProduct(seller_id, product_id int, product ProductFromSeller) (ProductFromSeller, error) {
	db := getConn()
	var flag int
	err := db.QueryRow(context.Background(), "select count(*) from product_seller where product_id = $1 and seller_id = $2", product_id, seller_id).Scan(&flag)
	if err != nil {
		return product, err
	}
	statement := "insert into product_seller (quantity, cost, product_id, seller_id, published) values($1, $2, $3, $4, $5) returning *"
	if flag > 0 {
		statement = "update product_seller set quantity = $1, cost = $2, published = $5 where product_id = $3 and seller_id = $4 returning *"
	}
	row := db.QueryRow(context.Background(), statement, product.Quantity, product.Cost, product_id, seller_id, product.Published)
	return scanProductFromSeller(row)
}

func DeleteSellerProduct(product_id, seller_id int) (ProductFromSeller, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "delete from product_seller WHERE product_id = $1 and seller_id = $2 returning *", product_id, seller_id)
	return scanProductFromSeller(row)
}

func PublishSellerProduct(product_id, seller_id int, value, inverse bool) (ProductFromSeller, error) {
	db := getConn()
	statement_value := "not published"
	if !inverse {
		statement_value = strconv.FormatBool(value)
	}
	query := fmt.Sprintf("update product_seller set published = %s where product_id = %d and seller_id = %d returning *\n", statement_value, product_id, seller_id)
	row := db.QueryRow(context.Background(), query)
	return scanProductFromSeller(row)
}

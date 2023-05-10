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

func scanManyProductFromSeller(rows pgx.Rows) ([]ProductFromSeller, error) {
	sellerProducts := make([]ProductFromSeller, 0)
	for rows.Next() {
		sellerProduct, err := scanProductFromSeller(rows)
		if err != nil {
			return nil, err
		}
		sellerProducts = append(sellerProducts, sellerProduct)
	}
	return sellerProducts, rows.Err()
}

func scanProduct(row pgx.Row) (Product, error) {
	var product Product
	err := row.Scan(&product.Id, &product.Name, &product.Description)
	return product, err
}

func scanManyProduct(rows pgx.Rows) ([]Product, error) {
	products := make([]Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
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
	defer rows.Close()

	sellerProducts, err := scanManyProductFromSeller(rows)
	if err != nil {
		return nil, err
	}
	return sellerProducts, nil
}

func ListProducts() ([]Product, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products, err := scanManyProduct(rows)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func SearchProduct(name string) ([]Product, error) {
	db := getConn()
	rows, err := db.Query(context.Background(), fmt.Sprintf("SELECT * FROM products WHERE name LIKE '%%%s%%'", name))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products, err := scanManyProduct(rows)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func SellerProducts(seller_id int) ([]ProductFromSeller, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
        SELECT ps.product_id, ps.seller_id, ps.quantity, ps.cost, ps.published
        FROM product_seller ps
        JOIN users u ON ps.seller_id = u.id
        WHERE u.id = $1`, seller_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sellerProducts, err := scanManyProductFromSeller(rows)
	if err != nil {
		return nil, err
	}
	return sellerProducts, nil
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
	product, err = scanProductFromSeller(row)
	return product, err
}

func DeleteSellerProduct(product_id, seller_id int) (ProductFromSeller, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "delete from product_seller WHERE product_id = $1 and seller_id = $2 returning *", product_id, seller_id)
	product, err := scanProductFromSeller(row)
	return product, err
}

func PublishSellerProduct(product_id, seller_id int, value, inverse bool) (ProductFromSeller, error) {
	db := getConn()
	statement_value := "not published"
	if !inverse {
		statement_value = strconv.FormatBool(value)
	}
	query := fmt.Sprintf("update product_seller set published = %s where product_id = %d and seller_id = %d returning *\n", statement_value, product_id, seller_id)
	row := db.QueryRow(context.Background(), query)
	product, err := scanProductFromSeller(row)
	return product, err
}

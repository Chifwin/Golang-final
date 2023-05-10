package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ProductFromSeller struct {
	ProductId int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Cost      float64 `json:"cost"`
}

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func scanProductFromSeller(row pgx.Row) (ProductFromSeller, error) {
	var product ProductFromSeller
	err := row.Scan(&product.ProductId, &product.Quantity, &product.Cost)
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

func SellerProducts(seller_id int) ([]ProductFromSeller, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
        SELECT ps.product_id, ps.quantity, ps.cost
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

func ProductSellers(product_id int) ([]ProductFromSeller, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
	SELECT ps.seller_id, ps.quantity, ps.cost
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

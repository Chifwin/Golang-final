package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	// "time"
)

func scanProduct(rows pgx.Row) (Products, error) {
	var product Products
	err := rows.Scan(&product.ProductId, &product.Cost, &product.Quantity)
	return product, err
}

func UpdateSellerProductDB(seller_id, product_id int, product Products, value string) error {
	db := getConn()
	row := db.QueryRow(context.Background(), "select exists(select 1 from product_seller where product_id = $1 and seller_id = $2)", product_id, seller_id)
	var flag bool
	row.Scan(&flag)
	var pr pgx.Row
	if flag {
		pr = db.QueryRow(context.Background(), "update product_seller set quantity = $1, cost = $2, published = $5 where product_id = $3 and seller_id = $4 and product_id in (select id from products) returning *",
			product.Quantity, product.Cost, product_id, seller_id, value == "true")
	} else {
		pr = db.QueryRow(context.Background(), "insert into product_seller (quantity, cost, product_id, seller_id, published) values($1, $2, $3, $4, $5)",
			product.Quantity, product.Cost, product_id, seller_id,  value == "true")
	}
	_, err := scanProduct(pr)

	return err
}

func DeleteSellerProductDB(id, seller_id int) error {
	db := getConn()
	row := db.QueryRow(context.Background(), "DELETE FROM product_seller WHERE product_id = $1 and seller_id = $2 and product_id in (select id from products) returning *", id, seller_id)
	_, err := scanComment(row)
	return err
}


func PurchasesSellerDB(seller_id int) ([]Purchase, error){
	db := getConn()
	
	rows, err := db.Query(context.Background(), "SELECT * FROM purchases where seller_id=$1", seller_id)
	purchases := make([]Purchase, 0)
	if err != nil {
		if err == pgx.ErrNoRows {
			return purchases, nil
		}
		return nil, err
	}

	for rows.Next() {
		purchase, err := scanPurchase(rows)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}

	return purchases, nil
}
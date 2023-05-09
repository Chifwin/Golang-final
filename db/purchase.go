package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Purchase struct {
	ID        int       `json:"purchase_id"`
	UserID    int       `json:"user_id"`
	ProductId int       `json:"product_id" binding:"required"`
	SellerId  int       `json:"seller_id" binding:"required"`
	Date      time.Time `json:"date"`
	Quantity  int       `json:"quantity" binding:"required"`
}


func scanPurchase(rows pgx.Row) (Purchase, error) {
	var purchase Purchase
	err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.ProductId, &purchase.SellerId, &purchase.Date, &purchase.Quantity)
	return purchase, err
}

func GetBuyerPurchases(buyer_id int) ([]Purchase, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM purchases where buyer_id=$1", buyer_id)
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

func CreatePurchase(purchase Purchase) (Purchase, error) {
	db := getConn()
	fmt.Println(purchase)
	row := db.QueryRow(context.Background(), "select * from buy($1, $2, $3, $4)",
		purchase.UserID, purchase.SellerId, purchase.ProductId, purchase.Quantity)
	purchase, err := scanPurchase(row)
	return purchase, err
}

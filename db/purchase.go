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
	return purchase, handleError(err)
}

func filterPurchases(condition string) ([]Purchase, error) {
	db := getConn()
	rows, err := db.Query(context.Background(), "SELECT * FROM purchases where "+condition)
	return scanManyData(rows, err, scanPurchase)
}

func BuyerPurchases(buyer_id int) ([]Purchase, error) {
	return filterPurchases(fmt.Sprintf("buyer_id = %d", buyer_id))
}

func SellerPurchases(seller_id int) ([]Purchase, error) {
	return filterPurchases(fmt.Sprintf("seller_id = %d", seller_id))
}

func Buy(purchase Purchase) (Purchase, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "select * from buy($1, $2, $3, $4)",
		purchase.UserID, purchase.SellerId, purchase.ProductId, purchase.Quantity)
	purchase, err := scanPurchase(row)
	return purchase, err
}

func LastPurchases() ([]Purchase, error) {
	db := getConn()
	rows, err := db.Query(context.Background(), "select * from purchases order by date desc ")
	return scanManyData(rows, err, scanPurchase)
}

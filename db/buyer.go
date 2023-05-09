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

func GetPurchases(buyer_id int) ([]Purchase, error) {
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

func UpdateCommentDB(id int) (Scores, error) {
	db := getConn()
	var score Scores
	_, err := db.Exec(context.Background(), "UPDATE scores SET purchase_id=$1, rating=$2, comment=$3 WHERE product_id=$4",
		score.ProductId, score.Rating, score.Comment, id, BUYER)
	if err != nil {
		return score, err
	}

	return score, nil
}

func GetCommentDB() ([]Scores, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM scores", BUYER)
	if err != nil {
		return nil, err
	}

	var scores []Scores
	for rows.Next() {
		var score Scores
		err = rows.Scan(&score.ProductId, &score.Rating, &score.Comment)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}

func CreateCommentDB() (Scores, error) {
	db := getConn()
	var score Scores
	_, err := db.Exec(context.Background(), "INSERT INTO scores (purchase_id, rating, comment) VALUES ($1, $2, $3)",
		score.ProductId, score.Rating, score.Comment, BUYER)
	if err != nil {
		return score, err
	}

	return score, nil
}

func DeleteCommentDB(id int) (Scores, error) {
	db := getConn()
	var score Scores
	_, err := db.Exec(context.Background(), "DELETE FROM scores WHERE purchase_id=$1", id, BUYER)
	if err != nil {
		return score, err
	}

	return score, nil
}

package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Comment struct {
	PurchaseId int     `json:"purchase_id"`
	Rating     float64 `json:"rating" binding:"required"`
	Comment    string  `json:"comment"`
}

var ErrBadUser error = errors.New("the user have not access")

func scanComment(rows pgx.Row) (Comment, error) {
	var comment Comment
	err := rows.Scan(&comment.PurchaseId, &comment.Rating, &comment.Comment)
	return comment, err
}

func checkUserBelongComment(user_id, purchase_id int) bool {
	db := getConn()
	var real_user_id int
	err := db.QueryRow(context.Background(), "SELECT buyer_id from purchases where id=$1", purchase_id).Scan(&real_user_id)
	if err != nil || real_user_id != user_id {
		return false
	}
	return true
}

func GetBuyerComments(buyer_id int) ([]Comment, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM scores WHERE purchase_id in (SELECT id FROM purchases where buyer_id=$1)", buyer_id)
	comments := make([]Comment, 0)
	if err != nil {
		if err == pgx.ErrNoRows {
			return comments, nil
		}
		return nil, err
	}

	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func CreateComment(purchase_id, buyer_id int, comment Comment) (Comment, error) {
	db := getConn()
	if !checkUserBelongComment(buyer_id, purchase_id) {
		return comment, ErrBadUser
	}
	row := db.QueryRow(context.Background(), "insert into scores (purchase_id, rating, comment) values ($1, $2, $3) returning *",
		purchase_id, comment.Rating, comment.Comment)
	comment, err := scanComment(row)
	return comment, err
}

func UpdateComment(purchase_id, buyer_id int, comment Comment) (Comment, error) {
	db := getConn()
	if !checkUserBelongComment(buyer_id, purchase_id) {
		return comment, ErrBadUser
	}
	row := db.QueryRow(context.Background(), "UPDATE scores SET rating=$1, comment=$2 WHERE purchase_id=$3 returning *",
		comment.Rating, comment.Comment, purchase_id)
	comment, err := scanComment(row)
	return comment, err
}

func DeleteComment(purchase_id, buyer_id int) (Comment, error) {
	db := getConn()
	var comment Comment
	if !checkUserBelongComment(buyer_id, purchase_id) {
		return comment, ErrBadUser
	}
	row := db.QueryRow(context.Background(), "DELETE FROM scores WHERE purchase_id = $1 returning *", purchase_id)
	comment, err := scanComment(row)
	return comment, err
}

func SellerComments(seller_id int) ([]Comment, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
		SELECT p.id, s.rating, s.comment
		FROM scores s
		JOIN purchases p ON s.purchase_id = p.id
		WHERE p.seller_id = $1`, seller_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

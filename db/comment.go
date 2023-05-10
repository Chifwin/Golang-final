package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrBadUser error = errors.New("the user have not access")

func scanComment(rows pgx.Row) (Scores, error) {
	var scores Scores
	err := rows.Scan(&scores.PurchaseId, &scores.Rating, &scores.Comment)
	return scores, err
}

func checkUserBelongComment(user_id, purchase_id int) bool {
	db := getConn()
	var realUserId int
	err := db.QueryRow(context.Background(), "SELECT buyer_id from purchases where purchase_id=$1", purchase_id).Scan(&realUserId)
	if err != nil || realUserId != user_id {
		return false
	}
	return true
}

func GetBuyerComments(buyer_id int) ([]Scores, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM scores WHERE purchase_id in (SELECT id FROM purchase where buyer_id=$1)", buyer_id)
	scores := make([]Scores, 0)
	if err != nil {
		if err == pgx.ErrNoRows {
			return scores, nil
		}
		return nil, err
	}

	for rows.Next() {
		score, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}

func CreateComment(score Scores, buyer_id int) (Scores, error) {
	db := getConn()
	if !checkUserBelongComment(buyer_id, score.PurchaseId) {
		return score, ErrBadUser
	}
	row := db.QueryRow(context.Background(), "insert into score (purchase_id, rating, comment) value($1, $2, $3)",
		score.PurchaseId, score.Rating, score.Comment)
	score, err := scanComment(row)
	return score, err
}

func UpdateCommentDB(purchase_id int, buyer_id int, scores Scores) (Scores, error) {
	db := getConn()
	var score Scores
	if !checkUserBelongComment(buyer_id, purchase_id) {
		return score, ErrBadUser
	}
	row := db.QueryRow(context.Background(), "UPDATE scores SET rating=$1, comment=$2 WHERE purchase_id=$3 returning *",
		score.Rating, score.Comment, purchase_id)
	score, err := scanComment(row)
	return score, err
}

func DeleteCommentDB(purchase_id, buyer_id int) (Scores, error) {
	db := getConn()
	var score Scores
	if !checkUserBelongComment(buyer_id, purchase_id) {
		return score, ErrBadUser
	}
	row := db.QueryRow(context.Background(), "DELETE FROM scores WHERE purchase_id = $1", purchase_id)
	score, err := scanComment(row)
	return score, err
}

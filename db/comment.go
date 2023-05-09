package db

import (
	"context"
	"fmt"


	"github.com/jackc/pgx/v5"
)

func scanComment(rows pgx.Row) (Scores, error) {
	var scores Scores
	err := rows.Scan(&scores.ProductId, &scores.Rating, &scores.Comment)
	return scores, err
}

func GetBuyerComments(buyer_id int) ([]Scores, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM scores WHERE purchase_id=$1", buyer_id)
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

func CreateComment(score Scores) (Scores, error) {
	db := getConn()
	fmt.Println(score)
	row := db.QueryRow(context.Background(), "select * from score($1, $2, $3)",
		score.ProductId, score.Rating, score.Comment)
	score, err := scanComment(row)
	return score, err
}

func UpdateCommentDB(id int, scores Scores) (Scores, error) {
	db := getConn()
	var score Scores
	_, err := db.Exec(context.Background(), "UPDATE scores SET purchase_id=$1, rating=$2, comment=$3 WHERE product_id=$4",
		score.ProductId, score.Rating, score.Comment, id, BUYER)
	if err != nil {
		return score, err
	}

	return score, nil
}


// func CreateCommentDB() (Scores, error) {
// 	db := getConn()
// 	var score Scores
// 	_, err := db.Exec(context.Background(), "INSERT INTO scores (purchase_id, rating, comment) VALUES ($1, $2, $3)",
// 		score.ProductId, score.Rating, score.Comment, BUYER)
// 	if err != nil {
// 		return score, err
// 	}

// 	return score, nil
// }

func DeleteCommentDB(id int) (Scores, error) {
	db := getConn()
	var score Scores
	_, err := db.Exec(context.Background(), "DELETE FROM scores WHERE purchase_id=$1", id, BUYER)
	if err != nil {
		return score, err
	}

	return score, nil
}

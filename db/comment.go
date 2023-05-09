package db

import (
	"context"
)

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

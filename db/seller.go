package db

import (
	"context"
)

type Sellers struct {
	ID   int
	Name string
}

type Scores struct {
	ProductId int
	Rating    float64
	Comment   string
}

func ListSellers() ([]Sellers, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT id, name FROM users WHERE role = $1", SELLER)
	if err != nil {
		return nil, err
	}

	var sellers []Sellers
	for rows.Next() {
		var seller Sellers
		err = rows.Scan(&seller.ID, &seller.Name)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func FindSellerById(id int) (*Sellers, error) {
	db := getConn()
	var seller Sellers
	err := db.QueryRow(context.Background(), "select id, name from users where role = $1 and id = $2", SELLER, id).Scan(&seller.ID, &seller.Name)
	if err != nil {
		return nil, err
	}

	return &seller, nil

}

func SellerScores(id int) ([]Scores, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
		SELECT p.id, s.rating, s.comment
		FROM scores s
		JOIN purchases p ON s.purchase_id = p.id
		WHERE p.seller_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []Scores
	for rows.Next() {
		var score Scores
		err = rows.Scan(&score.ProductId, &score.Rating, &score.Comment)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return scores, nil
}

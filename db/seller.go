package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Seller struct {
	ID   int
	Name string
}

func scanSeller(row pgx.Row) (Seller, error) {
	var seller Seller
	err := row.Scan(&seller.ID, &seller.Name)
	return seller, err
}

func ListSellers() ([]Seller, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT id, name FROM users WHERE role = $1", SELLER)
	if err != nil {
		return nil, err
	}
	return scanManyData(rows, scanSeller)
}

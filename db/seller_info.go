package db

import (
	"context"
)

type Sellers struct {
	ID   int
	Name string
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

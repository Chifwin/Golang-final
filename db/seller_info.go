package db

import (
	"context"
)

type Sellers struct {
	ID   int
	Name string
}



type Products struct {
	ProductId int     `json:"productId"`
	Quantity  uint64  `json:"quantity"`
	Cost      float64 `json:"cost"`
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

func SellerProducts(id int) ([]Products, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
        SELECT ps.product_id, ps.quantity, ps.cost
        FROM product_seller ps
        JOIN users u ON ps.seller_id = u.id
        WHERE u.id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sellerProducts []Products
	for rows.Next() {
		var sellerProduct Products
		err = rows.Scan(&sellerProduct.ProductId, &sellerProduct.Quantity, &sellerProduct.Cost)
		if err != nil {
			return nil, err
		}
		sellerProducts = append(sellerProducts, sellerProduct)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return sellerProducts, nil

}

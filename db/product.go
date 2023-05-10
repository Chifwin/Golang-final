package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ProductFromSeller struct {
	ProductId int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Cost      float64 `json:"cost"`
}

func scanProduct(rows pgx.Row) (ProductFromSeller, error) {
	var product ProductFromSeller
	err := rows.Scan(&product.ProductId, &product.Quantity, &product.Cost)
	return product, err
}

func SellerProducts(seller_id int) ([]ProductFromSeller, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
        SELECT ps.product_id, ps.quantity, ps.cost
        FROM product_seller ps
        JOIN users u ON ps.seller_id = u.id
        WHERE u.id = $1`, seller_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sellerProducts := make([]ProductFromSeller, 0)
	for rows.Next() {
		sellerProduct, err := scanProduct(rows)
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

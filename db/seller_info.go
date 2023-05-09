package db

import (
	"context"
	"time"
)

type Sellers struct {
	ID   int
	Name string
}

type Scores struct {
	ProductId int       `json:"purchase_id" binding:"required"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
	Date      time.Time `json:"date"`
}

type Products struct {
	ProductId int     `json:"productId"`
	Quantity  uint64  `json:"quantity"`
	Cost      float64 `json:"cost"`
}

type ListProduct struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

func SellerScores(id int) ([]Scores, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
		SELECT p.id, s.rating, s.comment , p.date
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
		err = rows.Scan(&score.ProductId, &score.Rating, &score.Comment, &score.Date)
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

func ListProducts() ([]ListProduct, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	var products []ListProduct
	for rows.Next() {
		var product ListProduct
		err = rows.Scan(&product.Id, &product.Name, &product.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func SearchProduct() ([]ListProduct, error){
	db := getConn()

	rows, err := db.Query(context.Background(), "SELECT name FROM products")
	if err != nil {
		return nil, err
	}

	var products []ListProduct
	for rows.Next() {
		var product ListProduct
		err = rows.Scan(&product.Id, &product.Name, &product.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func ProductSellers(id int) ([]Products, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
        SELECT ps.seller_id, ps.quantity, ps.cost
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

func ProductScores(id int) ([]Scores, error) {
	db := getConn()

	rows, err := db.Query(context.Background(), `
		SELECT s.id, p.rating, p.comment , p.date
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
		err = rows.Scan(&score.ProductId, &score.Rating, &score.Comment, &score.Date)
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


func ProductSearchDB(name string) ([]ListProduct, error) {
	db := getConn()
	name = "%" + name + "%"
	rows, err := db.Query(context.Background(), `
        SELECT * from products where name like $1`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ListProduct
	for rows.Next() {
		var product ListProduct
		err = rows.Scan(&product.Id, &product.Name, &product.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}
package db

import (
	"context"
	"fmt"
)

type UserRole string

const (
	ADMIN  UserRole = "admin"
	SELLER UserRole = "seller"
	BUYER  UserRole = "buyer"
)

type UserRet struct {
	ID       int
	Username string
	Name     string
	Role     UserRole
}

type UserCred struct {
	Username string
	Name     string
	Role     UserRole
	Password string
}

func AddUser(user UserCred) error {
	db := getConn()
	_, err := db.Exec(context.Background(), "call add_user($1, $2, $3, $4)", user.Username, user.Password, user.Password, user.Role)
	return err
}

func AuthoriseUser(username, password string) (*UserRet, error) {
	var res UserRet
	db := getConn()
	fmt.Printf("Login with username: %s and password %s\n", username, password)
	err := db.QueryRow(context.Background(), "select * from authorise_user($1, $2)", username, password).Scan(&res.ID, &res.Username, &res.Name, &res.Role)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

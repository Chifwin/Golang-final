package db

import (
	"context"
	"errors"
	"final/structs"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("not found")
)

func AddUser(user structs.UserCred) error {
	db := getConn()
	_, err := db.Exec(context.Background(), "call add_user($1, $2, $3, $4)", user.Username, user.Password, user.Password, user.Role)
	return err
}

func AuthoriseUser(username, password string) (*structs.UserRet, error) {
	var res structs.UserRet
	db := getConn()
	fmt.Printf("Login with username: %s and password %s\n", username, password)
	err := db.QueryRow(context.Background(), "select * from authorise_user($1, $2)", username, password).Scan(&res.ID, &res.Username, &res.Name, &res.Role)
	switch err {
	case nil:
		return &res, nil
	case pgx.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

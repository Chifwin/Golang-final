package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRole string

const (
	ADMIN  UserRole = "admin"
	SELLER UserRole = "seller"
	BUYER  UserRole = "buyer"
)

type UserCred struct {
	Username string   `json:"username"  binding:"required"`
	Name     string   `json:"name"  binding:"required"`
	Role     UserRole `json:"role"  binding:"required"`
	Password string   `json:"password"  binding:"required"`
}

type UserRet struct {
	ID       int
	Username string
	Name     string
	Role     UserRole
}

func scanUserRet(row pgx.Row) (UserRet, error) {
	var user UserRet
	err := row.Scan(&user.ID, user.Username, user.Name, user.Role)
	return user, err
}

func AddUser(user UserCred) error {
	db := getConn()
	_, err := db.Exec(context.Background(), "call add_user($1, $2, $3, $4)", user.Username, user.Password, user.Password, user.Role)
	return err
}

func AuthoriseUser(username, password string) (UserRet, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "select * from authorise_user($1, $2)", username, password)
	return scanUserRet(row)
}

func GetAllUsers() ([]UserRet, error) {
	db := getConn()
	rows, err := db.Query(context.Background(), "select id, username, name, role from users")
	if err != nil {
		return nil, err
	}
	return scanManyData(rows, scanUserRet)
}

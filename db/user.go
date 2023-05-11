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

	// Not in database
	NO_ROLE  UserRole = ""
	ANY_ROLE UserRole = "any"
)

var ROLES [3]UserRole = [...]UserRole{ADMIN, SELLER, BUYER}

func ValidRole(s string) (UserRole, bool) {
	for _, role := range ROLES {
		if s == string(role) {
			return role, true
		}
	}
	return NO_ROLE, false
}

type UserCred struct {
	Username string   `json:"username" binding:"required"`
	Name     string   `json:"name" binding:"required"`
	Role     UserRole `json:"role"`
	Password string   `json:"password" binding:"required"`
}

type UserRet struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Role     UserRole `json:"role"`
}

func scanUserRet(row pgx.Row) (UserRet, error) {
	var user UserRet
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Role)
	return user, handleError(err)
}

func GetUser(user_id int) (UserRet, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "select id, username, name, role from users where id = $1", user_id)
	return scanUserRet(row)
}

func AddUser(user UserCred) (UserRet, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "insert into users (username, password, name, role) values ($1, $2, $3, $4) returning id, username, name, role", user.Username, user.Password, user.Name, user.Role)
	return scanUserRet(row)
}

func UpdateUser(user_id int, user UserCred) (UserRet, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "update users set username = $1, password = $2, name = $3 where id = $4 returning id, username, name, role", user.Username, user.Password, user.Name, user_id)
	return scanUserRet(row)
}

func DeleteUser(user_id int) (UserRet, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "delete from users where id = $1 returning id, username, name, role", user_id)
	return scanUserRet(row)
}

func AuthoriseUser(username, password string) (UserRet, error) {
	db := getConn()
	row := db.QueryRow(context.Background(), "select * from authorise_user($1, $2)", username, password)
	return scanUserRet(row)
}

func GetAllUsers() ([]UserRet, error) {
	db := getConn()
	rows, err := db.Query(context.Background(), "select id, username, name, role from users")
	return scanManyData(rows, err, scanUserRet)
}

package db

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Errors
var (
	ErrNotFound     error = errors.New("not founded")
	ErrUniqueFailed error = errors.New("unique constraint failed")
	ErrBuyError     error = errors.New("seller do not have such product or it is not published")
)

func handleError(err error) error {
	if err == pgx.ErrNoRows {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName != "" {
			return ErrUniqueFailed
		}
		if strings.ToLower(pgErr.Message) == ErrBuyError.Error() {
			return ErrBuyError
		}
	}
	return err
}

var postgres_conn *pgx.Conn

func getConn() *pgx.Conn {
	// Reuse existing connection because of connection limit to database
	if postgres_conn == nil {
		var err error
		postgres_conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	return postgres_conn
}

func scanManyData[T any](rows pgx.Rows, rows_err error, scanOne func(pgx.Row) (T, error)) ([]T, error) {
	datas := make([]T, 0)
	defer rows.Close()
	if rows_err != nil {
		return datas, handleError(rows_err)
	}
	for rows.Next() {
		data, err := scanOne(rows)
		if err != nil {
			return nil, handleError(err)
		}
		datas = append(datas, data)
	}
	return datas, handleError(rows.Err())
}

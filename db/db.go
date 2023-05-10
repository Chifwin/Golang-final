package db

import (
	"context"
	"log"
	"os"

	pgx "github.com/jackc/pgx/v5"
)

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

func scanManyData[T any](rows pgx.Rows, scanOne func(pgx.Row) (T, error)) ([]T, error) {
	datas := make([]T, 0)
	defer rows.Close()
	for rows.Next() {
		data, err := scanOne(rows)
		if err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}
	return datas, rows.Err()
}

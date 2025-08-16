package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Connectdb() (pgx.Conn, error) {
	pglink := "postgres://a4bhi:a4bhi@localhost:5432/echoai"
	conn, err := pgx.Connect(context.Background(), pglink)

	return *conn, err
}

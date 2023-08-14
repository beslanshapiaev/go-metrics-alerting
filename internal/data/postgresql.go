package data

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PgModule struct {
	conn *pgx.Conn
}

func NewPgModule(connectionString string) *PgModule {
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprint(os.Stderr, "Unable to connect to database: ", err)
		os.Exit(1)
	}
	return &PgModule{
		conn: conn,
	}
}

func (m *PgModule) Ping() bool {
	if err := m.conn.Ping(context.Background()); err != nil {
		return false
	}
	return true
}

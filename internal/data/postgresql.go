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
		return nil
	}
	return &PgModule{
		conn: conn,
	}
}

func (m *PgModule) Ping() bool {
	if m.conn != nil {
		return false
	}
	if err := m.conn.Ping(context.Background()); err != nil {
		return false
	}
	return true
}

func (m *PgModule) Close() {
	m.conn.Close(context.Background())
}
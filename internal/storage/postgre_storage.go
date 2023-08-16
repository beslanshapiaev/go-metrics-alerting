package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/beslanshapiaev/go-metrics-alerting/common"
	"github.com/jackc/pgx/v5"
)

type PostgreStorage struct {
	mu       sync.RWMutex
	filePath string
	conn     pgx.Conn
}

func NewPostgreStorage(connString string, filePath string) *PostgreStorage {
	connection, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer conn.Close(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	// defer cancel()
	connection.Exec(context.Background(), "CREATE SCHEMA IF NOT EXISTS practicum AUTHORIZATION pg_database_owner; \n")
	connection.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS practicum.metrics"+
		"(id varchar(40), type varchar(40), delta integer, value double precision)")
	return &PostgreStorage{
		conn:     *connection,
		filePath: filePath,
	}
}

func (s *PostgreStorage) AddGaugeMetric(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var err error

	// conn, err := pgx.Connect(context.Background(), s.connectionString)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = s.conn.Exec(ctx, "DELETE FROM practicum.metrics where id = $1", name)
	if err != nil {
		panic(err)
	}

	_, err = s.conn.Exec(ctx, "INSERT INTO practicum.metrics (id, type, value) VALUES ($1, $2, $3)",
		name, "gauge", value)
	if err != nil {
		panic(err)
	}
}

func (s *PostgreStorage) AddCounterMetric(name string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = s.conn.Exec(ctx, "DELETE FROM practicum.metrics where id = $1", name)
	if err != nil {
		panic(err)
	}
	_, err = s.conn.Exec(ctx, "INSERT INTO practicum.metrics (id, type, delta) VALUES ($1, $2, $3)",
		name, "counter", value)
	if err != nil {
		panic(err)
	}
}

func (s *PostgreStorage) AddMetricsBatch(metrics []common.Metric) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	// ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	// defer cancel()

	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, v := range metrics {
		_, err := tx.Exec(context.Background(), "delete from practicum.metrics where id = $1", v.ID)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
		_, err = tx.Exec(context.Background(), "insert into practicum.metrics (id, type, delta, value) values ($1, $2, $3, $4)", v.ID, v.MType, v.Delta, v.Value)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (s *PostgreStorage) GetGaugeMetric(name string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.conn.QueryRow(ctx, "SELECT * FROM practicum.metrics where type = $1 and id = $2", "gauge", name)

	var value sql.NullFloat64
	var typeValue string
	var id string
	var delta sql.NullInt64
	err = row.Scan(&id, &typeValue, &delta, &value)
	if err != nil {
		return 0, false
	}
	if !value.Valid {
		return 0, false
	}
	return value.Float64, true
}

func (s *PostgreStorage) GetCounterMetric(name string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.conn.QueryRow(ctx, "SELECT * FROM practicum.metrics where type = $1 and id = $2", "counter", name)

	var value sql.NullFloat64
	var typeValue string
	var id string
	var delta sql.NullInt64
	err = row.Scan(&id, &typeValue, &delta, &value)
	if err != nil {
		return 0, false
	}
	if !delta.Valid {
		return 0, false
	}
	return delta.Int64, true
}

func (s *PostgreStorage) GetAllMetrics() map[string]interface{} {
	resultMap := make(map[string]interface{}, 30)

	s.mu.RLock()
	defer s.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.conn.Query(ctx, "select * from practitcum.metrics")
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var m common.Metric
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return nil
		}
		if m.MType == "gauge" {
			resultMap[m.ID] = m.Value
		} else if m.MType == "counter" {
			resultMap[m.ID] = m.Delta
		}
	}
	return resultMap
}

func (s *PostgreStorage) SaveToFile() error {
	return nil
}

func (s *PostgreStorage) RestoreFromFile() error {
	return nil
}

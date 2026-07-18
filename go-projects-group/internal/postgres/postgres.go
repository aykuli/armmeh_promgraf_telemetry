package postgres

import (
	"context"
	"fleet-app-gr/internal/payload"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	createEventsSql = `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		vehicle_id INT NOT NULL,
		vehicle_type VARCHAR(50) NOT NULL,
		event_type VARCHAR(100) NOT NULL,
		severity VARCHAR(30) NOT NULL,
		description TEXT,
		code VARCHAR(20),
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`
	insertEventSql = `
	INSERT INTO events (vehicle_id, vehicle_type, event_type, severity, description, code)
	VALUES ($1, $2, $3, $4, $5, $6);`
)

type DBStorage struct {
	instance *pgxpool.Pool
}

func NewStorage(dsn string) (*DBStorage, error) {
	ctx := context.Background()
	var s DBStorage
	tryCount := 0
	createConn := func() error {
		word := "try"
		if tryCount > 0 {
			word = "retry"
		}
		fmt.Printf("%s to connect to database, probe %d\n", word, tryCount)
		tryCount++

		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return fmt.Errorf("could not connect to database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			return fmt.Errorf("could not connect to database: %v", err)
		}

		fmt.Printf("  | -- connected to database %s\n", dsn)
		s.instance = pool

		return nil
	}
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 12 * time.Second
	if err := backoff.Retry(createConn, expBackoff); err != nil {
		return nil, fmt.Errorf("\nfailed to connect to database after retrying %d times: %v", tryCount, err)
	}
	if err := s.initTables(ctx); err != nil {
		return &s, err
	}

	return &s, nil
}

func (s *DBStorage) initTables(ctx context.Context) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Query(ctx, createEventsSql)
	return err
}

func (s *DBStorage) Close() {
	s.instance.Close()
}

func (s *DBStorage) Ping(ctx context.Context) error {
	return s.instance.Ping(ctx)
}

func (s *DBStorage) SaveEvent(ctx context.Context, p payload.TelemetryPayload) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, event := range p.Events {
		_, err = conn.Exec(ctx, insertEventSql, p.VehicleID, p.VehicleType,
			event.EventType, event.Severity, event.Description, event.Code)
	}

	return err
}

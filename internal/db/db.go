package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func Connect(cfg Config) (*pgxpool.Pool, error) {
	dns := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to database %s", cfg.Database)

	err = CreateTables(pool)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return pool, nil
}

func CreateTables(pool *pgxpool.Pool) error {
	rows, err := pool.Query(context.Background(), "CREATE TABLE IF NOT EXISTS tasks (id SERIAL PRIMARY KEY, name TEXT, schedule TEXT)")
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

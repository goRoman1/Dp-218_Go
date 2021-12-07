package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool     *pgxpool.Pool
}

func (db *Postgres) QueryResult(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return db.Pool.Query(ctx, query, args...)
}

func (db *Postgres) QueryResultRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.Pool.QueryRow(ctx, query, args...)
}

func (db *Postgres) QueryExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, query, args...)
}

func NewConnection(connectionString string) (*Postgres, error) {
	dbPg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("Postgres. Error parsing config url - %v", err)
	}
	poolConfig.MaxConns = int32(dbPg.maxPoolSize)
	for dbPg.connAttempts > 0 {
		dbPg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		log.Printf("Postgres. Trying to connect, attempts left: %d", dbPg.connAttempts)
		time.Sleep(dbPg.connTimeout)
		dbPg.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("Postgres. Failed to connect: %v", err)
	}

	return dbPg, nil
}

func (db *Postgres) CloseDB() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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

	QuerySQL string
	Pool     *pgxpool.Pool
}

func (pg *Postgres) QueryResult(ctx context.Context, args ...interface{}) (pgx.Rows, error) {
	return pg.Pool.Query(ctx, pg.QuerySQL, args...)
}

func (pg *Postgres) QueryResultRow(ctx context.Context, args ...interface{}) pgx.Row {
	return pg.Pool.QueryRow(ctx, pg.QuerySQL, args...)
}

func (pg *Postgres) QueryExec(ctx context.Context, args ...interface{}) (pgconn.CommandTag, error) {
	return pg.Pool.Exec(ctx, pg.QuerySQL, args...)
}

func NewPostgres(dbUri string) (*Postgres, error) {
	dbPg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(dbUri)
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

func (pg *Postgres) CloseDB() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}

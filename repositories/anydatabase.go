package repositories

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type AnyDatabase interface {
	QueryResult(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryResultRow(context.Context, string, ...interface{}) pgx.Row
	QueryExec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	CloseDB()
}
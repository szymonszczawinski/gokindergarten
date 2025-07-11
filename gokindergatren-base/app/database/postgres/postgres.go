// Package postgres
package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"golang.org/x/sync/errgroup"
)

type postgresDatabase struct {
	dbpool *pgxpool.Pool
	ctx    context.Context
	eg     *errgroup.Group
}

func NewPostgresSqlDatabase(eg *errgroup.Group, ctx context.Context) postgresDatabase {
	return postgresDatabase{
		eg:  eg,
		ctx: ctx,
	}
}

func (db *postgresDatabase) Close() {
	db.dbpool.Close()
}

func (db *postgresDatabase) Open() {
	db.dbpool = openDatabase(db.ctx)
}

func openDatabase(ctx context.Context) *pgxpool.Pool {
	tracer := &tracelog.TraceLog{
		Logger:   myLogger{},
		LogLevel: tracelog.LogLevelTrace,
	}

	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse connString: %v\n", err)
		os.Exit(1)
	}

	dbConfig.ConnConfig.Tracer = tracer
	dbpool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	// dbpool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	err = dbpool.Ping(ctx)
	if err != nil {

		fmt.Fprintf(os.Stderr, "Unable to PING connection pool: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

type myLogger struct{}

func (ll myLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	slog.Info("PGX", "level=", level, "msg=", msg, "args=", data)
}

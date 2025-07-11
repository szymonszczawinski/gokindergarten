// Package database
package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq" // posgres driver
)

type GenericDB struct {
	db *sql.DB
}

func NewGenericDb() GenericDB {
	return GenericDB{}
}

func (gdb *GenericDB) Open(connectionString string) (*sql.DB, error) {
	slog.Info("opening generic db with:", "connStr", connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not ping db: %w", err)
	}
	gdb.db = db
	return db, nil
}

func (gdb *GenericDB) Close() error {
	if gdb.db != nil {
		err := gdb.db.Close()
		if err != nil {
			return fmt.Errorf("could not close database: %w", err)
		}
	} else {
		return fmt.Errorf("database not initialized")
	}
	return nil
}

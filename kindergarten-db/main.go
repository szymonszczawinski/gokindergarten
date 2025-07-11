package main

import (
	"kindergarten-db/db"
	"kindergarten-db/db/migrations"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	slog.Info("Start DB migrations")
	runDBMigrations()
	slog.Info("DB migrations finished")
}

func runDBMigrations() {
	dbConnectionString := os.Getenv("DB_URL")
	gdb := db.NewGenericDb()
	sqldb, err := gdb.Open(dbConnectionString)
	if err != nil {
		log.Fatalf("could not open db: %v", err.Error())
	}

	err = migrations.RunMigrations(sqldb)
	if err != nil {
		slog.Error("migrations error", "err", err.Error())
	}
	gdb.Close()
}

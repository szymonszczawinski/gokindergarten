// Package app
package app

import (
	"context"
	"gokindergarten/app/database"
	"gokindergarten/app/database/postgres"
	"gokindergarten/db/migrations"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func Start(args []string) {
	slog.Info("Starting app")
	baseContext, cancel := context.WithCancel(context.Background())
	signalChannel := registerShutdownHook(cancel)
	mainGroup, mainContext := errgroup.WithContext(baseContext)

	postgres.NewPostgresSqlDatabase(mainGroup, mainContext)
	runDbMigrations()
	if err := mainGroup.Wait(); err == nil {
		slog.Info("Closing App")
	}

	defer close(signalChannel)
}

func runDbMigrations() {
	dbConnectionString := os.Getenv("DB_URL")
	gdb := database.NewGenericDb()
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

func registerShutdownHook(cancel context.CancelFunc) chan os.Signal {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		// wait until receiving the signal
		<-sigCh
		cancel()
	}()

	return sigCh
}

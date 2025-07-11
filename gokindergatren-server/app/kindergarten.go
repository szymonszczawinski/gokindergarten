// Package app
package app

import (
	"context"
	"embed"
	"gokindergarten/app/api"
	"gokindergarten/app/database"
	"gokindergarten/app/database/postgres"
	"gokindergarten/app/home"
	"gokindergarten/app/http"
	"gokindergarten/db/migrations"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func Start(args []string, publicDir embed.FS) {
	slog.Info("Starting app")
	baseContext, cancel := context.WithCancel(context.Background())
	signalChannel := registerShutdownHook(cancel)
	mainGroup, mainContext := errgroup.WithContext(baseContext)

	runDBMigrations()

	postgres.NewPostgresSqlDatabase(mainGroup, mainContext)

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		httpPort = api.DefaultHTTPServerPort
	}
	httpServer := http.NewHTTPServer(mainContext, mainGroup, httpPort, publicDir)

	homeHandler := home.NewHomeHandler()

	httpServer.AddHandler("/", homeHandler)
	httpServer.Start()

	if err := mainGroup.Wait(); err == nil {
		slog.Info("Closing App")
	}
	defer close(signalChannel)
}

func runDBMigrations() {
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

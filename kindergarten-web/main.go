package main

import (
	"embed"
	"kindergarten-web/app"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

//go:embed static/*
var publicDir embed.FS

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	slog.Info("Hello GO Kindergatren")
	app.Start(os.Args, publicDir)
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/NekichGO/Graduate-work/pkg/db"
	"github.com/NekichGO/Graduate-work/pkg/server"
	"github.com/joho/godotenv"
)

const (
	webDir        = "./web"
	defaultPort   = "7540"
	dbFileDefault = "scheduler.db"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := loadEnv(); err != nil {
		return err
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = dbFileDefault
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}
	defer db.Close()

	if err := server.Run(port, webDir); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func loadEnv() error {
	err := godotenv.Load(".env")
	switch {
	case err == nil:
		return nil
	case errors.Is(err, os.ErrNotExist):
		log.Println("No .env file found, environment variables or default values are used")
		return nil
	default:
		return fmt.Errorf("load .env: %w", err)
	}
}

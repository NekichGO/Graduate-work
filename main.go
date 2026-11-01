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

const webDir = "./web"

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
		return errors.New("TODO_DBFILE env variable empty")
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		return errors.New("variable TODO_PORT is empty")
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
	if err != nil {
		return fmt.Errorf("load .env: %w", err)
	}
	return nil
}

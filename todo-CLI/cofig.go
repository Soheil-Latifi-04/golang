package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN string
}

func LoadConfig() *Config {
	// .env is optional — fine if it's missing (e.g. in prod where
	// env vars are set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is not set")
	}

	return &Config{DatabaseDSN: dsn}
}
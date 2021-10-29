package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func setupFlags(cfg *config) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.StringVar(&cfg.port, "port", os.Getenv("PORT"), "API Server PORT")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "Server Environment {development|production}")
	flag.StringVar(&cfg.domain, "domain", os.Getenv("DOMAIN"), "Application domain")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DUKKAN_DB_DSN"), "Database DSN")
	flag.Parse()
}

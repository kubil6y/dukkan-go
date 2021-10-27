package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Infof("%s server is running on port :%s", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

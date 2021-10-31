package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	//$ go run ./cmd/api -cors-trusted-origins="https://www.example.com https://staging.example.com"
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

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

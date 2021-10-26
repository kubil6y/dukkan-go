package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func setupFlags(cfg *config) {
	flag.IntVar(&cfg.port, "port", 4000, "API Server PORT")
	flag.StringVar(&cfg.env, "env", "development", "Server Environment {development|production}")
	flag.Parse()
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Infof("%s server is running on port :%d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

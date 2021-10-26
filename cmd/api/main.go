package main

import (
	"time"

	"github.com/kubil6y/dukkan-go/internal/data"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config  config
	logger  *zap.SugaredLogger
	models  data.Models
	version string
}

func main() {
	var cfg config
	setupFlags(&cfg)

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.Stamp)
	logger, _ := config.Build()
	sugar := logger.Sugar()

	app := &application{
		config:  cfg,
		logger:  sugar,
		version: version,
	}

	if err := app.serve(); err != nil {
		app.logger.Fatalf("failed to start %s server", app.config.env)
	}
}

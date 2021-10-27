package main

import (
	"sync"
	"time"

	"github.com/kubil6y/dukkan-go/internal/data"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const version = "1.0.0"

type config struct {
	port string
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
	wg      sync.WaitGroup
}

func main() {
	var cfg config
	setupFlags(&cfg)

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.Stamp)
	logger, _ := loggerConfig.Build()
	sugar := logger.Sugar()

	db, err := connectDatabase(cfg)
	if err != nil {
		sugar.Fatal(err)
	}
	autoMigrate(db)

	app := &application{
		config:  cfg,
		logger:  sugar,
		version: version,
		models:  data.NewModels(db),
	}

	if err := app.serve(); err != nil {
		app.logger.Fatalf("failed to start %s server", app.config.env)
	}
}

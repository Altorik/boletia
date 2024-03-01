package bootstrap

import (
	"boletia/internal/obtain"
	"boletia/internal/platform/http/api"
	"boletia/internal/platform/server"
	"boletia/internal/platform/storage/postgres"
	"boletia/kit/logger"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"time"
)

func Run() error {
	var cfg config
	err := envconfig.Process("boletia", &cfg)
	if err != nil {
		return err
	}
	log.Println("cfg", cfg)
	databaseURI := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
	dbpool, err := pgxpool.New(context.Background(), databaseURI)
	if err != nil {
		return err
	}

	defer dbpool.Close()

	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	httpClient := &http.Client{Transport: httpTransport, Timeout: cfg.ApiTimeout}

	loggerApp, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	loggerApp = loggerApp.With("service", "boletia")

	currencyRepository := api.NewCurrencyRepository(httpClient, loggerApp, cfg.ApiUrl, cfg.ApiKey)
	databaseRepository := postgres.NewDatabaseRepository(dbpool, cfg.DbTimeout, cfg.DbTableName)
	getCurrencyService := obtain.NewCurrencyService(currencyRepository, databaseRepository)

	ctx, srv := server.New(context.Background(), cfg.Host, cfg.Port, cfg.ShutdownTimeout, getCurrencyService)
	return srv.Run(ctx)
}

type config struct {
	// Server configuration
	Host            string        `default:"0.0.0.0"`
	Port            uint          `default:"80"`
	ShutdownTimeout time.Duration `default:"10s"`
	// Database configuration
	DbUser      string        `default:"altorik"`
	DbPass      string        `default:"superpass"`
	DbHost      string        `default:"localhost"`
	DbName      string        `default:"boletia"`
	DbPort      uint          `default:"5432"`
	DbTimeout   time.Duration `default:"6s"`
	DbTableName string        `default:"temp_currency_rates"`
	// API config
	ApiUrl     string        `default:"https://api.currencyapi.com/v3/latest?apikey="`
	ApiKey     string        `default:"cur_live_NCuvBT4JiXLzvjPg7kubfgpmgLU4x9g8Wu1zJjd9"`
	ApiTimeout time.Duration `default:"3s"`
}

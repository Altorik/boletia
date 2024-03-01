package bootstrap

import (
	"boletia/api/internal/obtain"
	"boletia/api/internal/platform/server"
	"boletia/api/internal/platform/storage/postgres"
	cache "boletia/api/internal/platform/storage/redis"
	"boletia/api/kit/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"log"
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + fmt.Sprint(cfg.RedisPort),
		Password: "",
		DB:       cfg.RedisDb,
	})

	loggerApp, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	loggerApp = loggerApp.With("service", "boletia")

	databaseRepository := postgres.NewDatabaseRepository(dbpool, cfg.DbTimeout, cfg.DbTableName, cfg.DbWordAllCurrency)
	cacheRepository := cache.NewCacheRepository(rdb, cfg.RedisExpirationTime, cfg.RedisNamespace)
	getCurrencyService := obtain.NewCurrencyService(databaseRepository, cacheRepository)

	ctx, srv := server.New(context.Background(), cfg.Host, cfg.Port, cfg.ShutdownTimeout, getCurrencyService)
	return srv.Run(ctx)
}

type config struct {
	// Server configuration
	Host            string        `default:"localhost"`
	Port            uint          `default:"8080"`
	ShutdownTimeout time.Duration `default:"10s"`
	// Database configuration
	DbUser            string        `default:"altorik"`
	DbPass            string        `default:"superpass"`
	DbHost            string        `default:"db"`
	DbName            string        `default:"boletia"`
	DbPort            uint          `default:"5432"`
	DbTimeout         time.Duration `default:"6s"`
	DbTableName       string        `default:"temp_currency_rates"`
	DbWordAllCurrency string        `default:"ALLC"`
	// Cache config
	RedisHost           string        `default:"redis"`
	RedisPort           uint          `default:"6379"`
	RedisNamespace      string        `default:"BOL"`
	RedisExpirationTime time.Duration `default:"30s"`
	RedisDb             int           `default:"0"`
}

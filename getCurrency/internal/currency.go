package bole

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type CurrencyData struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]Currency `json:"data"`
	ID   string              `json:"id"`
}

// Currency is the data structure that represents a currency.
type Currency struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

// CurrencyRepository defines the expected behaviour from a currency api.
type CurrencyRepository interface {
	Get(ctx context.Context) (*CurrencyData, ApiCall, error)
}

//go:generate mockery --case=snake --outpkg=apimocks --output=platform/http/apimocks --name=CurrencyRepository

// DatabaseRepository defines the expected behaviour from a currency api.
type DatabaseRepository interface {
	Save(ctx context.Context, data *CurrencyData) error
	SaveCall(ctx context.Context, data ApiCall) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=DatabaseRepository

// DBIface defines the expected behaviour from a currency api.
type DBIface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=DBIface

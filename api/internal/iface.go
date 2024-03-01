package bole

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// IDatabaseRepository defines the expected behaviour from a currency api.
type IDatabaseRepository interface {
	Get(ctx context.Context, criteria Criteria) ([]Currency, error)
	GetCode(ctx context.Context, criteria Criteria) (bool, error)
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=IDatabaseRepository

// DBIface defines the expected behaviour from a currency api.
type DBIface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=DBIface

// ICacheRepository defines the expected behaviour from a cache api.
type ICacheRepository interface {
	Get(ctx context.Context, hash string) ([]Currency, error)
	Set(ctx context.Context, hash string, data []Currency) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=ICacheRepository

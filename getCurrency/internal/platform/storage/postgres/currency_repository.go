package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"

	bole "boletia/internal"
)

// DatabaseRepository is a Postgres bole.DatabaseRepository implementation.
type DatabaseRepository struct {
	db        bole.DBIface
	dbTimeout time.Duration
	tableName string
}

// NewDatabaseRepository initializes a Postgres-based implementation of bole.DatabaseRepository.
func NewDatabaseRepository(db bole.DBIface, dbTimeout time.Duration, tableName string) *DatabaseRepository {
	return &DatabaseRepository{
		db:        db,
		dbTimeout: dbTimeout,
		tableName: tableName,
	}
}

// Save implements the bole.DatabaseRepository interface.
func (r *DatabaseRepository) Save(ctx context.Context, data *bole.CurrencyData) error {
	if data.ID == "" {
		return errors.New("invalid batch id")
	}
	if data.Meta.LastUpdatedAt == (time.Time{}) {
		return errors.New("invalid last updated at time")
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "CREATE TEMP TABLE IF NOT EXISTS temp_currency_rates AS SELECT * FROM currency_rates WITH NO DATA;")
	if err != nil {
		return err
	}
	rows := make([][]interface{}, len(data.Data))
	var index = 0
	for _, obj := range data.Data {
		rows[index] = []interface{}{obj.Code, obj.Value, data.Meta.LastUpdatedAt, data.ID}
		index++
	}

	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{r.tableName},
		[]string{"code", "value", "last_updated_at", "batch_id"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		if err = tx.Rollback(ctx); err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, `
	   INSERT INTO currency_rates (code, value, last_updated_at, batch_id)
	   SELECT code, value, last_updated_at, batch_id
	   FROM temp_currency_rates
	   ON CONFLICT (code, last_updated_at) DO NOTHING;
	`)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `TRUNCATE TABLE temp_currency_rates;`)
	if err != nil {
		return err
	}

	return nil
}

// SaveCall implements the bole.DatabaseRepository interface.
func (r *DatabaseRepository) SaveCall(ctx context.Context, data bole.ApiCall) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO api_calls (id, status_code, response_time, timeout, error_message, called_at)
		VALUES ($1, $2, $3, $4, $5,NOW())`,
		data.ID, data.StatusCode, data.ResponseTime, data.Timeout, data.ErrorMessage,
	)
	if err != nil {
		return err
	}
	return nil
}

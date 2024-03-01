package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	bole "boletia/api/internal"
)

// DatabaseRepository is a Postgres bole.IDatabaseRepository implementation.
type DatabaseRepository struct {
	db              bole.DBIface
	dbTimeout       time.Duration
	tableName       string
	allCurrencyWord string
}

// NewDatabaseRepository initializes a Postgres-based implementation of bole.IDatabaseRepository.
func NewDatabaseRepository(db bole.DBIface, dbTimeout time.Duration, tableName, allCurrencyWord string) *DatabaseRepository {
	return &DatabaseRepository{
		db:              db,
		dbTimeout:       dbTimeout,
		tableName:       tableName,
		allCurrencyWord: allCurrencyWord,
	}
}

// Save implements the bole.IDatabaseRepository interface.
func (r *DatabaseRepository) Get(ctx context.Context, criteria bole.Criteria) ([]bole.Currency, error) {
	queryParams, query := r.getParams(criteria)
	query += " ORDER BY last_updated_at LIMIT 300"
	rows, err := r.db.Query(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []bole.Currency
	for rows.Next() {
		var item bole.Currency
		if err = rows.Scan(&item.Code, &item.Value, &item.LastModifiedAt); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if len(results) == 0 {
		return nil, bole.ErrorNotFound
	}
	return results, nil
}

// GetCode implements the bole.IDatabaseRepository interface.
func (r *DatabaseRepository) GetCode(ctx context.Context, criteria bole.Criteria) (bool, error) {
	if r.allCurrencyWord == criteria.CurrencyCode {
		return true, nil
	}
	query := "SELECT 1 FROM currency_rates WHERE code = $1 LIMIT 1"

	var exists int
	err := r.db.QueryRow(ctx, query, criteria.CurrencyCode).Scan(&exists)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *DatabaseRepository) getParams(criteria bole.Criteria) ([]interface{}, string) {
	var queryParams []interface{}
	query := "SELECT code, value, last_updated_at FROM currency_rates"

	paramCounter := 1
	if criteria.CurrencyCode != r.allCurrencyWord {
		query += fmt.Sprintf(" WHERE code = $%d", paramCounter)
		queryParams = append(queryParams, criteria.CurrencyCode)
		paramCounter++
	}

	if !criteria.StartDate.IsZero() {
		if paramCounter == 1 {
			query += " WHERE"
		} else {
			query += " AND"
		}
		query += fmt.Sprintf(" last_updated_at >= $%d", paramCounter)
		queryParams = append(queryParams, criteria.StartDate)
		paramCounter++
	}
	if !criteria.EndDate.IsZero() {
		if paramCounter == 1 {
			query += " WHERE"
		} else {
			query += " AND"
		}
		query += fmt.Sprintf(" last_updated_at <= $%d", paramCounter)
		queryParams = append(queryParams, criteria.EndDate)
	}
	return queryParams, query
}

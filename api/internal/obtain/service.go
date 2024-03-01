package obtain

import (
	bol "boletia/api/internal"
	"context"
)

type ICurrencyService interface {
	GetCurrency(ctx context.Context, criteria bol.Criteria) ([]bol.Currency, error)
}

//go:generate mockery --case=snake --outpkg=servicemocks --output=mock --name=ICurrencyService

// CurrencyService is the default CurrencyService interface
// implementation returned by obtain.NewCurrencyService.
type CurrencyService struct {
	database bol.IDatabaseRepository
	cache    bol.ICacheRepository
}

// NewCurrencyService returns the default Service interface implementation.
func NewCurrencyService(database bol.IDatabaseRepository, cache bol.ICacheRepository) CurrencyService {
	return CurrencyService{
		database: database,
		cache:    cache,
	}
}

// GetCurrency implements to obtain.CurrencyService interface.
func (s CurrencyService) GetCurrency(ctx context.Context, criteria bol.Criteria) ([]bol.Currency, error) {
	exist, err := s.database.GetCode(ctx, criteria)
	if err != nil {
		return []bol.Currency{}, err
	}
	if !exist {
		return []bol.Currency{}, bol.ErrorCurrencyCode
	}

	data, err := s.cache.Get(ctx, criteria.Hash())
	if err != nil {
		return []bol.Currency{}, err
	}
	if len(data) > 0 {
		return data, nil
	}

	data, err = s.database.Get(ctx, criteria)
	if err != nil {
		return []bol.Currency{}, err
	}
	if err = s.cache.Set(ctx, criteria.Hash(), data); err != nil {
		return []bol.Currency{}, err
	}

	return data, nil
}

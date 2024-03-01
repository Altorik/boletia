package obtain

import (
	bol "boletia/internal"
	"context"
)

type ICurrencyService interface {
	GetCurrency(ctx context.Context) error
}

//go:generate mockery --case=snake --outpkg=servicemocks --output=mock --name=ICurrencyService

// CurrencyService is the default CurrencyService interface
// implementation returned by obtain.NewCurrencyService.
type CurrencyService struct {
	currencyRepository bol.CurrencyRepository
	database           bol.DatabaseRepository
}

// NewCurrencyService returns the default Service interface implementation.
func NewCurrencyService(currencyRepository bol.CurrencyRepository, database bol.DatabaseRepository) CurrencyService {
	return CurrencyService{
		currencyRepository: currencyRepository,
		database:           database,
	}
}

// GetCurrency implements to obtain.CurrencyService interface.
func (s CurrencyService) GetCurrency(ctx context.Context) error {
	currency, apicall, err := s.currencyRepository.Get(ctx)
	if err = s.database.SaveCall(ctx, apicall); err != nil {
		return err
	}
	if len(currency.Data) < 1 {
		return nil
	}
	if err = s.database.Save(ctx, currency); err != nil {
		return err
	}

	return nil
}

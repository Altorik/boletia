package currency

import (
	"boletia/internal/obtain"
	"context"
)

// ObtainHandler returns an error or not for currency creation.
func ObtainHandler(ctx context.Context, getCurrency obtain.ICurrencyService) error {
	err := getCurrency.GetCurrency(ctx)
	if err != nil {
		return err
	}
	return nil
}

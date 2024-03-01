package obtain

import (
	bol "boletia/internal"
	"boletia/internal/platform/http/apimocks"
	"boletia/internal/platform/storage/storagemocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_CurrencyService_Get_RepositoryError(t *testing.T) {
	currencyRepositoryMock := new(apimocks.CurrencyRepository)
	currencyRepositoryMock.On("Get", mock.Anything).Return(bol.CurrencyData{}, errors.New("something unexpected happened"))

	databaseRepositoryMock := new(storagemocks.DatabaseRepository)

	currencyService := NewCurrencyService(currencyRepositoryMock, databaseRepositoryMock)

	err := currencyService.GetCurrency(context.Background())

	currencyRepositoryMock.AssertExpectations(t)
	databaseRepositoryMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_CurrencyService_Get_Succeed(t *testing.T) {
	data := map[string]bol.Currency{
		"USD": {Code: "USD", Value: 123.12},
		"EUR": {Code: "EUR", Value: 12.12},
	}
	currData := bol.CurrencyData{Data: data}
	currencyRepositoryMock := new(apimocks.CurrencyRepository)
	currencyRepositoryMock.On("Get", mock.Anything).Return(currData, nil)

	databaseRepositoryMock := new(storagemocks.DatabaseRepository)
	databaseRepositoryMock.On("Save", mock.Anything, mock.AnythingOfType("bole.CurrencyData")).Return(nil)

	currencyService := NewCurrencyService(currencyRepositoryMock, databaseRepositoryMock)

	err := currencyService.GetCurrency(context.Background())

	currencyRepositoryMock.AssertExpectations(t)
	databaseRepositoryMock.AssertExpectations(t)

	assert.NoError(t, err)
}

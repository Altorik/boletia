package obtain

import (
	bol "boletia/api/internal"
	"boletia/api/internal/platform/storage/storagemocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_CurrencyService_Get_RepositoryError(t *testing.T) {
	criteria := bol.Criteria{
		CurrencyCode: "MXN",
		StartDate:    time.Now(),
		EndDate:      time.Now(),
	}
	cacheRepositoryMock := new(storagemocks.ICacheRepository)
	cacheRepositoryMock.On("Get", mock.Anything, mock.AnythingOfType("string")).Return([]bol.Currency{}, errors.New("something unexpected happened"))

	databaseRepositoryMock := new(storagemocks.IDatabaseRepository)

	currencyService := NewCurrencyService(databaseRepositoryMock, cacheRepositoryMock)

	_, err := currencyService.GetCurrency(context.Background(), criteria)

	cacheRepositoryMock.AssertExpectations(t)
	databaseRepositoryMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_CurrencyService_Get_Cache_Succeed(t *testing.T) {
	criteria := bol.Criteria{
		CurrencyCode: "MXN",
		StartDate:    time.Now(),
		EndDate:      time.Now(),
	}
	mockData := []bol.Currency{
		{
			Code:           "MXN",
			Value:          120,
			LastModifiedAt: time.Now(),
		}, {
			Code:           "MXN",
			Value:          10,
			LastModifiedAt: time.Now(),
		},
	}

	cacheRepositoryMock := new(storagemocks.ICacheRepository)
	cacheRepositoryMock.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(mockData, nil)

	databaseRepositoryMock := new(storagemocks.IDatabaseRepository)

	currencyService := NewCurrencyService(databaseRepositoryMock, cacheRepositoryMock)

	resData, err := currencyService.GetCurrency(context.Background(), criteria)
	require.NoError(t, err)
	cacheRepositoryMock.AssertExpectations(t)
	databaseRepositoryMock.AssertExpectations(t)
	assert.Equal(t, resData, mockData)
}

func Test_CurrencyService_Get_Succeed(t *testing.T) {
	ctx := context.TODO()
	criteria := bol.Criteria{
		CurrencyCode: "MXN",
		StartDate:    time.Now(),
		EndDate:      time.Now(),
	}
	mockData := []bol.Currency{
		{
			Code:           "MXN",
			Value:          120,
			LastModifiedAt: time.Now(),
		}, {
			Code:           "MXN",
			Value:          10,
			LastModifiedAt: time.Now(),
		},
	}

	databaseRepositoryMock := new(storagemocks.IDatabaseRepository)
	cacheRepositoryMock := new(storagemocks.ICacheRepository)
	hash := criteria.Hash()
	cacheRepositoryMock.On("Get", ctx, hash).Return([]bol.Currency{}, nil) // Primera llamada devuelve caché vacío
	databaseRepositoryMock.On("Get", ctx, criteria).Return(mockData, nil)  // Base de datos devuelve datos
	cacheRepositoryMock.On("Set", ctx, hash, mockData).Return(nil)

	currencyService := NewCurrencyService(databaseRepositoryMock, cacheRepositoryMock)

	resData, err := currencyService.GetCurrency(ctx, criteria)
	require.NoError(t, err)
	cacheRepositoryMock.AssertExpectations(t)
	databaseRepositoryMock.AssertExpectations(t)
	assert.Equal(t, resData, mockData)
}

package currency

import (
	bol "boletia/api/internal"
	servicemocks "boletia/api/internal/obtain/mock"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Currency_ServiceError(t *testing.T) {
	criteria := bol.Criteria{
		CurrencyCode: "MXN",
		StartDate:    time.Now(),
		EndDate:      time.Now(),
	}
	currencyService := new(servicemocks.ICurrencyService)
	currencyService.On("GetCurrency", mock.Anything, criteria).Return(errors.New("something unexpected happened"))
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/currencies/:currency", ObtainHandler(currencyService))

	t.Run("given an invalid fend it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/currencies/mxn?finit=2024-02-28T23:50:56&fend=2024-03-29 23:50:59", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("given an invalid finit it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/currencies/mxn?finit=2024-02-28T23:50:66&fend=2024-03-29T23:50:59", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("given an invalid dates it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/currencies/mxn?finit=2024-02-28T23:50:56&fend=2024-02-27T23:50:59", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestHandler_Currency_ServiceSuccess(t *testing.T) {
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
	currencyService := new(servicemocks.ICurrencyService)
	currencyService.On("GetCurrency", mock.AnythingOfType("*gin.Context"), mock.Anything).Return(mockData, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/currencies/:currency", ObtainHandler(currencyService))

	t.Run("only currency code and return 201", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/currencies/mxn", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("currency code finit and return 201", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/currencies/mxn?finit=2024-02-28T23:50:56", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("currency code fend and return 201", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/currencies/mxn?fend=2024-03-29T23:50:59", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

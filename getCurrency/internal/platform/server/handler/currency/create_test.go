package currency

import (
	servicemocks "boletia/internal/obtain/mock"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Currency_ServiceError(t *testing.T) {
	currencyService := new(servicemocks.ICurrencyService)
	currencyService.On("GetCurrency", mock.Anything).Return(errors.New("something unexpected happened"))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/getData", ObtainHandler(currencyService))

	t.Run("given an invalid request it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/getData", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestHandler_Currency_ServiceSuccess(t *testing.T) {
	currencyService := new(servicemocks.ICurrencyService)
	currencyService.On("GetCurrency", mock.Anything).Return(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/getData", ObtainHandler(currencyService))

	t.Run("return 201", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/getData", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}

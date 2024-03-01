package api

import (
	bol "boletia/internal"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurrencyRepository_Get(t *testing.T) {
	currencies := map[string]bol.Currency{
		"USD": {Code: "USD", Value: 1.0},
		"EUR": {Code: "EUR", Value: 0.9},
	}
	currData := &bol.CurrencyData{Data: currencies}
	currenciesJSON, err := json.Marshal(currData)
	require.NoError(t, err)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(currenciesJSON))
	}))
	defer testServer.Close()
	logger := zaptest.NewLogger(t).Sugar()
	repo := NewCurrencyRepository(testServer.Client(), logger, testServer.URL, "")

	result, apicall, err := repo.Get(context.Background())
	require.NoError(t, err)

	if len(result.Data) != len(currencies) {
		t.Errorf("Expected %d currencies, got %d", len(currencies), len(result.Data))
	}
	if apicall.ID != result.ID {
		t.Errorf("Expected apicall id %v, got %v", result.ID, apicall.ID)
	}

	for i, currency := range result.Data {
		if currency.Code != currencies[i].Code || currency.Value != currencies[i].Value {
			t.Errorf("Expected currency %v, got %v", currencies[i], currency)
		}
	}
}

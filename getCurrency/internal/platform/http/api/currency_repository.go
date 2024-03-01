package api

import (
	bol "boletia/internal"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"time"
)

// CurrencyRepository is the repository for the call api for get currencies.
type CurrencyRepository struct {
	client *http.Client
	logger *zap.SugaredLogger
	url    string
	key    string
}

// NewCurrencyRepository initializes a CurrencyRepository.
func NewCurrencyRepository(client *http.Client, logger *zap.SugaredLogger, url, key string) *CurrencyRepository {
	return &CurrencyRepository{
		client: client,
		logger: logger,
		url:    url,
		key:    key,
	}
}

// Get invoke to get data from currencyapi.
func (m *CurrencyRepository) Get(ctx context.Context) (*bol.CurrencyData, bol.ApiCall, error) {
	var apiCall bol.ApiCall
	apiCall.ID = uuid.New().String()
	start := time.Now()
	resp, err := m.client.Get(m.url + m.key)
	if err != nil {
		apiCall.ErrorMessage = err.Error()
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			apiCall.Timeout = true
		}
		return &bol.CurrencyData{}, apiCall, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		apiCall.ErrorMessage = err.Error()
		return &bol.CurrencyData{}, apiCall, err
	}
	elapsed := time.Since(start)
	apiCall.StatusCode = resp.StatusCode
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		apiCall.ErrorMessage = err.Error()
		return &bol.CurrencyData{}, apiCall, errors.New("error get data from currencyapi endpoint: " + string(rawData))
	}
	var data bol.CurrencyData
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		apiCall.ErrorMessage = err.Error()
		return &bol.CurrencyData{}, apiCall, errors.New("error parse response from currencyapi info endpoint: " + err.Error())
	}
	data.ID = apiCall.ID
	apiCall.ResponseTime = elapsed
	return &data, apiCall, nil
}

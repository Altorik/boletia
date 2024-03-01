package bole

import (
	"errors"
	"time"
)

// Currency is the data structure that represents a currency.
type Currency struct {
	Code           string    `json:"code"`
	Value          float64   `json:"value"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

var ErrorNotFound = errors.New("not value found")

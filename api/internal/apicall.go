package bole

import "time"

type ApiCall struct {
	ID           string `json:"id"`
	ResponseTime time.Duration
	StatusCode   int
	Timeout      bool
	ErrorMessage string
}

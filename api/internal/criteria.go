package bole

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"
)

type Criteria struct {
	CurrencyCode       string
	StartDate, EndDate time.Time
}

var ErrorValidation = errors.New("invalid date format")
var ErrorDates = errors.New("invalid finit and fend dates")
var ErrorCurrencyCode = errors.New("invalid currency code")

func (d Criteria) NewCriteria(code, fechaInicio, fechaFin string) (Criteria, error) {
	fInit, fEnd := time.Time{}, time.Time{}
	var err error
	if fechaInicio != "" {
		fInit, err = time.Parse("2006-01-02T15:04:05", fechaInicio)
		if err != nil {
			return Criteria{}, ErrorValidation
		}
	}
	if fechaFin != "" {
		fEnd, err = time.Parse("2006-01-02T15:04:05", fechaFin)
		if err != nil {
			return Criteria{}, ErrorValidation
		}
		if fInit.After(fEnd) {
			return Criteria{}, ErrorDates
		}
	}
	code = strings.ToUpper(code)
	if code == "" || len(code) < 3 || len(code) > 6 || !d.isAlpha(code) {
		return Criteria{}, ErrorCurrencyCode
	}
	return Criteria{
		CurrencyCode: code,
		StartDate:    fInit,
		EndDate:      fEnd,
	}, nil
}

func (d Criteria) Hash() string {
	hash := sha256.Sum256([]byte(d.CurrencyCode + ":" + d.StartDate.String() + ":" + d.EndDate.String()))
	return hex.EncodeToString(hash[:])
}

func (d Criteria) isAlpha(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z]+$`)
	return re.MatchString(s)
}

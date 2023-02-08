package models

import (
	"checkout-task/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AuthRequest struct {
	Email string `json:"email"`
}

func (a AuthRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}

type RefreshRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (a RefreshRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Token,
			validation.Required,
			validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
		),
	)
}

type PaymentResponse struct {
	PaymentID        string  `json:"payment_id"`
	MaskedCardNumber string  `json:"masked_card_number"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"`
}

type PaymentRequest struct {
	CardNumber string  `json:"card_number"`
	ExpMonth   int     `json:"expiry_month"`
	ExpYear    int     `json:"expiry_year"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	CVV        int     `json:"cvv"`
}

type ProcessPaymentResponse struct {
	PaymentIdentifier string
	Status            string
}

func ValidatePaymentRequest(request PaymentRequest) error {
	// Validate card number format
	if len(strings.ReplaceAll(request.CardNumber, " ", "")) != 16 {
		return errors.New("invalid card number")
	}

	// Validate expiry date
	now := time.Now()
	expDate := time.Date(request.ExpYear, time.Month(request.ExpMonth), 1, 0, 0, 0, 0, time.UTC)
	if expDate.Before(now) {
		return errors.New("expired card")
	}

	// Validate amount
	if request.Amount <= 0 || request.Amount > 1000 {
		return errors.New("invalid amount")
	}

	// Validate currency
	supportedCurrencies := []string{"USD", "EUR", "GBP"}
	if !utils.Contains(supportedCurrencies, request.Currency) {
		return errors.New("unsupported currency")
	}

	// Validate CVV
	if len(strconv.Itoa(request.CVV)) != 3 && len(strconv.Itoa(request.CVV)) != 4 {
		return errors.New("invalid CVV")
	}

	return nil
}

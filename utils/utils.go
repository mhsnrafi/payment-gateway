package utils

import (
	"strings"
	"time"
)

func Contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func ValidateCardNumber(cardNumber string) bool {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	if len(cardNumber) != 16 {
		return false
	}

	for i := 0; i < 16; i++ {
		if cardNumber[i] < '0' || cardNumber[i] > '9' {
			return false
		}
	}

	return true
}

func ValidateExpiryDate(expMonth int, expYear int) bool {
	now := time.Now()
	if expYear < now.Year() || (expYear == now.Year() && expMonth < int(now.Month())) {
		return false
	}

	return true
}

func ValidateAmount(amount float64) bool {
	if amount <= 0 {
		return false
	}

	return true
}

func ValidateCurrency(currency string) bool {
	supportedCurrencies := []string{"USD", "EUR", "GBP"}
	for _, c := range supportedCurrencies {
		if c == currency {
			return true
		}
	}

	return false
}

func ValidateCVV(cvv int) bool {
	if cvv < 100 || cvv > 999 {
		return false
	}

	return true
}

// maskCardNumber function
func MaskCardNumber(cardNumber string) string {
	if len(cardNumber) > 4 {
		return "**** **** **** " + cardNumber[len(cardNumber)-4:]
	}
	return cardNumber
}

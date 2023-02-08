package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContains(t *testing.T) {
	list := []string{"USD", "EU", "CAD"}

	// Test 1: Test the function with a valid item
	item := "EU"
	result := Contains(list, item)
	assert.Equal(t, true, result)

	item = "USD"
	result = Contains(list, item)
	assert.Equal(t, true, result)

	// Test 2: Test the function with an invalid item
	item = "test currency"
	result = Contains(list, item)
	assert.Equal(t, false, result)
}

func TestValidateCardNumber(t *testing.T) {
	// Test case 1: valid card number
	cardNumber := "2345 6789 0123 4567"
	result := ValidateCardNumber(cardNumber)

	assert.Equal(t, true, result)

	// Test case 2: invalid length of card number
	cardNumber = "23456789012345"
	result = ValidateCardNumber(cardNumber)
	assert.Equal(t, false, result)

	// Test case 3: invalid characters in card number
	cardNumber = "1234 5678 9012 3A6F"
	result = ValidateCardNumber(cardNumber)
	assert.Equal(t, false, result)
}

func TestValidateExpiryDate(t *testing.T) {
	// Test case 1: valid expiry date
	expMonth := 12
	expYear := 2030
	result := ValidateExpiryDate(expMonth, expYear)
	assert.Equal(t, true, result)

	// Test case 2: expiry date in the past
	expMonth = 2
	expYear = 2022
	result = ValidateExpiryDate(expMonth, expYear)
	assert.Equal(t, false, result)

	// Test case 3: current month and year
	now := time.Now()
	expMonth = int(now.Month())
	expYear = now.Year()
	result = ValidateExpiryDate(expMonth, expYear)
	assert.Equal(t, true, result)
}

func TestValidateAmount(t *testing.T) {
	// Test case 1: valid amount
	amount := 100.0
	result := ValidateAmount(amount)
	assert.Equal(t, true, result)

	// Test case 2: invalid amount (0)
	amount = 0.0
	result = ValidateAmount(amount)
	assert.Equal(t, false, result)

	// Test case 3: invalid amount (-100)
	amount = -100.0
	result = ValidateAmount(amount)
	assert.Equal(t, false, result)
}

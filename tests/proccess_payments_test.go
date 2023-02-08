package main

import (
	db "checkout-task/models/db"
	"checkout-task/services"

	"github.com/stretchr/testify/assert"
	"testing"

	"checkout-task/models"
)

func TestProcessPayment_Success(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	payment := models.PaymentRequest{
		CardNumber: "1234123412341234",
		ExpMonth:   12,
		ExpYear:    2024,
		CVV:        123,
		Amount:     10.0,
		Currency:   "USD",
	}

	response, err := services.ProcessPayment(payment)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "failure", response.Status)

	payment = models.PaymentRequest{
		CardNumber: "1234123413341234",
		ExpMonth:   12,
		ExpYear:    2024,
		CVV:        123,
		Amount:     10.0,
		Currency:   "USD",
	}

	response, err = services.ProcessPayment(payment)

	assert.NoError(t, err)
	assert.Equal(t, "failure", response.Status)

	// Check if payment details are stored in the database correctly
	var paymentInfo db.Payment
	err = services.DbConnection.Find(&paymentInfo).Where("card_number = ?", "1234123412341234").Error
	assert.NoError(t, err)
	assert.Equal(t, payment.CardNumber, paymentInfo.CardNumber)
	assert.Equal(t, payment.ExpMonth, paymentInfo.ExpMonth)
	assert.Equal(t, payment.ExpYear, paymentInfo.ExpYear)
	assert.Equal(t, payment.CVV, paymentInfo.CVV)
	assert.Equal(t, payment.Amount, paymentInfo.Amount)
	assert.Equal(t, payment.Currency, paymentInfo.Currency)
	assert.Equal(t, "In-Progress", paymentInfo.Status)
}

func TestRetrievePayment(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	paymentID := "9bbaf0c4-8d6f-4ae1-920a-b91299b68b0a"
	paymentInfo, err := services.RetrievePayment(paymentID)
	assert.NoError(t, err)
	assert.Equal(t, paymentID, paymentInfo.PaymentID)
	assert.Equal(t, "**** **** **** 3457", paymentInfo.MaskedCardNumber)
	assert.Equal(t, 34.00, paymentInfo.Amount)
	assert.Equal(t, "EUR", paymentInfo.Currency)
	assert.Equal(t, "Success", paymentInfo.Status)
}

func TestRetrievePaymentWithInvalidID(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	paymentID := "-1"
	_, err := services.RetrievePayment(paymentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve payment details")
}

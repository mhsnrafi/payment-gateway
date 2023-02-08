package services

import (
	"checkout-task/constants"
	"checkout-task/logger"
	"checkout-task/models"
	db "checkout-task/models/db"
	"checkout-task/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type ProcessPaymentResponse struct {
	PaymentIdentifier string
	Status            string
}

// ProcessPayment processes a payment and returns a response
func ProcessPayment(payment models.PaymentRequest) (ProcessPaymentResponse, error) {
	// Store payment details in database
	paymentInfo := db.Payment{
		PaymentID:  uuid.New().String(),
		CardNumber: payment.CardNumber,
		ExpMonth:   payment.ExpMonth,
		ExpYear:    payment.ExpYear,
		CVV:        payment.CVV,
		Amount:     payment.Amount,
		Currency:   payment.Currency,
		Status:     constants.INPROGRESS,
	}
	if err := DbConnection.Create(&paymentInfo).Error; err != nil {
		logger.Error("failed to store payment details", zap.Error(err))
		return ProcessPaymentResponse{}, fmt.Errorf("failed to store payment details")
	}

	// Simulate acquiring bank
	status := simulateAcquiringBank(payment)

	if status == constants.SUCCESS {
		if err := DbConnection.Model(&db.Payment{}).Where("id = ?", paymentInfo.ID).Update(db.Payment{Status: constants.SUCCESS}).Error; err != nil {
			logger.Error("failed to store payment details", zap.Error(err))
			return ProcessPaymentResponse{}, fmt.Errorf("failed to store payment details")
		}
	}

	logger.Info("Payment processed", zap.String("status", status))
	//  Return success or failure response
	return ProcessPaymentResponse{
		Status:            status,
		PaymentIdentifier: paymentInfo.PaymentID,
	}, nil

}

func RetrievePayment(paymentID string) (*models.PaymentResponse, error) {
	var payment db.Payment
	if err := DbConnection.Where("payment_id = ?", paymentID).First(&payment).Error; err != nil {
		logger.Error("failed to retrieve payment details", zap.String("error", err.Error()), zap.String("payment_id", paymentID))
		return &models.PaymentResponse{}, fmt.Errorf("failed to retrieve payment details")
	}

	// Get payment details from cache
	paymentDetails, isPaymentDetailsFound, err := GetPaymentDetailsFromCache(paymentID)
	if err != nil {
		logger.Error("failed to retrieve payment details from cache", zap.String("error", err.Error()), zap.String("payment_id", paymentID))
	}

	if isPaymentDetailsFound {
		return paymentDetails, nil
	}

	// Mask card number
	maskedCardNumber := utils.MaskCardNumber(payment.CardNumber)
	payment.CardNumber = maskedCardNumber

	// Return payment details
	paymentDetails = &models.PaymentResponse{
		PaymentID:        payment.PaymentID,
		MaskedCardNumber: payment.CardNumber,
		Amount:           payment.Amount,
		Currency:         payment.Currency,
		Status:           payment.Status,
	}

	// Set payment details in cache
	SetPaymentDetails(paymentID, paymentDetails)
	return paymentDetails, nil
}

func simulateAcquiringBank(payment models.PaymentRequest) string {
	// Validate card details
	validCardNumber := utils.ValidateCardNumber(payment.CardNumber)
	validExpiryDate := utils.ValidateExpiryDate(payment.ExpMonth, payment.ExpYear)
	validAmount := utils.ValidateAmount(payment.Amount)
	validCurrency := utils.ValidateCurrency(payment.Currency)
	validCVV := utils.ValidateCVV(payment.CVV)

	if !validCardNumber || !validExpiryDate || !validAmount || !validCurrency || !validCVV {
		logger.Error("Payment failed due to invalid card details")
		return constants.FAILURE
	}

	// Check if card has sufficient funds
	hasSufficientFunds := CheckFunds(payment.CardNumber, payment.Amount)
	if !hasSufficientFunds {
		logger.Error("Payment failed due to insufficient funds")
		return constants.FAILURE
	}

	// Check for fraud
	isFraud := CheckForFraud(payment.CardNumber, payment.Amount, payment.Currency)
	if isFraud {
		logger.Error("Payment failed due to potential fraud")
		return constants.FAILURE
	}

	// Approve payment
	err := ApprovePayment(payment.CardNumber, payment.Amount)
	if err != nil {
		logger.Error("Payment failed while approving payment", zap.Error(err))
		return constants.FAILURE
	}

	logger.Info("Payment approved successfully")
	// Return status
	return constants.SUCCESS
}

func CheckFunds(cardNumber string, amount float64) bool {
	var card db.Card
	if err := DbConnection.Where("card_number = ?", cardNumber).First(&card).Error; err != nil {
		fmt.Println("Error finding card:", err)
		return false
	}

	return card.Balance >= amount
}

func CheckForFraud(cardNumber string, amount float64, currency string) bool {
	var fraud db.Fraud
	if err := DbConnection.Where("card_number = ? AND amount >= ? AND currency = ?", cardNumber, amount, currency).First(&fraud).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		fmt.Println("Error finding fraud:", err)
		return true
	}

	return true
}

func ApprovePayment(cardNumber string, amount float64) error {
	// Find the card in the database
	var card db.Card
	if err := DbConnection.Where("card_number = ?", cardNumber).First(&card).Error; err != nil {
		return err
	}

	// Update the card balance
	card.Balance -= amount
	if err := DbConnection.Save(&card).Error; err != nil {
		return err
	}

	return nil
}

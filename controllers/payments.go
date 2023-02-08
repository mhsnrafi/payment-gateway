package controllers

import (
	"checkout-task/constants"
	"checkout-task/models"
	"checkout-task/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sony/gobreaker"
	"net/http"
)

// create a circuit breaker to handle failures in the payment processing service
var cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{})

// ProcessPayment handles the processing of a payment request.
// @Summary Handle the processing of a payment request.
// @Description Handle the processing of a payment request by validating the payment details, processing the payment and returning the response.
// @Tags Payments
// @Accept json
// @Produce json
// @Param paymentReq body models.PaymentRequest true "Payment Request"
// @Success 201 {object} models.Response
// @Success 400 {object} models.Response
// @Router /payment [post]
func ProcessPayment(c *gin.Context) {
	var paymentReq models.PaymentRequest
	_ = c.ShouldBindBodyWith(&paymentReq, binding.JSON)

	// Validate payment details
	// Validate payment details
	if err := models.ValidatePaymentRequest(paymentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// process the payment using the circuit breaker
	paymentResponse, err := cb.Execute(func() (interface{}, error) {
		return services.ProcessPayment(paymentReq)
	})
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// cast the payment response to the expected type
	pr, ok := paymentResponse.(models.ProcessPaymentResponse)
	if !ok {
		response.Message = "Error processing payment"
		response.SendResponse(c)
		return
	}

	// // Return success response
	if pr.Status == constants.SUCCESS {
		response.Success = true
		response.StatusCode = http.StatusCreated
		response.Data = gin.H{
			"Payment Status": pr.Status,
			"Payment ID":     pr.PaymentIdentifier,
			"Message":        "Payment processed successfully",
		}
		response.SendResponse(c)
	} else {
		response.StatusCode = http.StatusBadRequest
		response.Success = false
		response.Data = gin.H{
			"Status":  pr.Status,
			"Message": "Payment failed",
		}
		response.SendResponse(c)
	}

}

// GetPaymentDetails retrieves payment details for a specified payment.
// @Summary Retrieve payment details for a specified payment.
// @Description Retrieve payment details for a specified payment using the payment ID provided in the request URL query parameters.
// @Tags Payments
// @Accept  json
// @Produce  json
// @Param payment_id query string true "Payment ID"
// @Success 200 {object} models.Response
// @Success 400 {object} models.Response
// @Router /payment [get]
func GetPaymentDetails(c *gin.Context) {
	paymentID := c.Request.URL.Query().Get("payment_id")

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// retrieve the previous made payment details
	paymentDetails, err := cb.Execute(func() (interface{}, error) {
		return services.RetrievePayment(paymentID)
	})
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// cast the payment response to the expected type
	paymentsInfo, ok := paymentDetails.(models.PaymentResponse)
	if !ok {
		response.Message = "Error processing payment"
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"Payment Details": paymentsInfo,
	}
	response.SendResponse(c)
}

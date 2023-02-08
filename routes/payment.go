package routes

import (
	"checkout-task/controllers"
	"checkout-task/middlewares"
	"github.com/gin-gonic/gin"
)

func Payments(router *gin.RouterGroup) {
	auth := router.Group("/")
	{
		auth.POST(
			"/process-payment",
			middlewares.JWTMiddleware(),
			controllers.ProcessPayment,
		)
		auth.GET(
			"/get-payment",
			middlewares.JWTMiddleware(),
			controllers.GetPaymentDetails,
		)
	}
}

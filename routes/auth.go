package routes

import (
	"checkout-task/controllers"
	"checkout-task/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST(
			"/generate_access_token",
			validators.AuthValidator(),
			controllers.GenerateAccessToken,
		)

		auth.POST(
			"/refresh",
			validators.RefreshValidator(),
			controllers.Refresh,
		)
	}
}

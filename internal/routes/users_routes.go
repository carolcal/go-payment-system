package routes

import (
	"qr-payment/internal/handlers"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func SetUpUserRoutes(router *gin.Engine, uhandler handlers.UserHandlers) {

	users := router.Group("/users")
	{
		users.GET("", uhandler.GetAllUsersHandler)
	}

	user := router.Group("/user")
	{
		user.POST("", uhandler.CreateUserHandler)
		user.GET("/:id", uhandler.GetUserByIdHandler)
		user.PUT("/:id/balance", uhandler.UpdateBalanceHandler)
		user.DELETE("/:id", uhandler.RemoveUserHandler)
	}
}
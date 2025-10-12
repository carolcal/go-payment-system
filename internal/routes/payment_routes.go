package routes

import (
	"qr-payment/internal/handlers"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func SetUpPaymentRoutes(router *gin.Engine, phandler handlers.PaymentHandlers) {

	payments := router.Group("/payments")
	{
		payments.GET("", phandler.GetAllPaymentsHandler)
		payments.GET("/:user_id/:user_type", phandler.GetAllPaymentsByUserIdHandler)
	}

	payment := router.Group("/payment")
	{
		payment.POST("", phandler.CreatePaymentHandler)
		payment.GET("/:id", phandler.GetPaymentByIdHandler)
		payment.POST("/:user_id/pay", phandler.ProcessPaymentHandler)
		payment.DELETE("/:id", phandler.RemovePaymentHandler)
	}
}
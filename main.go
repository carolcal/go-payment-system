package main

import (
	"qr-payment/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/payments", handlers.GetAllPaymentsHandler)

	router.GET("/payment/:id", handlers.GetPaymentsByIdHandler)

	router.POST("/payment", handlers.CreatePaymentHandler)

	router.POST("/payment/:id/pay", handlers.MakePaymentHandler)

	router.DELETE("/payment/:id", handlers.RemovePaymentHandler)

	router.Run()
}

package main

import (
	"fmt"
	"database/sql"

	"qr-payment/handlers"
	"qr-payment/storage"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func main() {


	db, err := storage.NewDatabase()
    if err != nil { fmt.Print(err) }
    defer db.Close()


	router := gin.Default()
	setUpRoutes(router,db)
	router.Run()
}

func setUpRoutes(router *gin.Engine, db *sql.DB) {
	handler := &handlers.PaymentHandler{DB: db}

	payments := router.Group("/payments")
	{
		payments.GET("", handler.GetAllPaymentsHandler)
	}

	payment := router.Group("/payment")
	{
		payment.POST("", handler.CreatePaymentHandler)
		payment.GET("/:id", handler.GetPaymentByIdHandler)
		payment.POST("/:id/pay", handler.MakePaymentHandler)
		payment.DELETE("/:id", handler.RemovePaymentHandler)
	}
}

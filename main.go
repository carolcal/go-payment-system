package main

import (
	"database/sql"
	"fmt"

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
	router.Run("0.0.0.0:8080")
}

func setUpRoutes(router *gin.Engine, db *sql.DB) {
	phandler := &handlers.PaymentHandler{DB: db}
	uhandler := &handlers.UserHandler{DB: db}

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")

	router.GET("/", func (ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})

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

	payments := router.Group("/payments")
	{
		payments.GET("", phandler.GetAllPaymentsHandler)
	}

	payment := router.Group("/payment")
	{
		payment.POST("", phandler.CreatePaymentHandler)
		payment.GET("/:id", phandler.GetPaymentByIdHandler)
		payment.POST("/:id/pay", phandler.MakePaymentHandler)
		payment.DELETE("/:id", phandler.RemovePaymentHandler)
	}
}

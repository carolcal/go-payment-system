package main

import (
	// "database/sql"
	"fmt"

	"qr-payment/internal/core/services"
	"qr-payment/internal/handlers"
	"qr-payment/internal/infrastructure/database"
	"qr-payment/internal/infrastructure/repository"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := database.NewDatabase()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	urepository := repository.NewUserRepository(db)
	uservices := services.NewUserService(urepository)
	uhandler := handlers.NewUserHandlers(uservices)
	prepository := repository.NewPaymentRepository(db)
	pservices := services.NewPaymentService(prepository, uservices)
	phandler := handlers.NewPaymentHandlers(pservices)

	router := gin.Default()
	setUpRoutes(router, phandler, uhandler)
	router.Run("0.0.0.0:8080")
}

func setUpRoutes(router *gin.Engine, phandler handlers.PaymentHandlers, uhandler handlers.UserHandlers) {

	router.LoadHTMLGlob("web/*.html")
	router.Static("css", "./web/css")
	router.Static("js", "./web/js")

	router.GET("/", func(ctx *gin.Context) {
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

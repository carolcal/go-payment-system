package main

import (
	"fmt"
	"log"

	"qr-payment/internal/core/services"
	"qr-payment/internal/handlers"
	"qr-payment/internal/infrastructure/database"
	"qr-payment/internal/infrastructure/repository"
	"qr-payment/internal/routes"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}
	defer func() {
        if cerr := db.Close(); cerr != nil {
            fmt.Printf("database close error: %v\n", cerr)
        }
    }()

	urepository := repository.NewUserRepository(db)
	uservices := services.NewUserService(urepository)
	uhandler := handlers.NewUserHandlers(uservices)
	prepository := repository.NewPaymentRepository(db)
	pservices := services.NewPaymentService(prepository, uservices)
	phandler := handlers.NewPaymentHandlers(pservices)

	router := gin.Default()
	routes.SetUpWebRoutes(router)
	routes.SetUpUserRoutes(router, uhandler)
	routes.SetUpPaymentRoutes(router, phandler)
	router.Run("0.0.0.0:8080")
}



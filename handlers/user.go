package handlers

	// users := router.Group("/users")
	// {
	// 	users.GET("", handler.GetAllUsersHandler)
	// }

	// user := router.Group("/user")
	// {
	// 	user.POST("", handler.CreateUserHandler)
	// 	user.GET("/:id", handler.GetUserByIdHandler)
	// 	user.PUT("/:id/balance", handler.UpdateBalanceHandler)
	// 	user.DELETE("/:id", handler.RemoveUserHandler)
	// }

import (
	"fmt"
	"database/sql"

	"qr-payment/models"
	"qr-payment/storage"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) GetAllUsersHandler(ctx *gin.Context) {
	fmt.Println("Received GET request on /users")
	payments, err := storage.GetAllUsers(h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, payments)
}

func (h *UserHandler) GetUserByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := storage.GetUserById(id, h.DB)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, user)
}

func (h *UserHandler) CreateUserHandler(ctx *gin.Context) {
	var req models.CreateUserData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := &models.UserData{
		Name: req.Name,
		CPF: req.CPF,
		Balance: int(req.Balance * 100),
		City: req.City,
	}
	err := storage.CreateUser(user, h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, user)
}

func (h *UserHandler) UpdateBalanceHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateBalanceData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := storage.UpdateBalance(id, req, h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "balance updated successfully"})
}


func (h *UserHandler) RemoveUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := storage.RemoveUser(id, h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "deleted user successfully"})
}
package handlers

import (
	"qr-payment/internal/core/models"
	"qr-payment/internal/core/services"

	"github.com/gin-gonic/gin"
)

type UserHandlers interface {
	GetAllUsersHandler(ctx *gin.Context)
	GetUserByIdHandler(ctx *gin.Context)
	CreateUserHandler(ctx *gin.Context)
	UpdateBalanceHandler(ctx *gin.Context)
	RemoveUserHandler(ctx *gin.Context)
}

type userHandlers struct {
	service	services.UserService
}

func NewUserHandlers(service services.UserService) UserHandlers {
	return &userHandlers{
		service: service,
	}
}

func (h *userHandlers) GetAllUsersHandler(ctx *gin.Context) {
	payments, err := h.service.GetAllUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, payments)
}

func (h *userHandlers) GetUserByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := h.service.GetUserById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, user)
}

func (h *userHandlers) CreateUserHandler(ctx *gin.Context) {
	var req models.CreateUserData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, user)
}

func (h *userHandlers) UpdateBalanceHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var req models.UpdateBalanceData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateBalance(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "balance updated successfully"})
}


func (h *userHandlers) RemoveUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.service.RemoveUser(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "deleted user successfully"})
}
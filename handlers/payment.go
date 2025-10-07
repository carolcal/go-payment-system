package handlers

import (
	"database/sql"
	"fmt"

	"qr-payment/models"
	"qr-payment/storage"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	DB *sql.DB
}

func (h *PaymentHandler) GetAllPaymentsHandler(ctx *gin.Context) {
	payments, err := storage.GetAllPayments(h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, payments)
}

func (h *PaymentHandler) GetAllPaymentsByUserIdHandler(ctx *gin.Context) {
	user_id := ctx.Param("user_id")
	user_type_param := ctx.Param("user_type")

	user_type, isValid := models.IsValidTypeUser(user_type_param)
	if !isValid {
		ctx.JSON(400, gin.H{"error": "user_type deve ser 'receiver_id' ou 'payer_id'"})
		return
	}

	payments, err := storage.GetAllPaymentsByUserId(user_type, user_id, h.DB)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, payments)
}

func (h *PaymentHandler) GetPaymentByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	payment, err := storage.GetPaymentById(id, h.DB)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, payment)
}

func (h *PaymentHandler) CreatePaymentHandler(ctx *gin.Context) {
	var req models.CreatePaymentData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("amount: %f, receiver: %s", req.Amount, req.ReceiverId)

	user, err := storage.GetUserById(req.ReceiverId, h.DB)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}

	payment := &models.PaymentData{
		Amount:     int(req.Amount * 100),
		ReceiverId: req.ReceiverId,
	}
	err = storage.CreatePayment(user, payment, h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, payment)
}

func (h *PaymentHandler) ProcessPaymentHandler(ctx *gin.Context) {
	user_id := ctx.Param("user_id")
	var req models.ProcessPaymentData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := storage.ProcessPayment(user_id, req.QRCodeData, h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "payment made successfully"})
}

func (h *PaymentHandler) RemovePaymentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := storage.RemovePayment(id, h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "deleted payment successfully"})
}

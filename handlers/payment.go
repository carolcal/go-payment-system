package handlers

import (
	"database/sql"
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

	payment := &models.PaymentData{
		Amount: int(req.Amount * 100),
	}
	err := storage.CreatePayment(payment, h.DB)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, payment)
}

func (h *PaymentHandler) MakePaymentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := storage.MakePayment(id, h.DB)
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


// func (h *PaymentHandler)HtmlHandler(ctx *gin.Context) {
// 	ctx.HTML(200, "index.html", gin.H{})
// }
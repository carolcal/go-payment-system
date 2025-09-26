package handlers

import (
	"qr-payment/models"
	"qr-payment/storage"

	"github.com/gin-gonic/gin"
)

func GetPaymentsHandler(ctx *gin.Context) {
	payments, err := storage.GetPayments()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, payments)
}

func GetPaymentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	payment, err := storage.GetPayment(id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "payment not found"})
		return
	}
	ctx.JSON(200, payment)
}

func CreatePaymentHandler(ctx *gin.Context) {
	var req models.CreatePaymentData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	payment := &models.PaymentData{
		Amount:	req.Amount,
	}
	err := storage.CreatePayment(payment)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, payment)
}

func MakePaymentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := storage.MakePayment(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "payment made successfully"})
}

func RemovePaymentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := storage.RemovePayment(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "deleted payment successfully"})
}
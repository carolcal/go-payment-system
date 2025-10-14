package handlers

import (
	"qr-payment/internal/core/models"
	"qr-payment/internal/core/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandlers interface {
	GetAllPaymentsHandler(ctx *gin.Context)
	GetAllPaymentsByUserIdHandler(ctx *gin.Context)
	GetPaymentByIdHandler(ctx *gin.Context)
	CreatePaymentHandler(ctx *gin.Context)
	ProcessPaymentHandler(ctx *gin.Context)
	RemovePaymentHandler(ctx *gin.Context)
}

type paymentHandlers struct {
	service	services.PaymentService
}

func NewPaymentHandlers(service services.PaymentService) PaymentHandlers {
	return &paymentHandlers{
		service: service,
	}
}

func (h *paymentHandlers) GetAllPaymentsHandler(ctx *gin.Context) {
	payments, err := h.service.GetAllPayments(ctx.Request.Context())
	if err != nil {
		ctx.JSON(models.HTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, payments)
}

func (h *paymentHandlers) GetAllPaymentsByUserIdHandler(ctx *gin.Context) {
	user_id := ctx.Param("user_id")
	user_type_param := ctx.Param("user_type")

	payments, err := h.service.GetAllPaymentsByUserId(ctx.Request.Context(), user_type_param, user_id)
	if err != nil {
		ctx.JSON(models.HTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, payments)
}

func (h *paymentHandlers) GetPaymentByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	payment, err := h.service.GetPaymentById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(models.HTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, payment)
}

func (h *paymentHandlers) CreatePaymentHandler(ctx *gin.Context) {
	var req models.CreatePaymentData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.CreatePayment(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(models.HTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, payment)
}

func (h *paymentHandlers) ProcessPaymentHandler(ctx *gin.Context) {
	user_id := ctx.Param("user_id")

	var req models.ProcessPaymentData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ProcessPayment(ctx.Request.Context(), user_id, &req)
	if err != nil {
		ctx.JSON(models.HTTPStatus(err), gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"status": "payment made successfully"})
}

func (h *paymentHandlers) RemovePaymentHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.service.RemovePayment(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(models.HTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"status": "deleted payment successfully"})
}

package main

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusPaid      PaymentStatus = "paid"
	StatusFailed    PaymentStatus = "failed"
	StatusExpired   PaymentStatus = "expired"
)

type Payment struct {
	ID			string			`json:"id"`
	Amount		float64			`json:"amount"`
	Status		PaymentStatus	`json:"status"`
	CreatedAt	time.Time		`json:"created_at"`
	ExpiresAt	time.Time		`json:"expires_at"`
	QRCodeData	string			`json:"qr_code_data"`
}

type CreatePayment struct {
	Amount		float64			`json:"amount"`
}

var paymentsDB = make(map[string]*Payment)

var mu sync.Mutex

func getPayment(id string) (*Payment, error) {
	mu.Lock()
	defer mu.Unlock()
	payment, exists := paymentsDB[id]
	if !exists {
		return nil, fmt.Errorf("payment not found")
	}
	return payment, nil
}

func getPayments() (map[string]*Payment, error) {
	mu.Lock()
	defer mu.Unlock()
	return paymentsDB, nil
}

func generateID() string {
	return "pay_" + fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateQRCode(id string) ([]byte, error) {
	var png []byte
	url := "http://localhost:8080/payments/" + id + "/pay"
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qrcode")
	}
	return png, nil
}

func createPayment(p *Payment) error {
	mu.Lock()
	defer mu.Unlock()
	id := generateID()
	p.ID = id
	p.CreatedAt = time.Now()
	p.ExpiresAt = time.Now().Add(15 * time.Minute)
	p.Status = StatusPending
	qrcodeBytes, err := generateQRCode(id)
	if err != nil {
		return err
	}
	p.QRCodeData = "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrcodeBytes)
	paymentsDB[id] = p
	return nil
}

func makePayment(id string) error {
	mu.Lock()
	defer mu.Unlock()
	payment, exists := paymentsDB[id]
	if !exists {
		return fmt.Errorf("payment not found")
	}
	switch payment.Status {
	case StatusPaid:
		return fmt.Errorf("payment already completed")
	case StatusFailed:
		return fmt.Errorf("payment failed previously")
	case StatusExpired:
		return fmt.Errorf("payment has expired")
	}
	
	if time.Now().After(payment.ExpiresAt) {
		payment.Status = StatusExpired
		return fmt.Errorf("payment has expired")
	}
	
	payment.Status = StatusPaid
	return nil
}

func removePayment(id string) (error) {
	mu.Lock()
	defer mu.Unlock()
	_, exists := paymentsDB[id]
	if !exists {
		return fmt.Errorf("payment not found")
	}
	delete(paymentsDB, id)
	return nil
}

func main() {
	router := gin.Default()

	router.GET("/payments", func(ctx *gin.Context) {
		payments, err := getPayments()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, payments)
	})

	router.GET("/payment/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		payment, err := getPayment(id)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Payment not found"})
			return
		}
		ctx.JSON(200, payment)
	})

	router.POST("/payment", func(ctx *gin.Context) {
		var req CreatePayment
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		payment := &Payment{
			Amount:	req.Amount,
		}
		err := createPayment(payment)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(201, payment)
	})

	router.POST("/payment/:id/pay", func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := makePayment(id)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"status": "Payment made successfully"})
	})

	router.DELETE("/payment/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := removePayment(id)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"status": "Deleted payment successfully"})
	})

	router.Run()
}
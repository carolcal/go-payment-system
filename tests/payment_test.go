package tests

import (
	"bytes"
	"net/http"
	"testing"
	"encoding/json"

	"qr-payment/models"
	_ "github.com/mattn/go-sqlite3"
)

var router = "http://localhost:8080/"

var createdPayment models.PaymentData

var createdPaymentsMap = make(map[string]*models.PaymentData)

func TestCreatePayment(t *testing.T) {
	payload := `{"Amount": 42.87}`
	url := router + "payment"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %d", resp.StatusCode)
	}

	var result models.PaymentData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expected := models.PaymentData{
		Amount: 4287,
		Status: models.StatusPending,
	}

	if (result.Amount != expected.Amount) || (result.Status != expected.Status) {
		t.Errorf("Expected values: { Amount: %d, Status: %s }; got: { Amount: %d, Status: %s }", result.Amount, result.Status, expected.Amount, expected.Status)
	}
	if result.ID == "" {
		t.Errorf("Expected non-empty ID, got: %s", result.ID)
	}
	if result.ExpiresAt.IsZero() {
		t.Errorf("Expected valid ExpiresAt timestamp, got: %v", result.ExpiresAt)
	}
	gap := int(result.ExpiresAt.Sub(result.CreatedAt).Minutes())
	if gap != 15 {
		t.Errorf("Expected ExpiresAt to be 15 minutes after CreatedAt, got: %v", gap)
	}
	if result.QRCodeData == "" {
		t.Errorf("Expected non-empty QRCodeData, got: %s", result.QRCodeData)
	}

	createdPayment = result
	createdPaymentsMap[createdPayment.ID] = &createdPayment
}

func TestGetPaymentById(t *testing.T) {
	url := router + "payment/" + createdPayment.ID

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	var result models.PaymentData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result != createdPayment {
		t.Errorf("Expected payment: %+v, got %+v", createdPayment, result)
	}
}

func TestGetAllPayments(t *testing.T) {
	url := router + "payments"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]*models.PaymentData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode responde: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected map with length 1, got: %d", len(result))
	}

	if *result[createdPayment.ID] != createdPayment {
		t.Errorf("Expected payment: %+v, got %+v", createdPayment, result[createdPayment.ID])
	}
}

func TestMakePayment(t *testing.T) {
	url := router + "payment/" + createdPayment.ID + "/pay"

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetPaymentByIdStatusPaid(t *testing.T) {
	url := router + "payment/" + createdPayment.ID

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	var result models.PaymentData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Status != models.StatusPaid {
		t.Errorf("Failed to change Status")
	}
}

func TestRemovePayment(t *testing.T) {
	url := router + "payment/" + createdPayment.ID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestGetEmptyPayments(t *testing.T) {
	url := router + "payments"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]*models.PaymentData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode responde: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected map with length 0, got: %d", len(result))
	}
}
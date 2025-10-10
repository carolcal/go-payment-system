package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"qr-payment/internal/core/models"

	_ "github.com/mattn/go-sqlite3"
)

var createdPayment models.PaymentData

var users = make(map[string]*models.UserData)

var receiverID string

var payerID string

func TestCreateUsers(t *testing.T){
	payload := `{
		"name": "Arthur Dent",
		"cpf": "11111111111",
		"balance": 5000.00,
		"city": "Earth"
	}`
	url := router + "user"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %d", resp.StatusCode)
	}

	var result1 models.UserData
	err = json.NewDecoder(resp.Body).Decode(&result1)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	users[result1.ID] = &result1
	receiverID = result1.ID

	//Add second user
	payload = `{
		"name": "Ford Prefect",
		"cpf": "22222222222",
		"balance": 3000.00,
		"city": "Betelgeuse"
	}`
	url = router + "user"

	resp, err = http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %d", resp.StatusCode)
	}
	
	var result2 models.UserData
	err = json.NewDecoder(resp.Body).Decode(&result2)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	users[result2.ID] = &result2
	payerID = result2.ID
}

func TestCreatePayment(t *testing.T) {
	if receiverID == "" {
        t.Skip("receiverID not initialized; ensure TestCreateUsers runs first")
    }

	payload := `{
		"amount": 100.00,
		"receiver_id": "` + receiverID + `"
	}`
	url := router + "payment"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %d: %s", resp.StatusCode, resp.Body)
	}

	var result models.PaymentData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expected := models.PaymentData{
		Amount: 10000,
		Status: models.StatusPending,
	}

	if (result.Amount != expected.Amount) || (result.Status != expected.Status) {
		t.Errorf("Expected values: { Amount: %d, Status: %s }; got: { Amount: %d, Status: %s }",expected.Amount, expected.Status, result.Amount, result.Status)
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

func TestProcessPayment(t *testing.T) {
	url := router + "payment/" + payerID + "/pay"
	payload := `{ "qr_code_data": "` + createdPayment.QRCodeData + `" }`

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d, %s", resp.StatusCode, body)
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

func TestUpdatedUsersBalance(t *testing.T) {
	url := router + "users"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var results map[string]*models.UserData
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 users, got %d", len(results))
	}

	expectedReceiverBalance := users[receiverID].Balance + 10000
	expectedPayerBalance := users[payerID].Balance - 10000

	if results[payerID].Balance != expectedPayerBalance {
		t.Errorf("Expected payer balance: %d, got %d", expectedPayerBalance, results[payerID].Balance)
	}

	if results[receiverID].Balance != expectedReceiverBalance {
		t.Errorf("Expected receiver balance: %d, got %d", expectedReceiverBalance, results[receiverID].Balance)
	}
}

func TestCleanupUsers(t *testing.T) {
	for id := range users {
		url := router + "user/" + id

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
}

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"qr-payment/internal/core/models"

	_ "github.com/mattn/go-sqlite3"
)


var createdUser models.UserData

var createdUsersMap = make(map[string]*models.UserData)

func TestCreateUser(t *testing.T){
	payload := `{
		"name": "Arthur Dent",
		"cpf": "123.721.280-43",
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
		fmt.Println(resp)
		t.Fatalf("Expected status code 201, got %d", resp.StatusCode)
	}

	var result models.UserData
	err =json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expected := models.UserData {
		Name: "Arthur Dent",
		CPF: "12372128043",
		Balance: 500000,
		City: "Earth",
	}

	if (result.Name != expected.Name) || (result.CPF != expected.CPF) || (result.Balance != expected.Balance) || (result.City != expected.City) {
		t.Errorf("Expected values: { Name: %s, CPF: %s, Balance: %d, City: %s }; got: { Name: %s, CPF: %s, Balance: %d, City: %s }", result.Name, result.CPF, result.Balance, result.City, expected.Name, expected.CPF, expected.Balance, expected.City)
	}
	if result.ID == "" {
		t.Errorf("Expected non-empty ID, got: %s", result.ID)
	}
	if result.CreatedAt.IsZero() {
		t.Errorf("Expected valid CreatedAt timestamp, got: %v", result.CreatedAt)
	}
	if result.UpdatedAt.IsZero() {
		t.Errorf("Expected valid UpdatedAt timestamp, got: %v", result.UpdatedAt)
	}

	createdUser = result
	createdUsersMap[createdUser.ID] = &createdUser
}

func TestGetUserByID(t *testing.T) {
	url := router + "user/" + createdUser.ID

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	var result models.UserData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expected := createdUser

	if (result.ID != expected.ID) || (result.Name != expected.Name) || (result.CPF != expected.CPF) || (result.Balance != expected.Balance) || (result.City != expected.City) {
		t.Errorf("Expected values: { ID: %s, Name: %s, CPF: %s, Balance: %d, City: %s }; got: { ID: %s, Name: %s, CPF: %s, Balance: %d, City: %s }", result.ID, result.Name, result.CPF, result.Balance, result.City, expected.ID, expected.Name, expected.CPF, expected.Balance, expected.City)
	}
	if result.CreatedAt.IsZero() {
		t.Errorf("Expected valid CreatedAt timestamp, got: %v", result.CreatedAt)
	}
	if result.UpdatedAt.IsZero() {
		t.Errorf("Expected valid UpdatedAt timestamp, got: %v", result.UpdatedAt)
	}
}

func TestUpdateUserBalance(t *testing.T) {
	payload := `{
		"diff": -1000.00
	}`
	url := router + "user/" + createdUser.ID + "/balance"

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	expectedBalance := createdUser.Balance - 100000

	createdUser.Balance = expectedBalance
	createdUsersMap[createdUser.ID] = &createdUser

}

func TestGetAllUsers(t *testing.T) {
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

	if len(results) == 0 {
		t.Errorf("Expected non-empty results, got length 0")
	}

	for id, user := range createdUsersMap {
		result, exists := results[id]
		if !exists {
			t.Errorf("Expected user ID %s to be in results, but it was missing", id)
			continue
		}
		if *result != *user {
			t.Errorf("For user ID %s, expected %+v, got %+v", id, *user, *result)
		}
	}
}

func TestRemoveUser(t *testing.T) {
	url := router + "user/" + createdUser.ID

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

func TestGetEmptyUsers(t *testing.T) {
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

	if len(results) != 0 {
		t.Errorf("Expected empty results, got length %d", len(results))
	}
}

# Payment System

This is my first project in my Golang learning journey.
It implements a simple payment system where you can create and process payments using QR codes.

---

## Features

The API exposes **5 endpoints**:

* **GET `/payments`** → returns all payments
* **GET `/payment/:id`** → returns a specific payment by ID
* **POST `/payment`** → creates a new payment with QR code (requires an `amount` in the request body)
* **POST `/payment/:id/pay`** → processes (marks as paid) a payment
* **DELETE `/payment/:id`** → removes a payment from the database

---

## Tech Stack

* **Language:** Go (Golang)
* **Database:** SQLite

---

## Running the Project

Make sure you have Go installed. Then run:

```bash
go run main.go
```

By default, the server will start on **port 8080**.

---

## Testing the API

You can interact with the project in two ways:

* **Frontend:**
  Access [http://localhost:8080](http://localhost:8080) in your browser to interact with the frontend.

* **API:**
  Use Postman, `curl`, or any HTTP client to test the API directly.
  Example:

  ```bash
  curl http://localhost:8080/payments
  ```

---

## Automated Tests

This project includes automated tests. Make sure the server is running on port 8080, then run:
```bash
go test tests/payment_test.go
```


# Payment System

This is my first project in my Golang learning journey.
It implements a simple payment system where you can create and process payments using QR codes.

---

## Features

The API exposes the following **endpoints**:

### /users
* **GET `/`** → Returns all users.

### /user
* **GET `/:id`** → Returns user filtered by parameter `id`.
* **POST `/`** → Creates a new user.
  - **Body:**
  ```json
  {
    "name": "string(required)", 
    "cpf": "string(required)", 
    "balance": "float,(required)",
    "city": "string(required)"
  }
  ```
* **PUT `/:id/balance`** → Updates user's balance.
  - **Body:**
  ```json
  {
    "diff": "float(required)"
  }
  ```
  - A positive value deposits money into the account, and a negative value withdraws from it.
* **DELETE `/:id`** → Deletes user by the given parameter `id`.

### /payments
* **GET `/`** → Returns all payments.
* **GET `/:user_id/:user_type`** → Returns all payments for the given `user_id` and `user_type` (`receiver_id` or `payer_id`)
  - Passing `receiver_id` returns received payments (deposits).
  - Passing `payer_id` returns sent payments (withdrawals).

### /payment
* **GET `/:id`** → Returns a specific payment filtered by `id`.
* **POST `/`** → Creates a new payment with a QR code.
  - **Body:**
  ```json
   {
     "receiver_id": "string(required)",
     "amount": "float(optional)`"
   }
  ```
  - If you don't provide `amount`, it creates a QR code payment without a fixed value. In that case, you must provide an `amount` when processing the payment (see endpoint below).
  - To make it simple, it uses the user's cpf as pix key.
* **POST `/:user_id/pay`** → Processes a payment by transferring money from one account to another and marking it as `paid`.
  - **Body:**
  ```json
  {
    "qr_code_data": "string(required)",
    "amount": "float(optional)`"
  }
  ```
  - In `user_id`, you must provide the ID of the user who is making the payment (payer_id).
  - If the 15-minute processing time limit has elapsed, the payment will not be processed and will be marked as `expired`.
  - If something goes wrong, it will be marked as `failed`.
  - The program will get `TransactionAmount` from the QR code. If the QR code has no predefined value, you must provide an `amount` in the request body for the payment to be processed.
* **DELETE `/:id`** → Deletes payment by the given parameter `id`.

---

## Tech Stack

* **Database:** SQLite
* **Backend:** Go (Golang)
* **Frontend:** Vanila Javascript

---

## Running the Project

Make sure you have Go installed. Then run:

```bash
go run main.go
```

By default, the server will start on **port 8080** and will create database in the project repository.

---

## Testing the API

You can interact with the project in two ways:

* **Frontend:**
  The frontend part is very simple and was created only to make testing the API easier.
  Access [http://localhost:8080](http://localhost:8080) in your browser to interact with it.

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
**Note:** You must delete any pre-created database file for the tests to run correctly.


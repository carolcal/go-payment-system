package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const usersTable string = `
	CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	name TEXT NOT NULL,
	cpf VARCHAR(11) NOT NULL,
	balance INTEGER
);`

const paymentTable string = `
	CREATE TABLE IF NOT EXISTS payments (
	id TEXT PRIMARY KEY,
	amount INTEGER NOT NULL,
	status VARCHAR(10) NOT NULL,
	created_at DATETIME NOT NULL,
	expired_at DATETIME NOT NULL,
	qrcode_data TEXT NOT NULL
);`

func NewDatabase() (*sql.DB, error){

	const file string = "payment-system.db"

	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
    }

	if _, err := db.Exec(usersTable); err != nil {
		return nil, fmt.Errorf("failed to execute schema creation: %w", err)
	}

	if _, err := db.Exec(paymentTable); err != nil {
		return nil, fmt.Errorf("failed to execute schema creation: %w", err)
	}

	return db, nil
}
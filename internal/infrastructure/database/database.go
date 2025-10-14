package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"qr-payment/internal/core/models"
)

const usersTable string = `
	CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	name TEXT NOT NULL,
	cpf VARCHAR(11) NOT NULL,
	balance INTEGER,
	city TEXT NOT NULL
);`

const paymentTable string = `
	CREATE TABLE IF NOT EXISTS payments (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL,
	expired_at DATETIME NOT NULL,
	amount INTEGER NOT NULL,
	status VARCHAR(10) NOT NULL,
	receiver_id TEXT NOT NULL,
	payer_id TEXT,
	qr_code_data TEXT NOT NULL
);`

func NewDatabase() (*sql.DB, error){

	const file string = "payment-system.db"

	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, &models.Err{
			Op: "NewDatabase", 
			Status: models.Dependency, 
			Msg: "Failed to open database.", 
			Err: err,
		}
    }

	if _, err := db.Exec(usersTable); err != nil {
		return nil, &models.Err{
			Op: "NewDatabase", 
			Status: models.Dependency,
			Msg: "Failed to execute users schema creation.", 
			Err: err,
		}
	}

	if _, err := db.Exec(paymentTable); err != nil {
		return nil, &models.Err{
			Op: "NewDatabase", 
			Status: models.Dependency, 
			Msg: "Failed to execute payments schema creation.", 
			Err: err,
		}
	}

	return db, nil
}
package storage

import (
	"fmt"
	"time"
	"database/sql"

	"qr-payment/models"
	"qr-payment/utils"
)

func GetUserById(id string, db *sql.DB) (*models.UserData, error) {
	var user models.UserData

	row := db.QueryRow(`SELECT * FROM users WHERE id=?`, id)
	err := utils.ScanUserRow(row, &user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao escanear usuário: %w", err)
	}

	return &user, nil
}

func GetAllUsers(db *sql.DB) (map[string]*models.UserData, error) {
	allUsers := make(map[string]*models.UserData)

	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserData
		err := utils.ScanUserRows(rows, &user)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear usuário: %w", err)
		}
		allUsers[user.ID] = &user
	}

	return allUsers, nil
}

func CreateUser(u *models.UserData, db *sql.DB) error {
	id := utils.GenerateID("user")
	u.ID = id
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	_, err := db.Exec("INSERT INTO users VALUES(?, ?, ?, ?, ?, ?, ?);", u.ID, u.CreatedAt, u.UpdatedAt, u.Name, u.CPF, u.Balance, u.City)
	if (err != nil) {
		return fmt.Errorf("falha ao criar novo usuário")
	}

	return nil
}

func UpdateBalance(id string, req models.UpdateBalanceData, db *sql.DB) error {
	var user models.UserData

	row := db.QueryRow(`SELECT * FROM users WHERE id=?`, id)
	err := utils.ScanUserRow(row, &user)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("usuário não encontrado")
		}
		return fmt.Errorf("erro ao escanear usuário: %w", err)
	}

	newBalance := user.Balance + int(req.Diff * 100)

	_, err = db.Exec("UPDATE users SET balance=? WHERE id=?", newBalance, id)
	if err != nil {
		return fmt.Errorf("falha ao atualizar usuário para saldo %d: %w", newBalance, err)
	}

	return nil
}

func RemoveUser(id string, db *sql.DB) error {
	var user models.UserData

	row := db.QueryRow(`SELECT * FROM users WHERE id=?`, id)
	err := utils.ScanUserRow(row, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("usuário não encontrado")
		}
		return fmt.Errorf("erro ao escanear usuário: %w", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("falha ao deletar usuário: %w", err)
	}

	return nil
}

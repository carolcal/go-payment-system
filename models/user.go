package models

import (
	"time"
)

type UserData struct {
	ID			string		`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"expires_at"`
	Name		string		`json:"name"`
	CPF			string		`json:"cpf"`
	Balance		int			`json:"balance"`
}

type CreateUserData struct {
	Name		string		`json:"name"`
	CPF			string		`json:"cpf"`
	Balance		float64		`json:"balance"`
}

type UpdateBalanceData struct {
	Diff		float64		`json:"diff"`
}
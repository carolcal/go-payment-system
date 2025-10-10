package models

import (
	"time"
)

type UserData struct {
	ID			string		`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Name		string		`json:"name"`
	CPF			string		`json:"cpf"`
	Balance		int			`json:"balance"`
	City		string		`json:"city"`
}

type CreateUserData struct {
	Name		string		`json:"name" binding:"required"`
	CPF			string		`json:"cpf" binding:"required"`
	Balance		float64		`json:"balance" binding:"required"`
	City		string		`json:"city" binding:"required"`
}

type UpdateBalanceData struct {
	Diff		float64		`json:"diff"`
}

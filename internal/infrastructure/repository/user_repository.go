package repository

import (
	"database/sql"

	"qr-payment/internal/core/models"
	"qr-payment/internal/utils"
)

type UserRepository interface {
	FindAll() (map[string]*models.UserData, error)
	FindById(id string) (*models.UserData, error)
	FindByNameAndCPF(name string, cpf string) (*models.UserData, error)
	Create(ud *models.UserData) error
	UpdateBalance(id string, newBalance int) error
	Delete(id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) FindById(id string) (*models.UserData, error) {
	var user models.UserData

	row := u.db.QueryRow(`SELECT * FROM users WHERE id=?`, id)
	err := utils.ScanUserRow(row, &user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.Err{Op: "UserRepository.FindById", Status: models.NotFound, Msg: "User Not Found."}
		}
		return nil, &models.Err{Op: "UserRepository.FindById", Status: models.Dependency, Msg: "Error when scanning payment.", Err: err}
	}

	return &user, nil
}

func (u *userRepository) FindAll() (map[string]*models.UserData, error) {
	allUsers := make(map[string]*models.UserData)

	rows, err := u.db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, &models.Err{Op: "UserRepository.FindAll", Status: models.Dependency, Msg: "Error executing user query.", Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		usr := new(models.UserData)
		err := utils.ScanUserRows(rows, usr)
		if err != nil {
			return nil, &models.Err{Op: "UserRepository.FindAll", Status: models.Dependency, Msg: "Error when scanning payment.", Err: err}
		}
		allUsers[usr.ID] = usr
	}
	if err := rows.Err(); err != nil {
		return nil, &models.Err{Op: "UserRepository.FindAll", Status: models.Dependency, Msg: "Error when iterating users", Err: err}
    }

	return allUsers, nil
}

func (u *userRepository) FindByNameAndCPF(name string, cpf string) (*models.UserData, error) {
	var user models.UserData

	row := u.db.QueryRow(`SELECT * FROM users WHERE name=? AND cpf=? LIMIT 1;`, name, cpf)
	err := utils.ScanUserRow(row, &user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.Err{Op: "UserRepository.FindByNameAndCPF", Status: models.NotFound, Msg: "User Not Found."}
		}
		return nil, &models.Err{Op: "UserRepository.FindByNameAndCPF", Status: models.Dependency, Msg: "Error when scanning payment.", Err: err}
	}

	return &user, nil
}

func (u *userRepository) Create(ud *models.UserData) error {
	_, err := u.db.Exec("INSERT INTO users VALUES(?, ?, ?, ?, ?, ?, ?);", ud.ID, ud.CreatedAt, ud.UpdatedAt, ud.Name, ud.CPF, ud.Balance, ud.City)
	if (err != nil) {
		return &models.Err{Op: "UserRepository.Create", Status: models.Dependency, Msg: "Fail to create new user.", Err: err}
	}

	return nil
}

func (u *userRepository) UpdateBalance(id string, newBalance int) error {
	_, err := u.db.Exec("UPDATE users SET balance=? WHERE id=?", newBalance, id)
	if err != nil {
		return &models.Err{Op: "UserRepository.UpdateBalance", Status: models.Dependency, Msg: "Fail to update user balance.", Err: err}
	}

	return nil
}

func (u *userRepository) Delete(id string) error {
	_, err := u.db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return &models.Err{Op: "UserRepository.Delete", Status: models.Dependency, Msg: "Fail to delete user.", Err: err}
	}

	return nil
}

package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"regexp"

	"qr-payment/internal/core/models"
	"qr-payment/internal/infrastructure/repository"
	"qr-payment/internal/utils"
)

type UserService interface {
	GetAllUsers(ctx context.Context) (map[string]*models.UserData, error)
	GetUserById(ctx context.Context, id string) (*models.UserData, error)
	GetUserByNameAndCPF(ctx context.Context, name string, cpf string) (*models.UserData, error)
	CreateUser(ctx context.Context, cud *models.CreateUserData) (*models.UserData, error)
	ValidateBalance(ctx context.Context, id string, amount int) (*models.UserData, error)
	UpdateBalance(ctx context.Context, id string, req models.UpdateBalanceData) error
	RemoveUser(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetUserById(ctx context.Context, id string) (*models.UserData, error) {
	return s.repo.FindById(id)
}

func (s *userService) GetAllUsers(ctx context.Context) (map[string]*models.UserData, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByNameAndCPF(ctx context.Context, name string, cpf string) (*models.UserData, error) {
	return s.repo.FindByNameAndCPF(name, cpf)
}

func (s *userService) CreateUser(ctx context.Context, cud *models.CreateUserData) (*models.UserData, error) {
	id := utils.GenerateID("user")
	re := regexp.MustCompile(`\D+`)
	ud := &models.UserData{
		ID:			id,
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cud.Name,
		CPF:		re.ReplaceAllString(cud.CPF, ""),
		Balance:	int(cud.Balance * 100),
		City:		cud.City,
	}
	if err := s.repo.Create(ud); err != nil {
		return nil, err
	}
	return ud, nil
}

func (s *userService) ValidateBalance(ctx context.Context, id string, amount int) (*models.UserData, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao escanear usuário1: %w", err)
	}
	newBalance := user.Balance + amount
	if newBalance < 0 {
		return nil, fmt.Errorf("saldo insuficiente")
	}
	return user, nil
}

func (s *userService) UpdateBalance(ctx context.Context, id string, req models.UpdateBalanceData) error {
	amount := int(req.Diff * 100)
	user, err := s.ValidateBalance(ctx, id, amount)
	if err != nil {
		return err
	}
	newBalance := user.Balance + amount
	return s.repo.UpdateBalance(id, newBalance)
}

func (s *userService) RemoveUser(ctx context.Context, id string) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("usuário não encontrado")
		}
		return fmt.Errorf("erro ao escanear usuário: %w", err)
	}
	return s.repo.Delete(id)
}

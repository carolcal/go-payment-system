package services

import (
	"context"
	"time"

	"qr-payment/internal/core/models"
	"qr-payment/internal/core/validators"
	"qr-payment/internal/infrastructure/repository"
	"qr-payment/internal/utils"
)

type UserService interface {
	GetAllUsers(ctx context.Context) (map[string]*models.UserData, error)
	GetUserById(ctx context.Context, id string) (*models.UserData, error)
	GetUserByNameAndCPF(ctx context.Context, name string, cpf string) (*models.UserData, error)
	CreateUser(ctx context.Context, cud *models.CreateUserData) (*models.UserData, error)
	UpdateBalance(ctx context.Context, id string, req models.UpdateBalanceData) error
	RemoveUser(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
	val validators.UserValidator
}

func NewUserService(repo repository.UserRepository, val validators.UserValidator) UserService {
	return &userService{
		repo: repo,
		val: val,
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
	cleanCPF, err := s.val.ValidateCPF(cud.CPF)
	if err != nil {
		return nil, err
	}
	ud := &models.UserData{
		ID:			id,
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cud.Name,
		CPF:		cleanCPF,
		Balance:	int(cud.Balance * 100),
		City:		cud.City,
	}
	if err := s.repo.Create(ud); err != nil {
		return nil, err
	}
	return ud, nil
}

func (s *userService) UpdateBalance(ctx context.Context, id string, req models.UpdateBalanceData) error {
	amount := int(req.Diff * 100)
	user, err := s.val.ValidateBalance(id, amount)
	if err != nil {
		return err
	}
	newBalance := user.Balance + amount
	return s.repo.UpdateBalance(id, newBalance)
}

func (s *userService) RemoveUser(ctx context.Context, id string) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

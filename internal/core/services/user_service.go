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
	ValidateCPF(cpf string) (string, error)
	CreateUser(ctx context.Context, cud *models.CreateUserData) (*models.UserData, error)
	ValidateBalance(id string, amount int) (*models.UserData, error)
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

func calculateDigit(cleanCPF string, digitIndex int) error {
	var numCheck	int
	var mod			int

	sum := 0
	j := 0
	for i := digitIndex + 1; i > 1; i-- {
		digit := int(cleanCPF[j] - '0')
		sum += digit * i
		j++
	}
	fmt.Println(`index :`, digitIndex)
	fmt.Println(`sum :`, sum)
	mod = sum % 11
	fmt.Println(`mod :`, mod)
	if mod < 2 {
		numCheck = 0
	} else {
		numCheck = 11 - mod
	}
	if numCheck != int(cleanCPF[digitIndex] - '0') {
		return fmt.Errorf("invalid CPF number")
	}
	return nil
}

func verifyDigitsEquals(cleanCPF string) error {
	firstDigit := int(cleanCPF[0] - '0')
	for i := range 9 {
		if int(cleanCPF[i] - '0') != firstDigit {
			return nil
		}
	}
	return fmt.Errorf("invalid CPF number")
}

func (s *userService) ValidateCPF(cpf string) (string, error) {
	re := regexp.MustCompile(`\D+`)
	cleanCPF := re.ReplaceAllString(cpf, "")
	if len(cleanCPF) != 11 {
		return "", fmt.Errorf("CPF must have 11 digits")
	}
	if equalsErr := verifyDigitsEquals(cleanCPF); equalsErr != nil {
		return "", equalsErr
	}

	if firstCheckErr := calculateDigit(cleanCPF, 9); firstCheckErr != nil {
		return "", firstCheckErr
	}
	if secondCheckErr := calculateDigit(cleanCPF, 10); secondCheckErr != nil {
		return "", secondCheckErr
	}
	
	return cleanCPF, nil
}

func (s *userService) CreateUser(ctx context.Context, cud *models.CreateUserData) (*models.UserData, error) {
	id := utils.GenerateID("user")
	cleanCPF, err := s.ValidateCPF(cud.CPF)
	if err != nil {
		fmt.Println(err)
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

func (s *userService) ValidateBalance(id string, amount int) (*models.UserData, error) {
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
	user, err := s.ValidateBalance(id, amount)
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

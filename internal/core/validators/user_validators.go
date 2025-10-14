package validators

import (
	"regexp"

	"qr-payment/internal/core/models"
	"qr-payment/internal/infrastructure/repository"
)

type UserValidator interface {
	ValidateCPF(cpf string) (string, error)
	ValidateBalance(id string, amount int) (*models.UserData, error)
}

type userValidator struct {
	repo repository.UserRepository
}

func NewUserValidator(repo repository.UserRepository) UserValidator {
	return &userValidator{
		repo: repo,
	}
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
	mod = sum % 11
	if mod < 2 {
		numCheck = 0
	} else {
		numCheck = 11 - mod
	}
	if numCheck != int(cleanCPF[digitIndex] - '0') {
		return &models.Err{Op: "ValidateCPF", Status: models.Invalid, Msg: "Invalid CPF number."}
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
	return &models.Err{Op: "ValidateCPF", Status: models.Invalid, Msg: "Invalid CPF number."}
}

func (v *userValidator) ValidateCPF(cpf string) (string, error) {
	re := regexp.MustCompile(`\D+`)
	cleanCPF := re.ReplaceAllString(cpf, "")
	if len(cleanCPF) != 11 {
		return "", &models.Err{Op: "ValidateCPF", Status: models.Invalid, Msg: "CPF must have 11 digits."}
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

func (v *userValidator) ValidateBalance(id string, amount int) (*models.UserData, error) {
	user, err := v.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	newBalance := user.Balance + amount
	if newBalance < 0 {
		return nil, &models.Err{Op: "ValidateBalance", Status: models.Precondition, Msg: "Insufficient Balance."}
	}
	return user, nil
}

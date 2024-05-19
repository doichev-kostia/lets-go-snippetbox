package mocks

import (
	"github.com/google/uuid"
	"snippetbox.doichevkostia.dev/internal/models"
)

var UserID = uuid.New()

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) (uuid.UUID, error) {
	switch email {
	case "dupe@example.com":
		return uuid.UUID{}, models.ErrDuplicateEmail
	default:
		return uuid.New(), nil
	}
}

func (m *UserModel) Authenticate(email, password string) (uuid.UUID, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return UserID, nil
	}

	return uuid.UUID{}, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id uuid.UUID) (bool, error) {
	switch id {
	case UserID:
		return true, nil
	default:
		return false, nil
	}
}

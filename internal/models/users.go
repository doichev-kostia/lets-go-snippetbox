package models

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModelInterface interface {
	Insert(name, email, password string) (uuid.UUID, error)
	Authenticate(email, password string) (uuid.UUID, error)
	Exists(id uuid.UUID) (bool, error)
}

type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	HashedPassword []byte
	CreateTime     time.Time
}

type UserModel struct {
	PasswordCost int
	DB           *sql.DB
}

func (m *UserModel) Insert(name, email, password string) (uuid.UUID, error) {
	// TODO: lock the table to prevent race conditions
	exists, err := m.EmailExists(email)
	if err != nil {
		return uuid.UUID{}, err
	}

	if exists {
		return uuid.UUID{}, ErrDuplicateEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), m.PasswordCost)
	if err != nil {
		return uuid.UUID{}, err
	}

	stmt := `insert into "users" ("id", "name", "email", "hashed_password")
	values (?, ?, ?, ?)`

	id := uuid.New()
	_, err = m.DB.Exec(stmt, id, name, email, string(hashedPassword))
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (m *UserModel) Authenticate(email, password string) (uuid.UUID, error) {
	usr, err := m.ByEmail(email)

	if err != nil {
		return uuid.UUID{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(usr.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.UUID{}, ErrInvalidCredentials
		} else {
			return uuid.UUID{}, err
		}
	}

	return usr.ID, nil
}

func (m *UserModel) Exists(id uuid.UUID) (bool, error) {
	var exists bool
	stmt := `select exists(select true from "users" where id = ?)`

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) ByEmail(email string) (User, error) {
	stmt := `select "id", "name", "email", "hashed_password", "create_time" from "users" where email = ?`

	var u User
	err := m.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		} else {
			return User{}, err
		}
	}

	return u, nil
}

func (m *UserModel) EmailExists(email string) (bool, error) {
	stmt := `select count(*) from "users" where email = ?`

	var count int
	err := m.DB.QueryRow(stmt, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

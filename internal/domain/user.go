package domain

import (
	"errors"
	"time"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("user name is required")
	}
	if u.Email == "" {
		return errors.New("user email is required")
	}
	return nil
}

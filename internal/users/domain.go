package users

import (
	"errors"
	"time"
)

type User struct {
	ID        int64
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
}

var ErrEmailAlreadyExists = errors.New("email already exists")

package users

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

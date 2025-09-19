package users

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, name, email, password string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, name, email, password string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	newUser := &User{
		Email:    email,
		Name:     name,
		Password: string(hashed),
	}
	created, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	// Do not return password to callers
	created.Password = ""
	return created, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get by email: %w", err)
	}
	if u != nil {
		u.Password = ""
	}
	return u, nil
}

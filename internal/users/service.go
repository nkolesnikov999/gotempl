package users

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, name, email, password string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Authenticate(ctx context.Context, email, password string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, name, email, password string) (*User, error) {
	// Pre-check for email uniqueness to return a friendly error before hashing
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check existing email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

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

func (s *service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("authenticate: get by email: %w", err)
	}
	if u == nil {
		return nil, fmt.Errorf("authenticate: invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("authenticate: invalid credentials")
	}
	u.Password = ""
	return u, nil
}

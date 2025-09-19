package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRepository struct {
	pool *pgxpool.Pool
}

func NewPgxRepository(pool *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{pool: pool}
}

func (r *PgxRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	const query = `
        INSERT INTO users (email, name, password)
        VALUES ($1, $2, $3)
        RETURNING id, email, name, password, createdat
    `

	row := r.pool.QueryRow(ctx, query, user.Email, user.Name, user.Password)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt); err != nil {
		return nil, fmt.Errorf("scan inserted user: %w", err)
	}
	return &u, nil
}

func (r *PgxRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
        SELECT id, email, name, password, createdat
        FROM users
        WHERE email = $1
        LIMIT 1
    `

	row := r.pool.QueryRow(ctx, query, email)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by email: %w", err)
	}
	return &u, nil
}

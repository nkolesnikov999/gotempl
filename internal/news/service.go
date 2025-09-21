package news

import (
	"context"
)

type Service interface {
	Search(ctx context.Context, params SearchParams) (SearchResult, error)
	GetByAlias(ctx context.Context, alias string) (*Item, error)
	Popular(ctx context.Context, limit int) ([]Item, error)
	Latest(ctx context.Context, limit int) ([]Item, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Search(ctx context.Context, params SearchParams) (SearchResult, error) {
	return s.repo.Search(ctx, params)
}

func (s *service) GetByAlias(ctx context.Context, alias string) (*Item, error) {
	return s.repo.GetByAlias(ctx, alias)
}

func (s *service) Popular(ctx context.Context, limit int) ([]Item, error) {
	return s.repo.GetPopular(ctx, limit)
}
func (s *service) Latest(ctx context.Context, limit int) ([]Item, error) {
	return s.repo.GetLatest(ctx, limit)
}

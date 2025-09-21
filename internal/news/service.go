package news

import (
	"context"
	"fmt"
)

type Service interface {
	Search(ctx context.Context, params SearchParams) (SearchResult, error)
	GetByAlias(ctx context.Context, alias string) (*Item, error)
	Popular(ctx context.Context, limit int) ([]Item, error)
	Latest(ctx context.Context, limit int) ([]Item, error)
	Create(ctx context.Context, title, preview, content string, userID int64) (*Item, error)
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

func (s *service) Create(ctx context.Context, title, preview, content string, userID int64) (*Item, error) {
	if title == "" || content == "" {
		return nil, fmt.Errorf("title and content are required")
	}
	alias := slugify(title)
	return s.repo.Create(ctx, title, preview, content, userID, alias, nil)
}

func slugify(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r == ' ' || r == '\t' || r == '\n':
			out = append(out, '-')
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_':
			out = append(out, r)
		case (r >= 'A' && r <= 'Z'):
			out = append(out, r+32)
		default:
			// skip others
		}
	}
	if len(out) == 0 {
		return "news"
	}
	return string(out)
}

package news

import "context"

type Repository interface {
	Search(ctx context.Context, params SearchParams) (SearchResult, error)
	GetByAlias(ctx context.Context, alias string) (*Item, error)
	GetPopular(ctx context.Context, limit int) ([]Item, error)
	GetLatest(ctx context.Context, limit int) ([]Item, error)
	Create(ctx context.Context, title, preview, content string, userID int64, alias string, categories []int64) (*Item, error)
}

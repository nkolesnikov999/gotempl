package news

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRepository struct{ pool *pgxpool.Pool }

func NewPgxRepository(pool *pgxpool.Pool) *PgxRepository { return &PgxRepository{pool: pool} }

func (r *PgxRepository) Search(ctx context.Context, p SearchParams) (SearchResult, error) {
	if p.Limit <= 0 {
		p.Limit = 12
	}
	if p.Offset < 0 {
		p.Offset = 0
	}

	var whereParts []string
	var args []any
	arg := 1

	if p.Keyword != "" {
		whereParts = append(whereParts, fmt.Sprintf("(n.title ILIKE '%%' || $%d || '%%' OR n.content ILIKE '%%' || $%d || '%%')", arg, arg))
		args = append(args, p.Keyword)
		arg++
	}
	if p.Category != "" {
		whereParts = append(whereParts, fmt.Sprintf("EXISTS (SELECT 1 FROM news_categories nc JOIN categories c ON c.id = nc.category_id WHERE nc.news_id = n.id AND c.alias = $%d)", arg))
		args = append(args, p.Category)
		arg++
	}
	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	orderBy := "n.createdat DESC"
	switch p.Sort {
	case "popular":
		orderBy = "n.views DESC, n.createdat DESC"
	}

	countSQL := "SELECT COUNT(*) FROM news n " + whereSQL
	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return SearchResult{}, fmt.Errorf("count news: %w", err)
	}

	listSQL := "SELECT n.id, n.title, n.alias, n.preview, n.content, n.user_id, n.createdat, n.views, u.name FROM news n JOIN users u ON u.id = n.user_id " + whereSQL + " ORDER BY " + orderBy + " LIMIT $" + fmt.Sprint(arg) + " OFFSET $" + fmt.Sprint(arg+1)
	args = append(args, p.Limit, p.Offset)

	rows, err := r.pool.Query(ctx, listSQL, args...)
	if err != nil {
		return SearchResult{}, fmt.Errorf("query news: %w", err)
	}
	defer rows.Close()

	items, err := r.scanItems(ctx, rows)
	if err != nil {
		return SearchResult{}, err
	}

	return SearchResult{Total: total, Items: items}, nil
}

func (r *PgxRepository) GetByAlias(ctx context.Context, alias string) (*Item, error) {
	const q = `
        SELECT n.id, n.title, n.alias, n.preview, n.content, n.user_id, n.createdat, n.views, u.name
        FROM news n JOIN users u ON u.id = n.user_id
        WHERE n.alias = $1
        LIMIT 1
    `
	row := r.pool.QueryRow(ctx, q, alias)
	var it Item
	if err := row.Scan(&it.ID, &it.Title, &it.Alias, &it.Preview, &it.Content, &it.AuthorID, &it.CreatedAt, &it.Views, &it.Author); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get news by alias: %w", err)
	}
	it.Categories, _ = r.getCategories(ctx, it.ID)
	return &it, nil
}

func (r *PgxRepository) GetPopular(ctx context.Context, limit int) ([]Item, error) {
	if limit <= 0 {
		limit = 6
	}
	const q = `
        SELECT n.id, n.title, n.alias, n.preview, n.content, n.user_id, n.createdat, n.views, u.name
        FROM news n JOIN users u ON u.id = n.user_id
        ORDER BY n.views DESC, n.createdat DESC
        LIMIT $1
    `
	rows, err := r.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("popular: %w", err)
	}
	defer rows.Close()
	return r.scanItems(ctx, rows)
}

func (r *PgxRepository) GetLatest(ctx context.Context, limit int) ([]Item, error) {
	if limit <= 0 {
		limit = 6
	}
	const q = `
        SELECT n.id, n.title, n.alias, n.preview, n.content, n.user_id, n.createdat, n.views, u.name
        FROM news n JOIN users u ON u.id = n.user_id
        ORDER BY n.createdat DESC
        LIMIT $1
    `
	rows, err := r.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("latest: %w", err)
	}
	defer rows.Close()
	return r.scanItems(ctx, rows)
}

func (r *PgxRepository) Create(ctx context.Context, title, preview, content string, userID int64, alias string, categories []int64) (*Item, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const ins = `INSERT INTO news (title, alias, preview, content, user_id)
                 VALUES ($1,$2,$3,$4,$5)
                 RETURNING id, title, alias, preview, content, user_id, createdat, views`
	var it Item
	if err := tx.QueryRow(ctx, ins, title, alias, preview, content, userID).
		Scan(&it.ID, &it.Title, &it.Alias, &it.Preview, &it.Content, &it.AuthorID, &it.CreatedAt, &it.Views); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("insert news pg: %s %s", pgErr.Code, pgErr.Message)
		}
		return nil, fmt.Errorf("insert news: %w", err)
	}

	for _, cid := range categories {
		if _, err := tx.Exec(ctx, `INSERT INTO news_categories(news_id, category_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, it.ID, cid); err != nil {
			return nil, fmt.Errorf("link category: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	it.Categories, _ = r.getCategories(ctx, it.ID)
	return &it, nil
}

func (r *PgxRepository) scanItems(ctx context.Context, rows pgx.Rows) ([]Item, error) {
	items := make([]Item, 0)
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Title, &it.Alias, &it.Preview, &it.Content, &it.AuthorID, &it.CreatedAt, &it.Views, &it.Author); err != nil {
			return nil, fmt.Errorf("scan news: %w", err)
		}
		cats, err := r.getCategories(ctx, it.ID)
		if err == nil {
			it.Categories = cats
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *PgxRepository) getCategories(ctx context.Context, newsID int64) ([]Category, error) {
	const q = `
        SELECT c.id, c.title, c.alias, c.createdat
        FROM news_categories nc JOIN categories c ON c.id = nc.category_id
        WHERE nc.news_id = $1
        ORDER BY c.title
    `
	rows, err := r.pool.Query(ctx, q, newsID)
	if err != nil {
		return nil, fmt.Errorf("news categories: %w", err)
	}
	defer rows.Close()
	var cats []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Title, &c.Alias, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		cats = append(cats, c)
	}
	return cats, nil
}

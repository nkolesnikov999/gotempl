package news

import "time"

type Category struct {
	ID        int64
	Title     string
	Alias     string
	CreatedAt time.Time
}

type Item struct {
	ID         int64
	Title      string
	Alias      string
	Preview    string
	Content    string
	AuthorID   int64
	Author     string
	CreatedAt  time.Time
	Views      int64
	Categories []Category
}

type SearchParams struct {
	Limit    int
	Offset   int
	Category string // alias
	Keyword  string
	Sort     string // "latest" | "popular"
}

type SearchResult struct {
	Total int64
	Items []Item
}

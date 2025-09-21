package views

import "nkpro/gotempl/views/components"

type NewsCard struct {
	Image       string
	Title       string
	Description string
	Alias       string
	User        components.CardUser
}

type CategoryPageProps struct {
	User     PageUser
	Category string
	Query    string
	Items    []NewsCard
	Page     int
	PageSize int
	Total    int
}

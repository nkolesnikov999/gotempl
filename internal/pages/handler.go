package pages

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type HomeHandler struct {
	router fiber.Router
}

type Category struct {
	Name string
}

func NewHandler(router fiber.Router) {
	h := &HomeHandler{
		router: router,
	}
	api := h.router.Group("/pages")
	api.Get("/", h.home)
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	// Log the incoming request
	slog.Info("Home page accessed",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.String("user_agent", c.Get("User-Agent")),
		slog.String("ip", c.IP()),
	)

	// Define categories for the menu
	categories := []Category{
		{Name: "Еда"},
		{Name: "Напитки"},
		{Name: "Машины"},
		{Name: "Одежда"},
		{Name: "Дом"},
		{Name: "Спорт"},
		{Name: "Развлечения"},
		{Name: "Другое"},
	}

	// Log successful response
	slog.Info("Home page response sent",
		slog.String("status", "200"),
		slog.String("template", "page"),
		slog.Int("categories_count", len(categories)),
	)

	// Debug logging
	for i, cat := range categories {
		slog.Debug("Category data",
			slog.Int("index", i),
			slog.String("name", cat.Name),
		)
	}

	// Try simple approach - just pass the slice directly
	return c.Render("page", categories)
}

package pages

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type HomeHandler struct {
	router fiber.Router
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

	// Log successful response
	slog.Info("Home page response sent",
		slog.String("status", "200"),
		slog.String("response", "Hi"),
	)

	return c.SendString("Hi")
}

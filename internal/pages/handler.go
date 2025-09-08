package pages

import (
	"log/slog"

	"nkpro/gotempl/pkg/tadapter"
	"nkpro/gotempl/views"

	"github.com/gofiber/fiber/v2"
)

type HomeHandler struct {
	router fiber.Router
}

func NewHandler(router fiber.Router) {
	h := &HomeHandler{
		router: router,
	}
	h.router.Get("/", h.home)
	h.router.Get("/404", h.notFound)
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

	component := views.Main()
	return tadapter.Render(c, component)
}

func (h *HomeHandler) notFound(c *fiber.Ctx) error {
	slog.Warn("Not found",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.String("user_agent", c.Get("User-Agent")),
		slog.String("ip", c.IP()),
		slog.String("status", "404"),
	)
	return c.SendStatus(fiber.StatusNotFound)
}

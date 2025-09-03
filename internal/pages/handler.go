package pages

import (
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
	return c.SendString("Hi")
}

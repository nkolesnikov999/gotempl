package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slog.Debug("session: get start", slog.String("path", c.Path()), slog.String("method", c.Method()))
		sess, err := store.Get(c)
		if err != nil {
			slog.Error("session: get failed", slog.String("path", c.Path()), slog.String("method", c.Method()), slog.String("error", err.Error()))
			c.Locals("email", "")
			return c.Next()
		}
		sessID := sess.ID()
		slog.Debug("session: get ok", slog.String("id", sessID))

		userEmail := ""
		if email, ok := sess.Get("email").(string); ok {
			userEmail = email
			slog.Debug("session: email present", slog.String("id", sessID), slog.String("email", userEmail))
		} else {
			slog.Debug("session: email missing", slog.String("id", sessID))
		}
		c.Locals("email", userEmail)
		return c.Next()
	}
}

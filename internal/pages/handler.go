package pages

import (
	"log/slog"

	"nkpro/gotempl/pkg/tadapter"
	"nkpro/gotempl/views"
	"nkpro/gotempl/views/widgets"

	"github.com/a-h/templ"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
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
	h.router.Get("/register", h.register)
	h.router.Get("/404", h.notFound)

	// API
	h.router.Post("/api/register", h.apiRegister)
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

func (h *HomeHandler) register(c *fiber.Ctx) error {
	component := views.Register()
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

// --- API ---

type registerForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (h *HomeHandler) apiRegister(c *fiber.Ctx) error {
	var form registerForm
	if err := c.BodyParser(&form); err != nil {
		return htmxRender(c.Status(fiber.StatusBadRequest), widgets.RegisterResult(false, "Некорректные данные формы"))
	}

	errs := validate.Validate(
		&validators.StringIsPresent{Field: form.Name, Name: "Имя"},
		&validators.EmailIsPresent{Field: form.Email, Name: "Email"},
		&validators.StringLengthInRange{Field: form.Password, Name: "Пароль", Min: 6, Max: 100},
	)
	if errs.HasAny() {
		return htmxRender(c.Status(fiber.StatusUnprocessableEntity), widgets.RegisterResult(false, errs.Error()))
	}

	return htmxRender(c.Status(fiber.StatusOK), widgets.RegisterResult(true, "Регистрация прошла успешно"))
}

// htmxRender renders a small component suitable for HTMX swaps
func htmxRender(c *fiber.Ctx, comp templ.Component) error {
	return tadapter.Render(c, comp)
}

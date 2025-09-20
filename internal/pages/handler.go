package pages

import (
	"log/slog"

	"nkpro/gotempl/internal/users"
	"nkpro/gotempl/pkg/tadapter"
	"nkpro/gotempl/views"
	"nkpro/gotempl/views/widgets"

	"github.com/a-h/templ"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type HomeHandler struct {
	router      fiber.Router
	store       *session.Store
	userService users.Service
}

type Option func(*HomeHandler)

func WithUserService(s users.Service) Option { return func(h *HomeHandler) { h.userService = s } }

func NewHandler(router fiber.Router, store *session.Store, opts ...Option) {
	h := &HomeHandler{
		router: router,
		store:  store,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.router.Get("/", h.home)
	h.router.Get("/register", h.register)
	h.router.Get("/login", h.login)
	h.router.Get("/404", h.notFound)

	// API
	h.router.Post("/api/register", h.apiRegister)
	h.router.Post("/api/login", h.apiLogin)
	h.router.Get("/logout", h.logout)
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

	user := h.currentUser(c)
	component := views.Main(views.MainProps{User: user})
	return tadapter.Render(c, component)
}

func (h *HomeHandler) register(c *fiber.Ctx) error {
	user := h.currentUser(c)
	component := views.Register(views.RegisterProps{User: user})
	return tadapter.Render(c, component)
}

func (h *HomeHandler) login(c *fiber.Ctx) error {
	user := h.currentUser(c)
	component := views.Login(views.LoginProps{User: user})
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

	if h.userService == nil {
		return htmxRender(c.Status(fiber.StatusInternalServerError), widgets.RegisterResult(false, "Сервис пользователей не инициализирован"))
	}

	u, err := h.userService.Register(c.UserContext(), form.Name, form.Email, form.Password)
	if err != nil {
		slog.Error("register failed", slog.String("error", err.Error()))
		return htmxRender(c.Status(fiber.StatusInternalServerError), widgets.RegisterResult(false, "Ошибка регистрации"))
	}

	// Save email in session after successful registration
	slog.Debug("session: register get start", slog.String("path", c.Path()))
	sess, err := h.store.Get(c)
	if err == nil {
		sid := sess.ID()
		slog.Debug("session: register get ok", slog.String("id", sid))
		sess.Set("email", u.Email)
		if err := sess.Save(); err != nil {
			slog.Error("session: register save failed", slog.String("id", sid), slog.String("error", err.Error()))
		} else {
			slog.Info("session: register save ok", slog.String("id", sid), slog.String("email", u.Email))
		}
	} else {
		slog.Error("session: register get failed", slog.String("error", err.Error()))
	}

	// HTMX redirect to home on success
	c.Set("HX-Redirect", "/")
	c.Set("HX-Location", "/")
	return c.SendStatus(fiber.StatusNoContent)
}

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (h *HomeHandler) apiLogin(c *fiber.Ctx) error {
	var form loginForm
	if err := c.BodyParser(&form); err != nil {
		return htmxRender(c.Status(fiber.StatusBadRequest), widgets.LoginResult(false, "Некорректные данные формы"))
	}

	errs := validate.Validate(
		&validators.EmailIsPresent{Field: form.Email, Name: "Email"},
		&validators.StringLengthInRange{Field: form.Password, Name: "Пароль", Min: 6, Max: 100},
	)
	if errs.HasAny() {
		return htmxRender(c.Status(fiber.StatusUnprocessableEntity), widgets.LoginResult(false, errs.Error()))
	}

	if h.userService == nil {
		return htmxRender(c.Status(fiber.StatusInternalServerError), widgets.LoginResult(false, "Сервис пользователей не инициализирован"))
	}

	u, err := h.userService.Authenticate(c.UserContext(), form.Email, form.Password)
	if err != nil {
		slog.Warn("login failed", slog.String("error", err.Error()))
		return htmxRender(c.Status(fiber.StatusUnauthorized), widgets.LoginResult(false, "Неверный email или пароль"))
	}

	slog.Debug("session: login get start", slog.String("path", c.Path()))
	sess, err := h.store.Get(c)
	if err != nil {
		slog.Error("session: login get failed", slog.String("error", err.Error()))
		return htmxRender(c.Status(fiber.StatusInternalServerError), widgets.LoginResult(false, "Ошибка сессии"))
	}
	sid := sess.ID()
	slog.Debug("session: login get ok", slog.String("id", sid))
	sess.Set("email", u.Email)
	if err := sess.Save(); err != nil {
		slog.Error("session: login save failed", slog.String("id", sid), slog.String("error", err.Error()))
		return htmxRender(c.Status(fiber.StatusInternalServerError), widgets.LoginResult(false, "Ошибка сессии"))
	}
	slog.Info("session: login save ok", slog.String("id", sid), slog.String("email", u.Email))

	// HTMX redirect to home on success
	c.Set("HX-Redirect", "/")
	c.Set("HX-Location", "/")
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *HomeHandler) logout(c *fiber.Ctx) error {
	slog.Debug("session: logout get start", slog.String("path", c.Path()))
	sess, err := h.store.Get(c)
	if err == nil {
		sid := sess.ID()
		if err := sess.Destroy(); err != nil {
			slog.Error("session: logout destroy failed", slog.String("id", sid), slog.String("error", err.Error()))
		} else {
			slog.Info("session: logout destroy ok", slog.String("id", sid))
		}
	} else {
		slog.Error("session: logout get failed", slog.String("error", err.Error()))
	}
	return c.Redirect("/", fiber.StatusSeeOther)
}

// htmxRender renders a small component suitable for HTMX swaps
func htmxRender(c *fiber.Ctx, comp templ.Component) error {
	return tadapter.Render(c, comp)
}

func (h *HomeHandler) currentUser(c *fiber.Ctx) views.PageUser {
	email := c.Locals("email").(string)
	if email == "" || h.userService == nil {
		return views.PageUser{}
	}
	u, err := h.userService.GetByEmail(c.UserContext(), email)
	if err != nil || u == nil {
		return views.PageUser{}
	}
	return views.PageUser{Email: u.Email, Name: u.Name}
}

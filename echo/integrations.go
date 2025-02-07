package echo

import (
	"github.com/DBenyukh/goform/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	csrfTokenCookieName = "csrf_token" // Имя cookie для CSRF-токена
)

// FormMiddleware возвращает middleware для автоматической привязки данных.
func FormMiddleware(model interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			form := core.NewForm(model)
			if err := form.Bind(c.Request()); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
			}
			c.Set("form", form)
			return next(c)
		}
	}
}

// CSRFMiddleware возвращает middleware для проверки CSRF-токена.
func CSRFMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получаем CSRF-токен из cookies
			expectedToken, err := c.Cookie(csrfTokenCookieName)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, "CSRF token not found")
			}

			// Получаем CSRF-токен из запроса (например, из формы или заголовка)
			receivedToken := c.FormValue("csrf_token")
			if receivedToken == "" {
				return echo.NewHTTPError(http.StatusForbidden, "CSRF token is required")
			}

			// Сравниваем токены
			if !isValidCSRFToken(receivedToken, expectedToken.Value) {
				return echo.NewHTTPError(http.StatusForbidden, "Invalid CSRF token")
			}

			return next(c)
		}
	}
}

// isValidCSRFToken проверяет, что переданный токен совпадает с ожидаемым.
func isValidCSRFToken(receivedToken, expectedToken string) bool {
	return receivedToken == expectedToken
}

// SetCSRFToken устанавливает CSRF-токен в cookies.
func SetCSRFToken(c echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = csrfTokenCookieName
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	c.SetCookie(cookie)
}

// RenderForm рендерит форму в контексте Echo.
func RenderForm(c echo.Context, form *core.Form) error {
	return c.Render(http.StatusOK, "templates/default.html", map[string]interface{}{
		"form": form.Render(),
	})
}

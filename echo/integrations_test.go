package echo

import (
	"github.com/DBenyukh/goform/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// RegistrationForm представляет форму регистрации.
type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
}

func TestFormMiddleware(t *testing.T) {
	e := echo.New()
	model := &RegistrationForm{}

	e.Use(FormMiddleware(model))
	e.GET("/", func(c echo.Context) error {
		form := c.Get("form").(*core.Form)
		return c.JSON(http.StatusOK, form)
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}
}

func TestCSRFMiddleware(t *testing.T) {
	e := echo.New()

	e.Use(CSRFMiddleware())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status code 403, got %d", rec.Code)
	}
}

// Тестируем рендеринг формы
func TestRenderForm(t *testing.T) {
	// Создаем новый экземпляр echo
	e := echo.New()

	// Инициализируем рендерер шаблонов
	templateRenderer := core.NewTemplateRenderer()

	// Регистрируем рендерер шаблонов
	e.Renderer = templateRenderer

	// Создаем новый запрос
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Устанавливаем контекст запроса
	c := e.NewContext(req, rec)

	// Пытаемся отрендерить шаблон
	err := templateRenderer.Render(rec, "default.html", nil, c) // Передаем имя шаблона без пути
	if err != nil {
		t.Errorf("Error rendering template: %v", err)
	}

	// Проверяем статус код и другие параметры
	assert.Equal(t, http.StatusOK, rec.Code) // Проверяем, что статус код 200

	// Логируем результат для отладки
	//log.Printf("Response body: %s", rec.Body.String())
}

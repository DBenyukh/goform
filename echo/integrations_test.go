package echo

import (
	"encoding/json"
	"errors"
	"github.com/DBenyukh/goform/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

type TestForm struct {
	Username string `form:"username" validate:"required,min=3" validate_msg:"Username must be at least 3 characters"`
	Email    string `form:"email" validate:"required,email" validate_msg:"Please provide a valid email address"`
	Password string `form:"password" validate:"required" validate_msg:"Password is required"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
}

// TestFormMiddlewareSuccess проверяет успешную привязку данных формы.
func TestFormMiddlewareSuccess(t *testing.T) {
	e := echo.New()
	model := &TestForm{
		Method: "POST",
		FormID: "test_form",
	}

	e.Use(FormMiddleware(model, model.Method, model.FormID))
	e.POST("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model)
	})

	formData := url.Values{
		"test_form_username": {"testuser"},
		"test_form_email":    {"test@example.com"},
		"test_form_password": {"password123"},
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var responseData TestForm
	err := json.Unmarshal(rec.Body.Bytes(), &responseData)
	assert.NoError(t, err)

	assert.Equal(t, "testuser", responseData.Username)
	assert.Equal(t, "test@example.com", responseData.Email)
	assert.Equal(t, "password123", responseData.Password)
}

// TestCSRFMiddleware проверяет корректность работы CSRFMiddleware.
func TestCSRFMiddleware(t *testing.T) {
	e := echo.New()

	e.Use(CSRFMiddleware())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

// TestRenderFormSuccess проверяет успешный рендеринг формы.
func TestRenderFormSuccess(t *testing.T) {
	e := echo.New()

	// Получаем абсолютный путь к директории templates
	projectDir, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("Failed to get project directory: %v", err)
	}
	templateDir := filepath.Join(projectDir, "templates")

	// Инициализируем рендерер
	renderer, err := core.NewTemplateRenderer(templateDir)
	if err != nil {
		t.Fatalf("Failed to create template renderer: %v", err)
	}

	// Устанавливаем рендерер
	e.Renderer = renderer

	model := TestForm{
		Method: "POST",
		FormID: "test_form",
	}
	form := core.NewForm(model, model.Method, model.FormID)
	form.RenderHTML = true // Устанавливаем флаг (true для монолита, false для разделённого фронта/бэка)

	// Заполняем форму данными
	form.Fields[0].Value = "testuser"         // Username
	form.Fields[1].Value = "test@example.com" // Email
	form.Fields[2].Value = "password123"      // Password

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = RenderForm(c, form)
	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), "test_form_username")
	assert.Contains(t, rec.Body.String(), "test_form_email")
	assert.Contains(t, rec.Body.String(), "test_form_password")
}

// TestAddCustomValidationMiddleware проверяет добавление кастомных правил валидации.
func TestAddCustomValidationMiddleware(t *testing.T) {
	e := echo.New()
	model := &TestForm{ // Используем указатель на модель
		Method: "POST",
		FormID: "test_form",
	}

	e.Use(FormMiddleware(model, model.Method, model.FormID))
	e.Use(AddCustomValidationMiddleware("username", func(value string) error {
		if len(value) < 5 {
			return errors.New("username too short")
		}
		return nil
	}))
	e.POST("/", func(c echo.Context) error {
		form := c.Get("form").(*core.Form)

		// Находим поле "username" и применяем кастомную валидацию
		for _, field := range form.Fields {
			if field.Name == "username" && field.CustomValidation != nil {
				if err := field.CustomValidation(field.Value.(string)); err != nil {
					form.AddError(field.Name, err.Error()) // Используем метод AddError
					responseData := form.ToResponse()      // Преобразуем Form в FormResponse
					return c.JSON(http.StatusBadRequest, responseData)
				}
			}
		}

		// Если кастомная валидация прошла успешно, возвращаем OK
		return c.NoContent(http.StatusOK)
	})

	// Тест с коротким именем
	formData := url.Values{
		"test_form_username": {"us"},
		"test_form_email":    {"testexample.com"},
		"test_form_password": {"password"},
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Тест с валидным именем
	formData.Set("test_form_username", "longusername")
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

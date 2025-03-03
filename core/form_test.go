package core

import (
	"net/http/httptest"
	"testing"
)

type TestForm struct {
	Username string `form:"username" validate:"required,min=3" validate_msg:"Username must be at least 3 characters"`
	Email    string `form:"email" validate:"required,email" validate_msg:"Please provide a valid email address"`
	Password string `form:"password" validate:"required,min=6" validate_msg:"Password is required"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
}

func TestNewForm(t *testing.T) {
	model := &TestForm{
		Method: "POST",
		FormID: "test_form",
	}
	form := NewForm(model, model.Method, model.FormID)

	if len(form.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(form.Fields))
	}

	if form.Method != "POST" {
		t.Errorf("Expected method 'POST', got '%s'", form.Method)
	}

	if form.FormID != "test_form" {
		t.Errorf("Expected form ID 'test_form', got '%s'", form.FormID)
	}
}

func TestFormRender(t *testing.T) {
	model := &TestForm{
		Method: "POST",
		FormID: "test_form",
	}
	form := NewForm(model, model.Method, model.FormID)
	form.RenderHTML = true // Устанавливаем флаг (true для монолита, false для разделённого фронта/бэка)

	responseData := form.ToResponse() // Используем ToResponse вместо Render
	// Приводим тип к FormResponse
	formResponse, ok := responseData.(FormResponse)
	if !ok {
		t.Fatal("Expected FormResponse, got different type")
	}

	// Проверяем, что структура FormResponse содержит поля формы
	if len(formResponse.Fields) == 0 {
		t.Error("Expected non-empty fields, got empty slice")
	}

	// Проверяем, что каждое поле имеет имя и тип
	for _, field := range formResponse.Fields {
		if field.Name == "" {
			t.Error("Expected field name, got empty string")
		}
		if field.Type == "" {
			t.Error("Expected field type, got empty string")
		}
	}
}

func TestFormBind(t *testing.T) {
	model := &TestForm{
		Method: "POST",
		FormID: "test_form",
	}
	form := NewForm(model, model.Method, model.FormID)

	req := httptest.NewRequest("POST", "/", nil)
	req.Form = map[string][]string{
		"test_form_username": {"testuser"},
		"test_form_email":    {"test@example.com"},
		"test_form_password": {"password123"},
	}

	if err := form.Bind(req); err != nil {
		t.Errorf("Bind failed: %v", err)
	}

	// Проверяем значение поля username
	if form.Fields[0].Value != "testuser" {
		t.Errorf("Expected username 'testuser', got '%v'", form.Fields[0].Value)
	}
}

func TestFormValidate(t *testing.T) {
	model := &TestForm{
		Method: "POST",
		FormID: "test_form",
	}
	form := NewForm(model, model.Method, model.FormID)

	// Устанавливаем значения полей
	form.Fields[0].Value = "ab" // Не соответствует min=3
	form.Fields[1].Value = "invalid-email"
	form.Fields[2].Value = "12345" // Не соответствует min=6

	// Выполняем валидацию
	if err := form.Validate(model); err == nil {
		t.Error("Expected validation error, got nil")
	}

	// Проверяем, что ошибка была установлена для всех полей
	for i, field := range form.Fields {
		if field.Error == "" {
			t.Errorf("Expected error for field %d, got none", i)
		}
	}

	// Проверяем текст ошибки для поля username
	expectedUsernameError := "Username must be at least 3 characters"
	if form.Fields[0].Error != expectedUsernameError {
		t.Errorf("Expected error message '%s', got '%s'", expectedUsernameError, form.Fields[0].Error)
	}

	// Проверяем текст ошибки для поля email
	expectedEmailError := "Please provide a valid email address"
	if form.Fields[1].Error != expectedEmailError {
		t.Errorf("Expected error message '%s', got '%s'", expectedEmailError, form.Fields[1].Error)
	}

	// Проверяем текст ошибки для поля password
	expectedPasswordError := "Password is required"
	if form.Fields[2].Error != expectedPasswordError {
		t.Errorf("Expected error message '%s', got '%s'", expectedPasswordError, form.Fields[2].Error)
	}
}

func TestFormToJSON(t *testing.T) {
	model := &TestForm{
		Method: "POST",
		FormID: "test_form",
	}
	form := NewForm(model, model.Method, model.FormID)

	// Устанавливаем значения полей
	form.Fields[0].Value = "testuser"
	form.Fields[1].Value = "test@example.com"
	form.Fields[2].Value = "password123"

	// Получаем JSON
	jsonData := form.ToJSONResponse()

	// Проверяем, что данные корректны
	if jsonData["username"].(map[string]interface{})["value"] != "testuser" {
		t.Errorf("Expected username 'testuser', got '%v'", jsonData["username"].(map[string]interface{})["value"])
	}
}

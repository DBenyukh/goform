package core

import (
	"net/http/httptest"
	"testing"
)

type TestForm struct {
	Username string `form:"username" validate:"required,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
}

func TestNewForm(t *testing.T) {
	model := &TestForm{}
	form := NewForm(model)

	if len(form.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(form.Fields))
	}
}

func TestFormRender(t *testing.T) {
	model := &TestForm{}
	form := NewForm(model)

	renderData := form.Render()

	// Проверяем, что структура RenderData содержит поля формы
	if len(renderData.Fields) == 0 {
		t.Error("Expected non-empty fields, got empty slice")
	}

	// Проверяем, что каждое поле имеет имя и тип
	for _, field := range renderData.Fields {
		if field.Name == "" {
			t.Error("Expected field name, got empty string")
		}
		if field.Type == "" {
			t.Error("Expected field type, got empty string")
		}
	}
}

func TestFormBind(t *testing.T) {
	model := &TestForm{}
	form := NewForm(model)

	req := httptest.NewRequest("POST", "/", nil)
	req.Form = map[string][]string{
		"username": {"testuser"},
		"email":    {"test@example.com"},
		"password": {"password123"},
	}

	if err := form.Bind(req); err != nil {
		t.Errorf("Bind failed: %v", err)
	}

	if form.Fields[0].Value != "testuser" {
		t.Errorf("Expected username 'testuser', got '%v'", form.Fields[0].Value)
	}
}

func TestFormValidate(t *testing.T) {
	model := &TestForm{}
	form := NewForm(model)

	form.Fields[0].Value = "ab" // Не соответствует min=3
	form.Fields[1].Value = "invalid-email"
	form.Fields[2].Value = "12345" // Не соответствует min=6

	if err := form.Validate(model); err == nil {
		t.Error("Expected validation error, got nil")
	}

	if form.Fields[0].Error == "" {
		t.Error("Expected error for username, got none")
	}
}

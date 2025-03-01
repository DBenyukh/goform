package main

import (
	"encoding/json"
	"errors"
	"github.com/DBenyukh/goform/core"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// RegistrationForm представляет форму регистрации.
type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3" validate_msg:"Username must be at least 3 characters"`
	Email    string `form:"email" validate:"required,email" validate_msg:"Please provide a valid email address"`
	Password string `form:"password" validate:"required" validate_msg:"Password is required"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
}

var tmpl *template.Template

func init() {
	projectDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("Error getting absolute project directory path: %v", err)
	}

	templateDir := filepath.Join(projectDir, "templates")
	renderer, err := core.NewTemplateRenderer(templateDir)
	if err != nil {
		log.Fatalf("Failed to create template renderer: %v", err)
	}

	tmpl = renderer.Templates
}

func isPasswordStrong(password string) error {
	// Проверка минимальной длины пароля
	if len(password) < 6 {
		return errors.New("Password must be at least 6 characters long")
	}

	// Проверка наличия специальных символов
	specialChars := "!@#$%^&*"
	hasSpecialChar := false
	for _, char := range specialChars {
		if strings.ContainsRune(password, char) {
			hasSpecialChar = true
			break
		}
	}

	if !hasSpecialChar {
		return errors.New("Password must contain at least one special character (!@#$%^&*)")
	}

	return nil
}

func isAjax(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		model := &RegistrationForm{
			Method: "POST",
			FormID: "register_form",
		}
		form := core.NewForm(model, model.Method, model.FormID)
		form.RenderHTML = true // Устанавливаем флаг (true для монолита, false для разделённого фронта/бэка)

		// Добавление кастомного правила валидации
		form.AddCustomValidation("password", isPasswordStrong)

		if r.Method == http.MethodGet {
			token, err := core.GenerateCSRFToken()
			if err != nil {
				http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
				return
			}

			form.AddCSRFToken(token)

			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    token,
				HttpOnly: true,
				Path:     "/",
			})

			// Возвращаем данные в зависимости от флага
			if form.RenderHTML {
				_ = tmpl.Execute(w, form.ToResponse())
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(form.ToResponse())
			}
			return
		}

		method := r.Method
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			methodOverride := r.FormValue("_method")
			if methodOverride != "" {
				method = methodOverride
			}
		}

		if method == http.MethodPost {
			csrfTokenFromForm := r.FormValue(model.FormID + "_csrf_token")
			csrfTokenFromCookie, err := r.Cookie("csrf_token")
			if err != nil {
				http.Error(w, "CSRF token missing in cookies", http.StatusForbidden)
				return
			}

			if csrfTokenFromCookie.Value != csrfTokenFromForm {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		if err := form.Bind(r); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		for _, field := range form.Fields {
			field.Value = r.FormValue(model.FormID + "_" + field.Name)
		}

		model.Username = r.FormValue(model.FormID + "_username")
		model.Email = r.FormValue(model.FormID + "_email")
		model.Password = r.FormValue(model.FormID + "_password")

		if err := form.Validate(model); err != nil {
			if isAjax(r) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				errors := make(map[string]string)
				for _, field := range form.Fields {
					if field.Error != "" {
						errors[field.Name] = field.Error
					}
				}
				json.NewEncoder(w).Encode(map[string]interface{}{
					"errors": errors,
				})
				return
			}

			// Генерация нового CSRF-токена
			newToken, err := core.GenerateCSRFToken()
			if err != nil {
				http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
				return
			}

			// Обновление CSRF-токена в форме
			form.AddCSRFToken(newToken)

			// Установка нового CSRF-токена в куки
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    newToken,
				HttpOnly: true,
				Path:     "/",
			})

			// Рендеринг формы с новым CSRF-токеном
			_ = tmpl.Execute(w, form.ToResponse())
			return
		}

		if isAjax(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
			return
		}

		switch method {
		case http.MethodPost:
			w.Write([]byte("User registered successfully!"))
		case http.MethodPut:
			w.Write([]byte("User updated successfully!"))
		case http.MethodDelete:
			w.Write([]byte("User deleted successfully!"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}

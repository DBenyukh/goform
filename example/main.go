package main

import (
	"github.com/DBenyukh/goform/core"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// RegistrationForm представляет форму регистрации.
type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3" validate_msg:"Username must be at least 3 characters"`
	Email    string `form:"email" validate:"required,email" validate_msg:"Please provide a valid email address"`
	Password string `form:"password" validate:"required,min=6,max=8" validate_msg:"Password must be at least 6 and no more than 8 characters long"`
}

var tmpl *template.Template

func init() {
	// Получаем абсолютный путь к корневой директории проекта
	projectDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("Error getting absolute project directory path: %v", err)
	}

	// Формируем путь к шаблонам относительно корня проекта
	templateDir := filepath.Join(projectDir, "templates")

	// Загружаем шаблон с абсолютным путем
	tmpl = template.Must(template.ParseFiles(filepath.Join(templateDir, "default.html")))
}

func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		model := &RegistrationForm{}
		form := core.NewForm(model)

		if r.Method == http.MethodGet {
			// Генерация CSRF-токена
			token, err := core.GenerateCSRFToken()
			if err != nil {
				http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
				return
			}

			// Добавление CSRF-токена в форму
			form.AddCSRFToken(token)

			// Установка CSRF-токена в куки
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    token,
				HttpOnly: true,
				Path:     "/",
			})

			// Рендеринг формы
			_ = tmpl.Execute(w, form.Render())
			return
		}

		// Проверка CSRF-токена для POST-запроса
		if r.Method == http.MethodPost {
			// Получаем CSRF-токен из формы
			csrfTokenFromForm := r.FormValue("csrf_token")

			// Получаем CSRF-токен из куки
			csrfTokenFromCookie, err := r.Cookie("csrf_token")
			if err != nil {
				http.Error(w, "CSRF token missing in cookies", http.StatusForbidden)
				return
			}

			// Сравниваем токены
			if csrfTokenFromCookie.Value != csrfTokenFromForm {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		// Обработка данных формы
		if err := form.Bind(r); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		if err := form.Validate(model); err != nil {
			_ = tmpl.Execute(w, form.Render())
			return
		}

		w.Write([]byte("User registered successfully!"))
	})

	http.ListenAndServe(":8080", nil)
}

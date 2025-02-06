package main

import (
	"goform/core"
	"html/template"
	"net/http"
)

// RegistrationForm представляет форму регистрации.
type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
}

var tmpl *template.Template

func init() {
	// Загружаем шаблон
	tmpl = template.Must(template.ParseFiles("templates/default.html"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		model := &RegistrationForm{}
		form := core.NewForm(model)
		form.Render() // Рендерим форму
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		model := &RegistrationForm{}
		form := core.NewForm(model)

		if r.Method == http.MethodGet {
			_ = tmpl.Execute(w, form)
			return
		}

		if err := form.Bind(r); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		if err := form.Validate(model); err != nil {
			_ = tmpl.Execute(w, form)
			return
		}

		w.Write([]byte("User registered successfully!"))
	})

	http.ListenAndServe(":8080", nil)
}

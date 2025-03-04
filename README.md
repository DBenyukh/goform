# GoForm

[![Go Reference](https://pkg.go.dev/badge/github.com/DBenyukh/goform.svg)](https://pkg.go.dev/github.com/DBenyukh/goform)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/DBenyukh/goform/tree/master/LICENSE)

`GoForm` — это пакет для работы с HTML-формами в Go. Он предоставляет удобный API для создания, рендеринга, валидации и обработки форм. Пакет поддерживает рендеринг форм в HTML и JSON, кастомную валидацию, обработку AJAX-запросов и многое другое.

---

## Оглавление

1. [Установка](#установка)
2. [Быстрый старт](#быстрый-старт)
3. [Основные функции](#основные-функции)
    - [Создание формы](#создание-формы)
    - [Рендеринг формы](#рендеринг-формы)
    - [Валидация формы](#валидация-формы)
    - [Обработка AJAX-запросов](#обработка-ajax-запросов)
4. [Расширенные возможности](#расширенные-возможности)
    - [Кастомная валидация](#кастомная-валидация)
    - [Поддержка нескольких форм](#поддержка-нескольких-форм)
    - [Рендеринг HTML и JSON](#рендеринг-html-и-json)
    - [Интеграция с Echo](#интеграция-с-echo)
    - [Кастомные сообщения об ошибках](#кастомные-сообщения-об-ошибках)
    - [Скрытые поля](#скрытые-поля)
    - [CSRF-токены](#csrf-токены)
5. [Пример HTML-шаблона](#пример-html-шаблона)
6. [Примеры](#примеры)
   - [Пример 1: Простая форма регистрации](#пример-1-простая-форма-регистрации)
   - [Пример 2: Кастомная валидация](#пример-2-кастомная-валидация)
   - [Пример 3: Разделённый фронтенд и бэкенд с net/http](#пример-3-разделённый-фронтенд-и-бэкенд-с-nethttp)
   - [Пример 4: Разделённый фронтенд и бэкенд с Echo](#пример-4-разделённый-фронтенд-и-бэкенд-с-echo)
7. [Лицензия](#лицензия)

---

## Установка
Для установки пакета выполните команду:

```bash
go get github.com/DBenyukh/goform
```

## Быстрый старт
Вот пример создания и рендеринга простой формы:

```go
package main

import (
	"net/http"
	"github.com/DBenyukh/goform/core"
)

type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
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
	renderer, err := core.NewTemplateRenderer(templateDir, "")
	if err != nil {
		log.Fatalf("Failed to create template renderer: %v", err)
	}

	tmpl = renderer.Templates
}

func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		model := &RegistrationForm{
			Method: "POST",
			FormID: "register_form",
		}
		form := core.NewForm(model, model.Method, model.FormID)
		form.RenderHTML = true

		if r.Method == http.MethodGet {
			// Рендеринг формы
			response := form.ToResponse()
			if form.RenderHTML {
				// Рендеринг HTML
				_ = tmpl.ExecuteTemplate(w, "имя_шаблона.html", response)
			} else {
				// Возврат JSON
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
			return
		}

		// Обработка POST-запроса
		if err := form.Bind(r); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		if err := form.Validate(model); err != nil {
			// Возврат ошибок валидации
			response := form.ToResponse()
			if form.RenderHTML {
				_ = tmpl.ExecuteTemplate(w, "имя_шаблона.html", response)
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
			return
		}

		// Обработка успешной отправки формы
		w.Write([]byte("User registered successfully!"))
	})

	http.ListenAndServe(":8080", nil)
}
```

## Основные функции
### Создание формы
Для создания формы используйте функцию `NewForm`:
```go
form := core.NewForm(model, method, formID)
```

`model` — структура, описывающая поля формы.

`method` — HTTP-метод формы (например, "POST").

`formID` — уникальный идентификатор формы.

---

### Рендеринг формы
Для рендеринга формы используйте метод `ToResponse`:
```go
response := form.ToResponse()
```
Этот метод возвращает данные формы в зависимости от флага RenderHTML:

- Если `RenderHTML = true`, возвращается структура FormResponse для рендеринга HTML.
- Если `RenderHTML = false`, возвращается JSON.

Пример использования:
```go
if form.RenderHTML {
_ = tmpl.ExecuteTemplate(w, "имя_шаблона.html", response)
} else {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

---

### Валидация формы
Для валидации данных формы используйте метод `Validate`:
```go
if err := form.Validate(model); err != nil {
    // Обработка ошибок валидации
}
```

---

### Обработка AJAX-запросов
Библиотека автоматически определяет AJAX-запросы и возвращает данные в формате JSON:

```go
if isAjax(r) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(form.ToResponse())
}
```

---

## Расширенные возможности
### Кастомная валидация

Вы можете добавлять кастомные правила валидации:

```go
form.AddCustomValidation("password", func(value string) error {
    if len(value) < 6 {
        return errors.New("password must be at least 6 characters long")
    }
    return nil
})
```

---

### Поддержка нескольких форм
Библиотека поддерживает несколько форм на одной странице. Убедитесь, что у каждой формы уникальный `formID`.

---

### Рендеринг HTML и JSON
Вы можете управлять форматом вывода с помощью флага `RenderHTML`:

```go
form.RenderHTML = false // Возвращает JSON
form.RenderHTML = true  // Возвращает HTML
```

---

### Интеграция с Echo
Пакет поддерживает интеграцию с фреймворком [Echo](https://github.com/labstack/echo). 
Для этого используйте рендерер шаблонов, предоставляемый `goform`.

Пример интеграции:
```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/DBenyukh/goform/core"
)

type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
}

func main() {
	e := echo.New()

	// Инициализация рендерера шаблонов
	templateDir := "templates"
	renderer, err := core.NewTemplateRenderer(templateDir, "default.html")
	if err != nil {
		e.Logger.Fatal("Failed to create template renderer:", err)
	}
	e.Renderer = renderer

	e.GET("/register", func(c echo.Context) error {
		model := &RegistrationForm{
			Method: "POST",
			FormID: "register_form",
		}
		form := core.NewForm(model, model.Method, model.FormID)
		form.RenderHTML = true

		// Рендеринг формы
		return c.Render(http.StatusOK, "default.html", form.ToResponse())
	})

	e.Logger.Fatal(e.Start(":8080"))
}
```

---

### Кастомные сообщения об ошибках
Вы можете указывать кастомные сообщения об ошибках валидации с помощью тега `validate_msg` в структуре формы. Например:

```go
type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3" validate_msg:"Username must be at least 3 characters"`
	Email    string `form:"email" validate:"required,email" validate_msg:"Please provide a valid email address"`
	Password string `form:"password" validate:"required" validate_msg:"Password is required"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
}
```
Эти сообщения будут использоваться при валидации и отображаться в форме, если данные не соответствуют правилам.

---

### Скрытые поля
Вы можете создавать скрытые поля, которые не отображаются в форме, но передаются на сервер. Для этого установите поле `Hidden` в `true`:
```go
type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
	HiddenField string `form:"hidden_field" hidden:"true"` // Скрытое поле
}
```

В HTML-шаблоне такие поля не будут отображаться, но их значения будут переданы на сервер.

---

### CSRF-токены
Пакет поддерживает генерацию и проверку CSRF-токенов для защиты от атак. Для этого используйте метод `AddCSRFToken`:
```go
token, err := core.GenerateCSRFToken()
if err != nil {
    http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
}
form.AddCSRFToken(token)
```

CSRF-токен автоматически добавляется в форму и проверяется при обработке POST-запросов.

Пример использования в обработчике:
```go
if r.Method == http.MethodPost {
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
```

---

## Пример HTML-шаблона
Пример шаблона `default.html` для рендеринга формы:
```html
<form method="{{ if eq .Method "GET" }}GET{{ else }}POST{{ end }}">
    {{ if and (ne .Method "GET") (ne .Method "POST")  }}
        <input type="hidden" name="_method" value="{{ .Method }}">
    {{ end }}
    <input type="hidden" name="form_id" value="{{ .FormID }}">
    {{ range .Fields }}
        {{ if not .Hidden }}
        <div>
            <label>{{ .Name }}</label>
            <input type="{{ .Type }}" name="{{ $.FormID }}_{{ .Name }}" value="{{ .Value }}">
            {{ if .Error }}
                <span style="color: red;">{{ .Error }}</span>
            {{ end }}
        </div>
        {{ end }}
    {{ end }}
    <input type="hidden" name="{{ .FormID }}_csrf_token" value="{{ .CSRF }}">
    <button type="submit">Submit</button>
</form>
```

---

## Примеры
### Пример 1: Простая форма регистрации
См. раздел [Быстрый старт](#быстрый-старт).

---

### Пример 2: Кастомная валидация
```go
form.AddCustomValidation("username", func(value string) error {
    if strings.Contains(value, " ") {
        return errors.New("username cannot contain spaces")
    }
    return nil
})
```

---

### Пример 3: Разделённый фронтенд и бэкенд с `net/http`
Бэкенд (Go)

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/DBenyukh/goform/core"
)

type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
}

func main() {
	// В HandleFunc первым аргументом указываем роут до вашего api регистрации
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		model := &RegistrationForm{
			Method: "POST",
			FormID: "register_form",
		}
		form := core.NewForm(model, model.Method, model.FormID)
		form.RenderHTML = false // Возвращаем JSON

		if r.Method == http.MethodGet {
			// Возвращаем данные формы в формате JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(form.ToResponse())
			return
		}

		if r.Method == http.MethodPost {
			// Привязка данных из запроса к форме
			if err := form.Bind(r); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			// Валидация данных
			if err := form.Validate(model); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(form.ToResponse())
				return
			}

			// Обработка успешной отправки формы
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Фронтенд (JavaScript + HTML)
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Registration Form</title>
</head>
<body>
    <div id="form-container"></div>
    <script>
        // Загрузка формы с бэкенда
        let apiRegister = 'http://localhost:8080/api/register' // указываем ссылку на ваш api
        fetch(apiRegister)
            .then(response => response.json())
            .then(data => {
                const formContainer = document.getElementById('form-container');
                let formHTML = `<form method="POST" action=`apiRegister`>`; // используем ссылку на api в action

                data.fields.forEach(field => {
                    formHTML += `
                        <div>
                            <label>${field.name}</label>
                            <input type="${field.type}" name="${field.name}" value="${field.value}">
                            ${field.error ? `<span style="color: red;">${field.error}</span>` : ''}
                        </div>`;
                });

                formHTML += `
                    <input type="hidden" name="csrf_token" value="${data.csrf}">
                    <button type="submit">Submit</button>
                </form>`;

                formContainer.innerHTML = formHTML;
            });

        // Обработка отправки формы
        document.addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);

            const response = await fetch(apiRegister, {
                method: 'POST',
                body: formData,
            });

            if (response.ok) {
                const result = await response.json();
                alert(result.message);
            } else {
                const errors = await response.json();
                alert(`Validation errors: ${JSON.stringify(errors)}`);
            }
        });
    </script>
</body>
</html>
```
---

### Пример 4: Разделённый фронтенд и бэкенд с `Echo`
Бэкенд (Go + Echo)

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/DBenyukh/goform/core"
)

type RegistrationForm struct {
	Username string `form:"username" validate:"required,min=3" validate_msg:"Username must be at least 3 characters"`
	Email    string `form:"email" validate:"required,email" validate_msg:"Please provide a valid email address"`
	Password string `form:"password" validate:"required" validate_msg:"Password is required"`
	Method   string `form:"-"`
	FormID   string `form:"-"`
}

func main() {
	e := echo.New()
	
	// В e.GET первым аргументом указываем роут до вашего api регистрации
	e.GET("/api/register", func(c echo.Context) error {
		model := &RegistrationForm{
			Method: "POST",
			FormID: "register_form",
		}
		form := core.NewForm(model, model.Method, model.FormID)
		form.RenderHTML = false // Возвращаем JSON

		return c.JSON(http.StatusOK, form.ToResponse())
	})
	
	// В e.POST первым аргументом указываем роут до вашего api регистрации
	e.POST("/api/register", func(c echo.Context) error {
		model := &RegistrationForm{
			Method: "POST",
			FormID: "register_form",
		}
		form := core.NewForm(model, model.Method, model.FormID)
		form.RenderHTML = false

		// Привязка данных из запроса к форме
		if err := form.Bind(c.Request()); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
		}

		// Валидация данных
		if err := form.Validate(model); err != nil {
			return c.JSON(http.StatusBadRequest, form.ToResponse())
		}

		// Обработка успешной отправки формы
		return c.JSON(http.StatusOK, map[string]string{"message": "User registered successfully!"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
```

Фронтенд (JavaScript + HTML)

Фронтенд остаётся таким же, как в предыдущем примере, так как API бэкенда не изменилось.

---

## Лицензия
Этот проект распространяется под лицензией MIT. Подробнее см. в файле [LICENSE](https://github.com/DBenyukh/goform/tree/master/LICENSE).
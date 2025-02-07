package core

import (
	"net/http"
	"reflect"
)

// Form представляет HTML-форму.
type Form struct {
	Fields []*Field          // Поля формы
	CSRF   string            // CSRF-токен
	errors map[string]string // Ошибки валидации
}

// RenderData представляет данные для рендеринга формы.
type RenderData struct {
	Fields []*Field
	Errors map[string]string
	CSRF   string
}

// NewForm создает новую форму на основе модели.
func NewForm(model interface{}) *Form {
	fields := parseModel(reflect.ValueOf(model))
	return &Form{
		Fields: fields,
		errors: make(map[string]string),
	}
}

// AddField добавляет поле в форму.
func (f *Form) AddField(field *Field) {
	f.Fields = append(f.Fields, field)
}

// Bind привязывает данные из запроса к форме.
func (f *Form) Bind(r *http.Request) error {
	return bindForm(r, f)
}

// Render возвращает HTML-код формы.
func (f *Form) Render() RenderData {
	return RenderData{
		Fields: f.Fields,
		Errors: f.errors,
		CSRF:   f.CSRF,
	}
}

// Validate проверяет данные формы.
func (f *Form) Validate(model interface{}) error {
	return validateForm(f, model)
}

// AddCSRFToken добавляет CSRF-токен в форму.
func (f *Form) AddCSRFToken(token string) {
	f.CSRF = token
}

// Errors возвращает список ошибок формы.
func (f *Form) Errors() []string {
	var errors []string
	for _, err := range f.errors {
		errors = append(errors, err)
	}
	return errors
}

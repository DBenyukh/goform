package core

import (
	"net/http"
	"reflect"
	"strings"
)

// Form представляет HTML-форму.
type Form struct {
	Fields []*Field          // Поля формы
	CSRF   string            // CSRF-токен
	errors map[string]string // Ошибки валидации
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

// Render возвращает HTML-код формы.
func (f *Form) Render() string {
	var sb strings.Builder

	sb.WriteString("<form method='POST'>")

	for _, field := range f.Fields {
		sb.WriteString("<div>")
		sb.WriteString("<label>" + field.Name + "</label>")

		value := ""
		if field.Value != nil {
			value = field.Value.(string)
		}

		sb.WriteString("<input type='text' name='" + field.Name + "' value='" + value + "'>")

		if field.Error != "" {
			sb.WriteString("<span style='color: red;'>" + field.Error + "</span>")
		}

		sb.WriteString("</div>")
	}

	sb.WriteString("<button type='submit'>Submit</button>")
	sb.WriteString("</form>")

	return sb.String()
}

// Bind привязывает данные из запроса к форме.
func (f *Form) Bind(r *http.Request) error {
	return bindForm(r, f)
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

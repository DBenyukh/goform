package core

import (
	"net/http"
	"reflect"
)

// Form представляет HTML-форму.
type Form struct {
	Fields     []*Field          // Поля формы
	CSRF       string            // CSRF-токен
	Errs       map[string]string // Ошибки валидации
	Method     string            // Метод HTTP (GET, POST и т.д.)
	FormID     string            // Идентификатор формы
	RenderHTML bool              // Флаг для рендеринга HTML
}

// FormResponse представляет данные формы для ответа.
type FormResponse struct {
	Fields []FieldResponse // Упрощенная версия полей формы
	Errs   map[string]string
	CSRF   string
	Method string
	FormID string
}

// FieldResponse представляет упрощенную версию Field для ответа.
type FieldResponse struct {
	Name   string
	Type   string
	Value  string
	Error  string
	Hidden bool
}

// NewForm создает новую форму на основе модели.
func NewForm(model interface{}, method, formID string) *Form {
	fields := parseModel(reflect.ValueOf(model))
	return &Form{
		Fields: fields,
		Errs:   make(map[string]string),
		Method: method,
		FormID: formID,
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

// Validate проверяет данные формы.
func (f *Form) Validate(model interface{}) error {
	return validateForm(f, model)
}

// AddCustomValidation метод для добавления кастомных правил валидации
func (f *Form) AddCustomValidation(fieldName string, fn ValidationFunc) {
	for _, field := range f.Fields {
		if field.Name == fieldName {
			field.CustomValidation = fn
			break
		}
	}
}

// AddCSRFToken добавляет CSRF-токен в форму.
func (f *Form) AddCSRFToken(token string) {
	f.CSRF = token
}

// Errors возвращает список ошибок формы.
func (f *Form) Errors() []string {
	var errors []string
	for _, err := range f.Errs {
		errors = append(errors, err)
	}
	return errors
}

// UpdateModelFromForm обновляет поля модели на основе данных из формы.
func UpdateModelFromForm(model interface{}, form *Form) error {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("form")

		// Если тег form:"-", пропускаем это поле
		if tag == "-" {
			continue
		}

		// Находим соответствующее поле в форме
		for _, formField := range form.Fields {
			if formField.Name == tag {
				// Обновляем поле модели значением из формы
				fieldValue := val.Field(i)
				if fieldValue.CanSet() {
					fieldValue.SetString(formField.Value.(string))
				}
				break
			}
		}
	}

	return nil
}

// ToResponse возвращает данные формы в зависимости от флага RenderHTML.
func (f *Form) ToResponse() interface{} {
	if f.RenderHTML {
		return f.ToHTMLResponse()
	}
	return f.ToJSONResponse()
}

// ToHTMLResponse возвращает данные для рендеринга HTML.
func (f *Form) ToHTMLResponse() FormResponse {
	fields := make([]FieldResponse, len(f.Fields))
	for i, field := range f.Fields {
		value := ""
		if field.Value != nil {
			value = field.Value.(string)
		}
		fields[i] = FieldResponse{
			Name:   field.Name,
			Type:   field.Type,
			Value:  value,
			Error:  field.Error,
			Hidden: field.Hidden,
		}
	}

	return FormResponse{
		Fields: fields,
		Errs:   f.Errs,
		CSRF:   f.CSRF,
		Method: f.Method,
		FormID: f.FormID,
	}
}

// ToJSONResponse возвращает данные формы в формате JSON.
func (f *Form) ToJSONResponse() map[string]interface{} {
	data := make(map[string]interface{})
	for _, field := range f.Fields {
		data[field.Name] = map[string]interface{}{
			"type":  field.Type,
			"value": field.Value,
			"error": field.Error,
		}
	}
	return data
}

// AddError добавляет ошибку для указанного поля.
func (f *Form) AddError(fieldName, errorMessage string) {
	f.Errs[fieldName] = errorMessage
}

// GetErrors возвращает мапу ошибок.
func (f *Form) GetErrors() map[string]string {
	return f.Errs
}

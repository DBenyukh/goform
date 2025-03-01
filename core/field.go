package core

// Field представляет поля формы.
type Field struct {
	Name             string         // Имя поля
	Type             string         // Тип поля (text, email, password и т.д.)
	Value            interface{}    // Значение поля
	Error            string         // Ошибка валидации
	Hidden           bool           // Скрытое поле
	CustomValidation ValidationFunc // Кастомная функция валидации
}

// NewField создает новое поле.
func NewField(name, fieldType string) *Field {
	return &Field{
		Name: name,
		Type: fieldType,
	}
}

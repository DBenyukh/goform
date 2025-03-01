package core

import (
	"reflect"
)

// parseModel парсит структуру и создает поля формы.
func parseModel(val reflect.Value) []*Field {
	var fields []*Field

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("form")

		// Если тег form:"-", пропускаем это поле
		if tag == "-" {
			continue
		}

		// Если тег form пустой, поле считается скрытым
		hidden := tag == ""

		// Создаем поле формы
		formField := NewField(tag, getFieldType(field.Type.Kind()))
		formField.Hidden = hidden // Устанавливаем, является ли поле скрытым
		formField.Value = ""      // Инициализируем значение пустой строкой

		// Добавляем поле в список полей формы
		fields = append(fields, formField)
	}

	return fields
}

// getFieldType возвращает тип поля формы на основе типа Go.
func getFieldType(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "text"
	case reflect.Int, reflect.Int64:
		return "number"
	case reflect.Bool:
		return "checkbox"
	case reflect.Float32, reflect.Float64:
		return "number"
	default:
		return "text"
	}
}

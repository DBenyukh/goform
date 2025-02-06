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
		if tag == "" {
			continue
		}

		formField := NewField(tag, getFieldType(field.Type.Kind()))
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

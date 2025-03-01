package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationFunc func(value string) error

// validateForm проверяет данные формы.
func validateForm(form *Form, model interface{}, fieldsToValidate ...string) error {
	val := reflect.ValueOf(model).Elem()
	typeOfModel := val.Type()

	for _, field := range form.Fields {
		// Пропуск полей, которые не нужно валидировать
		if len(fieldsToValidate) > 0 && !contains(fieldsToValidate, field.Name) {
			continue
		}

		value := field.Value.(string)

		// Вызов кастомной функции валидации
		if field.CustomValidation != nil {
			if err := field.CustomValidation(value); err != nil {
				field.Error = err.Error()
				form.Errs[field.Name] = field.Error
				continue // Пропустить стандартную валидацию, если кастомная валидация не прошла
			}
		}

		// Стандартная валидация
		rules := getValidationRules(model, field.Name)
		customMsg := getCustomErrorMessage(typeOfModel, field.Name)

		for _, rule := range rules {
			switch {
			case rule == "required" && value == "":
				field.Error = customMsg
				form.Errs[field.Name] = field.Error
			case strings.HasPrefix(rule, "min="):
				min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
				if len(value) < min {
					if containsFormatSpecifier(customMsg) {
						field.Error = fmt.Sprintf(customMsg, min)
					} else {
						field.Error = customMsg
					}
					form.Errs[field.Name] = field.Error
				}
			case strings.HasPrefix(rule, "max="):
				max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
				if len(value) > max {
					if containsFormatSpecifier(customMsg) {
						field.Error = fmt.Sprintf(customMsg, max)
					} else {
						field.Error = customMsg
					}
					form.Errs[field.Name] = field.Error
				}
			case rule == "email" && !strings.Contains(value, "@"):
				field.Error = customMsg
				form.Errs[field.Name] = field.Error
			}
		}
	}

	if len(form.Errs) > 0 {
		return fmt.Errorf("validation errors")
	}
	return nil
}

// getValidationRules возвращает правила валидации для поля.
func getValidationRules(model interface{}, fieldName string) []string {
	val := reflect.ValueOf(model).Elem()
	typeOfModel := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typeOfModel.Field(i)
		if field.Tag.Get("form") == fieldName {
			validateTag := field.Tag.Get("validate")
			if validateTag != "" {
				return strings.Split(validateTag, ",")
			}
		}
	}

	return nil
}

// getCustomErrorMessage извлекает кастомное сообщение об ошибке из тега validate_msg.
func getCustomErrorMessage(modelType reflect.Type, fieldName string) string {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Tag.Get("form") == fieldName {
			return field.Tag.Get("validate_msg")
		}
	}
	return ""
}

// containsFormatSpecifier проверяет, содержит ли строка форматирующие спецификаторы.
func containsFormatSpecifier(s string) bool {
	return strings.Contains(s, "%")
}

// contains проверяет, содержится ли строка в слайсе.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

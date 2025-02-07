package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// validateForm проверяет данные формы.
func validateForm(form *Form, model interface{}) error {
	val := reflect.ValueOf(model).Elem()
	typeOfModel := val.Type()

	for _, field := range form.Fields {
		value := field.Value.(string)
		rules := getValidationRules(model, field.Name)
		customMsg := getCustomErrorMessage(typeOfModel, field.Name)

		for _, rule := range rules {
			switch {
			case rule == "required" && value == "":
				field.Error = customMsg
				form.errors[field.Name] = field.Error
			case strings.HasPrefix(rule, "min="):
				min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
				if len(value) < min {
					if containsFormatSpecifier(customMsg) {
						field.Error = fmt.Sprintf(customMsg, min)
					} else {
						field.Error = customMsg
					}
					form.errors[field.Name] = field.Error
				}
			case strings.HasPrefix(rule, "max="):
				max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
				if len(value) > max {
					if containsFormatSpecifier(customMsg) {
						field.Error = fmt.Sprintf(customMsg, max)
					} else {
						field.Error = customMsg
					}
					form.errors[field.Name] = field.Error
				}
			case rule == "email" && !strings.Contains(value, "@"):
				field.Error = customMsg
				form.errors[field.Name] = field.Error
			}
		}
	}
	if len(form.errors) > 0 {
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

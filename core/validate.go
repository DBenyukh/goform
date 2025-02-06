package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// validateForm проверяет данные формы.
func validateForm(form *Form, model interface{}) error {
	for _, field := range form.Fields {
		value := field.Value.(string)
		rules := getValidationRules(model, field.Name)

		for _, rule := range rules {
			switch {
			case rule == "required" && value == "":
				field.Error = "This field is required"
				form.errors[field.Name] = field.Error
			case strings.HasPrefix(rule, "min="):
				min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
				if len(value) < min {
					field.Error = "Value is too short"
					form.errors[field.Name] = field.Error
				}
			case strings.HasPrefix(rule, "max="):
				max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
				if len(value) > max {
					field.Error = "Value is too long"
					form.errors[field.Name] = field.Error
				}
			case rule == "email" && !strings.Contains(value, "@"):
				field.Error = "Invalid email address"
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

package core

import (
	"net/http"
)

// bindForm привязывает данные из запроса к форме.
func bindForm(r *http.Request, form *Form) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	for _, field := range form.Fields {
		// Учитываем FormID при извлечении значений
		key := form.FormID + "_" + field.Name
		value := r.FormValue(key)
		field.Value = value // Устанавливаем значение поля
	}

	return nil
}

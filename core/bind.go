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

	// Проходим по полям формы и заполняем значениями из запроса
	for _, field := range form.Fields {
		value := r.FormValue(field.Name)
		field.Value = value
	}

	return nil
}

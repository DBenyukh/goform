package core

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// renderTemplate рендерит форму с использованием шаблона.
func renderTemplate(w http.ResponseWriter, templateName string, form *Form) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles(filepath.Join("templates", templateName))
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Выполняем рендеринг шаблона
	err = tmpl.Execute(w, form)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
	}
}

package core

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"path/filepath"
)

// TemplateRenderer содержит шаблоны для рендеринга
type TemplateRenderer struct {
	Templates *template.Template
}

// Render выполняет рендеринг шаблона
func (tr *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := tr.Templates.ExecuteTemplate(w, name, data) // Используем имя шаблона без пути
	if err != nil {
		log.Println("Error rendering template:", err)
		return err
	}
	return nil
}

// NewTemplateRenderer инициализирует и возвращает новый рендерер шаблонов
func NewTemplateRenderer() *TemplateRenderer {
	// Получаем абсолютный путь к корню проекта
	projectDir, err := filepath.Abs("..") // Переходим на один уровень выше, чтобы получить корень проекта
	if err != nil {
		log.Fatalf("Error getting absolute project directory path: %v", err)
	}

	// Путь к папке с шаблонами (в корне проекта)
	templateDir := filepath.Join(projectDir, "templates")

	// Используем ParseGlob для поиска шаблонов в указанной папке
	tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	// Возвращаем рендерер с загруженными шаблонами
	return &TemplateRenderer{
		Templates: tmpl,
	}
}

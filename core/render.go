package core

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"path/filepath"
)

// TemplateRenderer содержит шаблоны для рендеринга
type TemplateRenderer struct {
	Templates       *template.Template
	DefaultTemplate string // Имя шаблона по умолчанию
}

// Render выполняет рендеринг шаблона
func (tr *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Используем переданное имя шаблона или имя по умолчанию
	templateName := name
	if templateName == "" {
		templateName = tr.DefaultTemplate
	}

	// Выполняем рендеринг шаблона
	err := tr.Templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Println("Error rendering template:", err)
		return err
	}
	return nil
}

// NewTemplateRenderer инициализирует и возвращает новый рендерер шаблонов
func NewTemplateRenderer(templateDir string, defaultTemplate string) (*TemplateRenderer, error) {
	// Используем ParseGlob для поиска шаблонов в указанной папке
	tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error loading templates: %v", err)
	}

	// Возвращаем рендерер с загруженными шаблонами
	return &TemplateRenderer{
		Templates:       tmpl,
		DefaultTemplate: defaultTemplate,
	}, nil
}

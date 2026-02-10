package template_manager

import (
	"bytes"
	"fmt"
	"html/template"
	"marketplace/internal/notification-service/domain"
	"path/filepath"
)

type templateManager struct {
	basePath string
}

func NewTemplateManager(basePath string) domain.TemplateManager {
	return &templateManager{basePath: basePath}
}

func (t *templateManager) Render(templateName string, data interface{}) (string, error) {

	tmplPath := filepath.Join("internal", "notification-service", "templates", templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("template %s not found: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

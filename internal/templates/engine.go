package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"phpier/internal/config"
)

//go:embed files
var templateFS embed.FS

// Engine represents the template engine
type Engine struct {
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

// TemplateData represents data passed to templates
type TemplateData struct {
	Project *config.ProjectConfig
	Global  *config.GlobalConfig
}

// NewEngine creates a new template engine
func NewEngine() *Engine {
	engine := &Engine{
		templates: make(map[string]*template.Template),
		funcMap:   createFuncMap(),
	}

	if err := engine.loadTemplates(); err != nil {
		panic(fmt.Sprintf("Failed to load templates: %v", err))
	}

	return engine
}

// Render renders a template with the given data
func (e *Engine) Render(templateName string, data *TemplateData) (string, error) {
	tmpl, exists := e.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// RenderProjectDockerCompose renders the docker-compose.yml for a project
func (e *Engine) RenderProjectDockerCompose(projectCfg *config.ProjectConfig, globalCfg *config.GlobalConfig) (string, error) {
	data := &TemplateData{
		Project: projectCfg,
		Global:  globalCfg,
	}
	return e.Render("docker-compose/project.yml", data)
}

// RenderGlobalDockerCompose renders the docker-compose.yml for the global services
func (e *Engine) RenderGlobalDockerCompose(globalCfg *config.GlobalConfig) (string, error) {
	data := &TemplateData{
		Global: globalCfg,
	}
	return e.Render("docker-compose/global.yml", data)
}

// RenderPHPDockerfile renders the appropriate PHP Dockerfile based on version
func (e *Engine) RenderPHPDockerfile(projectCfg *config.ProjectConfig) (string, error) {
	templateName := e.selectPHPDockerfileTemplate(projectCfg.PHP)
	data := &TemplateData{
		Project: projectCfg,
	}
	return e.Render(templateName, data)
}

// RenderTraefikConfig renders the traefik.yml configuration
func (e *Engine) RenderTraefikConfig(globalCfg *config.GlobalConfig) (string, error) {
	data := &TemplateData{
		Global: globalCfg,
	}
	return e.Render("configs/traefik.yml", data)
}

// RenderTraefikDynamicConfig renders the traefik dynamic configuration
func (e *Engine) RenderTraefikDynamicConfig(globalCfg *config.GlobalConfig) (string, error) {
	data := &TemplateData{
		Global: globalCfg,
	}
	return e.Render("configs/traefik-dynamic.yml", data)
}

// RenderPHPConfig renders the php.ini configuration
func (e *Engine) RenderPHPConfig() (string, error) {
	data := &TemplateData{}
	return e.Render("configs/php.ini", data)
}

// RenderNginxConfig renders the main nginx.conf configuration
func (e *Engine) RenderNginxConfig(projectCfg *config.ProjectConfig) (string, error) {
	data := &TemplateData{
		Project: projectCfg,
	}
	return e.Render("configs/nginx.conf", data)
}

// RenderNginxSiteConfig renders the nginx site configuration (default.conf)
func (e *Engine) RenderNginxSiteConfig(projectCfg *config.ProjectConfig, globalCfg *config.GlobalConfig) (string, error) {
	data := &TemplateData{
		Project: projectCfg,
		Global:  globalCfg,
	}
	return e.Render("configs/nginx-site.conf", data)
}

func (e *Engine) selectPHPDockerfileTemplate(phpVersion string) string {
	switch phpVersion {
	case "5.6", "7.0", "7.1", "7.2", "7.3":
		return "dockerfiles/php56-73.Dockerfile"
	case "7.4", "8.0":
		return "dockerfiles/php74-80.Dockerfile"
	case "8.1", "8.2", "8.3", "8.4":
		return "dockerfiles/php81-84.Dockerfile"
	default:
		return "dockerfiles/php81-84.Dockerfile"
	}
}

// loadTemplates loads all template files from the embedded filesystem
func (e *Engine) loadTemplates() error {
	return fs.WalkDir(templateFS, "files", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".tpl" {
			return nil
		}

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		templateName := strings.TrimPrefix(path, "files/")
		templateName = strings.TrimSuffix(templateName, ".tpl")

		tmpl, err := template.New(filepath.Base(path)).Funcs(e.funcMap).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", templateName, err)
		}

		e.templates[templateName] = tmpl
		return nil
	})
}

// createFuncMap creates template functions
func createFuncMap() template.FuncMap {
	return template.FuncMap{
		"serviceEnabled": func(service string, globalCfg *config.GlobalConfig) bool {
			if globalCfg == nil {
				return false
			}
			switch service {
			case "redis":
				return globalCfg.Services.Cache.Redis.Enabled
			case "memcached":
				return globalCfg.Services.Cache.Memcached.Enabled
			case "phpmyadmin":
				return globalCfg.Services.Tools.PHPMyAdmin
			case "mailpit":
				return globalCfg.Services.Tools.Mailpit.Enabled
			case "pgadmin":
				return globalCfg.Services.Tools.PgAdmin
			}
			return false
		},
		"default": func(defaultValue interface{}, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"getHostRule": func(projectCfg *config.ProjectConfig, globalCfg *config.GlobalConfig) string {
			// Auto-generate domain from project name + global domain
			domain := projectCfg.Name + "." + globalCfg.Traefik.Domain
			return "Host(`" + domain + "`)"
		},
	}
}

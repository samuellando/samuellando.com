package template

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	templates *template.Template
}

type FuncMap template.FuncMap
type HTML template.HTML

func New(name string) *Template {
	templates := template.New("templates")
	return &Template{templates: templates}
}

func (temps *Template) Funcs(funcs FuncMap) *Template {
    temps.templates = temps.templates.Funcs(template.FuncMap(funcs))
    return temps
}

func (temps *Template) Lookup(name string) *template.Template {
	return temps.templates.Lookup(name)
}

func (temps *Template) ExecuteTemplate(wr io.Writer, name string, data any) error {
    if html, ok := data.(HTML); ok {
        data = template.HTML(html)
    }
	return temps.templates.ExecuteTemplate(wr, name, data)
}

func (temps *Template) ParseTemplates(templateDir string) *Template {
	layouts := make(map[string]string)
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// Setup the layout for this directory, copy the parent directory if it exists.
			if layout, ok := layouts[filepath.Dir(path)]; ok {
				layouts[path] = layout
			} else {
				layouts[path] = "{{slot}}"
			}
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if info.Name() == "+layout.html" {
			// Nest this layout into the layout from the parent directory.
			layout, ok := layouts[filepath.Dir(filepath.Dir(path))]
			if !ok {
				panic("Layout nesting error")
			}
			layout = strings.ReplaceAll(layout, "{{slot}}", string(data))
			layouts[filepath.Dir(path)] = layout
			return nil
		}
		if strings.HasSuffix(path, ".html") {
			var templateName string
			var template string
			if info.Name() == "+page.html" {
				var layout string
				var ok bool
				if layout, ok = layouts[filepath.Dir(path)]; !ok {
					layout = "{{slot}}"
				}
				templateName = filepath.Dir(path)[len(templateDir)+1:]
				fmt.Println(templateName)
				template = fmt.Sprintf(`{{define "%s"}}%s{{end}}`, templateName, layout)
				template = strings.ReplaceAll(template, "{{slot}}", string(data))
			} else {
				templateName = strings.ReplaceAll(info.Name(), ".html", "")
				template = fmt.Sprintf(`{{define "%s"}}%s{{end}}`, templateName, string(data))
			}
			_, err = temps.templates.Parse(template)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return temps
}

// This package provides additional parsing functions to the stdlib template
// Library.
// It parses a template directory recursively in a svelte like like fashion, as follows:
//
// - Any template with the name +layout.html will be stored in memory for later.
//   - Layouts should have a {{slot}} component inside them.
//   - Layouts will be inherited from parent directories, and nested into parent layouts.
//
// - Any template with the name +page.html will be parsed, and nested into the layout, and parent directory layouts.
//   - These will be loaded with the name corresponding to the directory path.
//
// - Any other x.html file will be loaded as is with the name `x`
//   - Components are essentially globally scoped.
//
// Other other function definitions function exactly like the template package in stdlib
package template

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	templates *template.Template
}

type FuncMap template.FuncMap

// ADDED FUNCTIONALITY

// Parse all templates in the specified directory, including
// layouts and components, and returns the Template instance.
func (temps *Template) ParseTemplates(templateDir string) *Template {
	fileSys := os.DirFS(templateDir)
	return temps.ParseFs(fileSys)
}

// Parse all templates from the provided file system, including layouts
// and components, and returns the Template instance.
func (temps *Template) ParseFs(files fs.FS) *Template {
	layouts := make(map[string]string)
	err := fs.WalkDir(files, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			initDirLayout(layouts, path)
			return nil
		}
		return temps.handleFile(files, layouts, path, info)
	})
	if err != nil {
		panic(err)
	}
	return temps
}

// WRAPPER FUNCTIONS

func New(name string) *Template {
	templates := template.New("templates")
	return &Template{templates: templates}
}

func (temps *Template) Funcs(funcs FuncMap) *Template {
	temps.templates = temps.templates.Funcs(template.FuncMap(funcs))
	return temps
}

func (temps *Template) Lookup(name string) *Template {
	t := temps.templates.Lookup(name)
	if t == nil {
		return nil
	}
	return &Template{templates: t}
}

func (temps *Template) Execute(wr io.Writer, data any) error {
	return temps.templates.Execute(wr, data)
}

func (temps *Template) ExecuteTemplate(wr io.Writer, name string, data any) error {
	return temps.templates.ExecuteTemplate(wr, name, data)
}

// HELPERS

func initDirLayout(layouts map[string]string, path string) {
	if parent, ok := layouts[filepath.Dir(path)]; ok {
		layouts[path] = parent
	} else {
		layouts[path] = "{{slot}}"
	}
}

func (temps *Template) handleFile(files fs.FS, layouts map[string]string, path string, info fs.DirEntry) error {
	data, err := fs.ReadFile(files, path)
	if err != nil {
		return err
	}
	if info.Name() == "+layout.html" {
		temps.loadLayoutFile(layouts, path, data)
		return nil
	}
	if info.Name() == "+page.html" {
		return temps.loadPageFile(layouts, path, data)
	}
	if strings.HasSuffix(path, ".html") {
		return temps.loadComponentFile(path, data)
	}
	return nil
}

func (temps *Template) loadLayoutFile(layouts map[string]string, path string, data []byte) {
	if parentLayout, ok := layouts[filepath.Dir(path)]; ok {
		layout := strings.ReplaceAll(parentLayout, "{{slot}}", string(data))
		layouts[filepath.Dir(path)] = layout
	} else {
		layouts[filepath.Dir(path)] = string(data)
	}
}

func (temps *Template) loadPageFile(layouts map[string]string, path string, data []byte) error {
	template := string(data)
	if layout, ok := layouts[filepath.Dir(path)]; ok {
		template = strings.ReplaceAll(layout, "{{slot}}", template)
	}
	return temps.loadTemplate(filepath.Dir(path), template)
}

func (temps *Template) loadComponentFile(path string, data []byte) error {
	template := string(data)
	name := filepath.Base(path)
	name = name[:len(name)-5]
	return temps.loadTemplate(name, template)
}

func (temps *Template) loadTemplate(name, content string) error {
	if temps.templates.Lookup(name) != nil {
		return fmt.Errorf("Template collision")
	}
	template := fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, content)
	_, err := temps.templates.Parse(template)
	return err
}

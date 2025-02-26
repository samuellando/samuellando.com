package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

// Test that functions add added and are callable from inside a template
func TestFuncs(t *testing.T) {
	tempDir := t.TempDir()
	pageContent := `{{toUpper "hello world"}}`

	os.Mkdir(filepath.Join(tempDir, "test"), 0755)
	os.WriteFile(filepath.Join(tempDir, "test", "comp.html"), []byte(pageContent), 0644)
	tmpl := New("test")
	funcs := FuncMap{
		"toUpper": strings.ToUpper,
	}
	tmpl.Funcs(funcs)
	tmpl.ParseTemplates(tempDir)
	w := new(strings.Builder)
	tmpl.ExecuteTemplate(w, "comp", nil)
	if w.String() != "HELLO WORLD" {
		t.Fail()
	}
}

// Test that +page.html templates are correctly named and loaded
func TestParseTemplatesPage(t *testing.T) {
	// Setup a temporary directory with test templates
	tempDir := t.TempDir()
	pageContent := "<a>Hello, World!</a>"

	os.Mkdir(filepath.Join(tempDir, "dir"), 0755)
	os.Mkdir(filepath.Join(tempDir, "dir/test"), 0755)
	os.WriteFile(filepath.Join(tempDir, "dir/test", "+page.html"), []byte(pageContent), 0644)

	tmpl := New("test")
	tmpl.ParseTemplates(tempDir)

	if tmpl.Lookup("dir/test") == nil {
		t.Error("Expected template 'dir/test' to be parsed")
	}
	w := new(strings.Builder)
	tmpl.ExecuteTemplate(w, "dir/test", nil)
	if w.String() != "<a>Hello, World!</a>" {
		t.Fatal("Content should match")
	}
}

// Test that +layout templates correctly nest arround pages, and not components.
func TestParseTemplatesLayouts(t *testing.T) {
	// Setup a temporary directory with test templates
	tempDir := t.TempDir()
	layoutContent := "<html>{{slot}}</html>"
	layout2Content := "<h2>OK</h2><h1>{{slot}}</h1>"
	pageContent := "<a>Hello, World!</a>"

	os.Mkdir(filepath.Join(tempDir, "dir"), 0755)
	os.Mkdir(filepath.Join(tempDir, "dir/test"), 0755)
	os.Mkdir(filepath.Join(tempDir, "dir/test/a"), 0755)
	os.Mkdir(filepath.Join(tempDir, "dir/test/a/b"), 0755)
	os.WriteFile(filepath.Join(tempDir, "dir", "+layout.html"), []byte(layoutContent), 0644)
	os.WriteFile(filepath.Join(tempDir, "dir/test/a/b", "+layout.html"), []byte(layout2Content), 0644)
	os.WriteFile(filepath.Join(tempDir, "dir/test/a/b", "+page.html"), []byte(pageContent), 0644)
	os.WriteFile(filepath.Join(tempDir, "dir/test/a/b", "component.html"), []byte(pageContent), 0644)

	tmpl := New("test")
	tmpl.ParseTemplates(tempDir)

	if tmpl.Lookup("dir/test/a/b") == nil {
		t.Error("Expected template 'dir/test' to be parsed")
	}
	w := new(strings.Builder)
    err := tmpl.ExecuteTemplate(w, "dir/test/a/b", nil)
    if err != nil {
        t.Fatal(err)
    }
	if w.String() != "<html><h2>OK</h2><h1><a>Hello, World!</a></h1></html>" {
		t.Fatal("Nesting should work", w.String())
	}
	w = new(strings.Builder)
    err = tmpl.ExecuteTemplate(w, "component", nil)
    if err != nil {
        t.Fatal(err)
    }
	if w.String() != "<a>Hello, World!</a>" {
		t.Fatal("Components should not be affected by layout")
	}
}

func TestHandleFile(t *testing.T) {
	temps := New("test")
	layouts := make(map[string]string)
	files := fstest.MapFS{
		"dir/+layout.html": &fstest.MapFile{Data: []byte("<html>{{slot}}</html>")},
		"dir/+page.html":   &fstest.MapFile{Data: []byte("<body>Content</body>")},
	}
	dir, err := files.ReadDir("dir")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range dir {
		if f.Name() == "+layout.html" {
			err = temps.handleFile(files, layouts, "dir/+layout.html", f)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if layouts["dir"] != "<html>{{slot}}</html>" {
				t.Error("Expected layout to be set for dir")
			}
		} else if f.Name() == "+page.html" {
			err = temps.handleFile(files, layouts, "dir/+page.html", f)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if temps.Lookup("dir") == nil {
				t.Error("Expected template 'dir' to be parsed")
			}
		} else {
			t.Fatal("weird", f.Name())
		}
	}
}

// If there is a parent layout, the layout should be nested in.
func TestLoadLayoutFileWithParent(t *testing.T) {
	temps := New("test")
	layouts := map[string]string{
		"parent/child": "<html>{{slot}}</html>",
	}
	temps.loadLayoutFile(layouts, "parent/child/+layout.html", []byte("<body>{{slot}}</body>"))
	expected := "<html><body>{{slot}}</body></html>"
	if layouts["parent/child"] != expected {
		t.Errorf("Expected layout to be %q, got %q", expected, layouts["parent/child"])
	}
}

// If there is no parent layout, the layout should be as is.
func TestLoadLayoutFileNoParent(t *testing.T) {
	temps := New("test")
	layouts := map[string]string{}
	temps.loadLayoutFile(layouts, "parent/child/+layout.html", []byte("<body>{{slot}}</body>"))
	expected := "<body>{{slot}}</body>"
	if layouts["parent/child"] != expected {
		t.Errorf("Expected layout to be %q, got %q", expected, layouts["parent/child"])
	}
}

// The page should use the provided layout
func TestLoadPageFileWithLayouts(t *testing.T) {
	temps := New("test")
	layouts := map[string]string{}
	temps.loadLayoutFile(layouts, "parent/child/+layout.html", []byte("<body>{{slot}}</body>"))
    temps.loadPageFile(layouts, "parent/child/+page.html", []byte("<h1>Hello</h1>"))
    expected := "<body><h1>Hello</h1></body>"
    w := new(strings.Builder)
    err := temps.ExecuteTemplate(w, "parent/child", nil)
    if err != nil {
       t.Fatal(err)
    }
	if w.String() != expected {
		t.Errorf("Expected layout to be %q, got %q", expected, w.String())
	}
}

// If there is no parent layout, the page should be as is.
func TestLoadPageFileWithoutLayout(t *testing.T) {
	temps := New("test")
	layouts := map[string]string{}
    temps.loadPageFile(layouts, "parent/child/+page.html", []byte("<h1>Hello</h1>"))
    expected := "<h1>Hello</h1>"
    w := new(strings.Builder)
    err := temps.ExecuteTemplate(w, "parent/child", nil)
    if err != nil {
       t.Fatal(err)
    }
	if w.String() != expected {
        t.Fatal("Unexpected result", w.String())
	}
}

// Layout should not affect components
func TestLoadComponentFileWithLayout(t *testing.T) {
	temps := New("test")
	layouts := map[string]string{}
	temps.loadLayoutFile(layouts, "p/c/+layout.html", []byte("<body>{{slot}}</body>"))
    temps.loadComponentFile("p/c/com.html", []byte("<h1>Hello</h1>"))
    expected := "<h1>Hello</h1>"
    w := new(strings.Builder)
    err := temps.ExecuteTemplate(w, "com", nil)
    if err != nil {
       t.Fatal(err)
    }
	if w.String() != expected {
        t.Fatal("Unexpected result", w.String())
	}
}

// Components should collide and throw error
func TestLoadComponentFileCollisions(t *testing.T) {
	temps := New("test")
    temps.loadComponentFile("p/c/com.html", []byte("<h1>Hello</h1>"))
    err := temps.loadComponentFile("p/c/com.html", []byte("<h1>Hello</h1>"))
    if err == nil {
       t.Fatal("Components should collide if they have the same name")
    }
}

func TestLoadTamplate(t *testing.T) {
	temps := New("test")
    temps.loadTemplate("a", "<h1>Hello</h1>")
    w := new(strings.Builder)
    expected := "<h1>Hello</h1>"
    err := temps.ExecuteTemplate(w, "a", nil)
    if err != nil {
       t.Fatal(err)
    }
	if w.String() != expected {
        t.Fatal("Unexpected result", w.String())
	}
}

func TestLoadTamplateCollision(t *testing.T) {
	temps := New("test")
    temps.loadTemplate("a", "<h1>Hello</h1>")
    err := temps.loadTemplate("a", "<h1>Hello</h1>")
    if  err == nil {
        t.Fatal("Collisons should throw errors")
    }
}

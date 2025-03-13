package template

import (
	"fmt"
	"log"
	"maps"
	"net/http"
	"path"
	"path/filepath"
)

type Context struct {
	values ContextValues
}

func newContext(values ...ContextValues) Context {
	ctx := Context{values: make(ContextValues)}
	for _, v := range values {
		ctx.addValues(v)
	}
	return ctx
}

func (c *Context) addValues(values ContextValues) {
	maps.Copy(c.values, values)
}

func (c Context) Get(k string) any {
	if v, ok := c.values[k]; ok {
		return v(c)
	}
	panic(k + " not found in context.")
}

type ContextValue func(Context) any
type ContextValues map[string]ContextValue

type Handler struct {
	Templates     Template
	ContextValues map[string]ContextValue
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctxt := newContext(h.ContextValues, ContextValues{
		"Req": func(ctx Context) any { return req },
		"Page": func(ctx Context) any {
			req := ctx.Get("Req").(*http.Request)
			return req.URL.Path
		},
	})
	h.renderTemplate(ctxt, w, req)
}

func (h *Handler) renderTemplate(ctxt Context, w http.ResponseWriter, req *http.Request) {
	template := path.Join("pages", ctxt.Get("Page").(string))
	// Check that the template exists
	if h.Templates.Lookup(template) == nil {
		// Check for a slug
		template = filepath.Dir(template) + "/[slug]"
		if h.Templates.Lookup(template) == nil {
			http.NotFound(w, req)
			return
		}
	}
	err := h.Templates.ExecuteTemplate(w, template, ctxt)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%s : %s", http.StatusText(500), err), 500)
	}
}

package main

import (
	"strconv"
	"strings"

	"samuellando.com/internal/store/document"
	"samuellando.com/internal/store/project"
)

type sortFunctionReference[T any] struct {
	Name string
	Func func(T, T) bool
}

type groupFunctionReference[T any] struct {
	Name string
	Func func(T) string
}

var PROJECT_SORT_FUNCTIONS = map[string]sortFunctionReference[*project.Project]{
	"byName": sortFunctionReference[*project.Project]{
		Name: "By Name",
		Func: func(p1, p2 *project.Project) bool {
			return strings.Compare(p1.Title(), p2.Title()) < 0
		},
	},
	"byCreated": sortFunctionReference[*project.Project]{
		Name: "By Date Created",
		Func: func(p1, p2 *project.Project) bool {
			return p1.Created().After(p2.Created())
		},
	},
	"byLastPush": sortFunctionReference[*project.Project]{
		Name: "By Last Time Pushed",
		Func: func(p1, p2 *project.Project) bool {
			return p1.Pushed().After(p2.Pushed())
		},
	},
}

var PROJECT_GROUP_FUNCTIONS = map[string]groupFunctionReference[*project.Project]{
	"byCreated": groupFunctionReference[*project.Project]{
		Name: "By Year Created",
		Func: func(p *project.Project) string {
			return strconv.Itoa(p.Created().Year())
		},
	},
	"byLastPush": groupFunctionReference[*project.Project]{
		Name: "By Year of Last Push",
		Func: func(p *project.Project) string {
			return strconv.Itoa(p.Pushed().Year())
		},
	},
}

var DOCUMENT_SORT_FUNCTIONS = map[string]sortFunctionReference[*document.Document]{
	"byCreated": sortFunctionReference[*document.Document]{
		Name: "By Date Created",
		Func: func(d1, d2 *document.Document) bool {
			return d1.Created().After(d2.Created())
		},
	},
}

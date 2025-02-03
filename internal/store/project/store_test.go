package project

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setup() (Store, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("testData/sample.json")
		if err != nil {
            panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	ps := CreateStore(func(o *Options) {
		o.Url = ts.URL
	})
	return ps, ts
}

func teardown(ts *httptest.Server) {
	ts.Close()
}

func TestGetByIdBase(t *testing.T) {
	ps, ts := setup()
	defer teardown(ts)
	doc, _ := ps.GetById(1296269)
	if doc.Title() != "Hello-World" {
		t.Fatal("Ttitle does not match")
	}
}

func TestGetAllBase(t *testing.T) {
	ps, ts := setup()
	defer teardown(ts)
	projects, _ := ps.GetAll()
	if len(projects) != 4 {
		t.Fatal("Should have 4 projects")
	}
}

func TestFilter(t *testing.T) {
	ps, ts := setup()
	defer teardown(ts)
	filtered := ps.Filter(func(d *Project) bool {
		return strings.HasSuffix(d.Title(), "two")
	})
	data, _ := filtered.GetAll()
	if len(data) != 2 {
		t.Errorf("%d should contain 2 elements", len(data))
	}
}

func TestGetIdFiltered(t *testing.T) {
	ps, ts := setup()
	defer teardown(ts)
	filtered := ps.Filter(func(d *Project) bool {
		return strings.HasSuffix(d.Title(), "two")
	})
	elem, err := filtered.GetById(1299)
	if err != nil {
		t.Error(err)
	}
	if elem.Title() != "Hello-World-two" {
		t.Error("Wrong element")
	}
	_, err = filtered.GetById(1296269)
	if err == nil {
		t.Error("Should not be included")
	}
}

func TestSort(t *testing.T) {
	ps, ts := setup()
	defer teardown(ts)
	sorted := ps.Sort(func(a, b *Project) bool {
		return a.Id() < b.Id()
	})
	data, _ := sorted.GetAll()
	if len(data) != 4 {
		t.Errorf("%d should contain 4 elements", len(data))
	}
	for i, c := range []string{"Bye-World-two", "Hello-World-two", "Hello-World", "Bye-World"} {
		if data[i].Title() != c {
			t.Error("Out of order")
		}
	}
}

func TestGroup(t *testing.T) {
	ps, ts := setup()
	defer teardown(ts)
	groups := ps.Group(func(d *Project) string {
		return string(d.Title()[0])
	})
	if groups.Len() != 2 {
		t.Error("Wrong number of groups")
	}
	expectedLens := map[string]int{"H": 2, "B": 2}
	for k, s := range groups.All() {
		data, _ := s.GetAll()
		if len(data) != expectedLens[k] {
			t.Errorf("%d should contain %d elements", len(data), expectedLens[k])
		}
	}
}

func TestStack(t *testing.T) {
	ps, ts := setup()
	defer teardown(ts)
	g, _ := ps.Sort(func(a, b *Project) bool {
		return a.Id() < b.Id()
	}).Group(func(p *Project) string {
        return strings.Split(p.Title(), "-")[0]
    }).Get("Bye")
    res, _ := g.Filter(func(p *Project) bool {
        return strings.HasSuffix(p.Title(), "two")
    }).GetAll()
    if res[0].Title() != "Bye-World-two" {
        t.Fatal("Wrong element")
    }
}

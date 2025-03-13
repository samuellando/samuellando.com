package project

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"samuellando.com/internal/db"
	"samuellando.com/internal/store/tag"
	"samuellando.com/internal/testutil"
)

func setup() (Store, *httptest.Server, *sql.DB) {
	con := db.ConnectPostgres(testutil.GetDbCredentials())
	if err := testutil.ResetDb(con, "projectTest"); err != nil {
		panic(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("testData/sample.json")
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	ps := CreateStore(con, func(o *Options) {
		o.Url = ts.URL
	})
	return ps, ts, con
}

func teardown(ts *httptest.Server, db *sql.DB) {
	ts.Close()
	db.Close()
}

func TestGetByIdBase(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	proj, _ := ps.GetById(1296269)
	if proj.Title() != "Hello-World" {
		t.Fatal("Ttitle does not match")
	}
	desc := "And your Last!"
	proj.Update(func(pp *ProtoProject) {
		pp.Description = &desc
		pp.Tags = []tag.ProtoTag{
			{Value: "one"},
			{Value: "two"},
		}
	})
	proj, _ = ps.GetById(1296269)
	if proj.Description() != desc {
		t.Fatal("description not loaded from internal")
	}
	if len(proj.Tags()) != 2 {
		t.Fatal("tags not loaded from internal")
	}
}

func TestGetAllBase(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	projects, err := ps.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 4 {
		t.Fatal("Should have 4 projects")
	}
	desc := "And your Last!"
	projects[0].Update(func(pp *ProtoProject) {
		pp.Description = &desc
		pp.Tags = []tag.ProtoTag{
			{Value: "one"},
			{Value: "two"},
		}
	})
	id := projects[0].Id()
	projects, _ = ps.GetAll()
	for _, proj := range projects {
		if proj.Id() == id {
			if proj.Description() != desc {
				t.Fatalf("description not loaded from internal '%s'", proj.Description())
			}
			if len(proj.Tags()) != 2 {
				t.Fatal("tags not loaded from internal")
			}
		}
	}
}

func TestGetDesc(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	id := int64(1296269)
	proj, _ := ps.GetById(id)
	if proj.Title() != "Hello-World" {
		t.Fatal("Wrong Title")
	}
	if proj.Description() != "This your first repo!" {
		t.Fatalf("Wrong desc %s", proj.Description())
	}
	desc := "And your Last!"
	proj.Update(func(pp *ProtoProject) {
		pp.Description = &desc
	})
	proj, _ = ps.GetById(id)
	if proj.Description() != "And your Last!" {
		t.Fatalf("Wrong desc %s", proj.Description())
	}
}

func TestGetTags(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	id := int64(1296269)
	doc, _ := ps.GetById(id)
	if len(doc.Tags()) != 0 {
		t.Fatal("There should be no tags")
	}
	doc.Update(func(pp *ProtoProject) {
		pp.Tags = []tag.ProtoTag{
			{Value: "one"},
			{Value: "two"},
		}
	})
	doc, _ = ps.GetById(id)
	if len(doc.Tags()) != 2 {
		t.Fatal("Wrong number of tags")
	}
	if doc.Tags()[0].Value != "one" {
		t.Fatal("Wrong tag value")
	}
	if doc.Tags()[1].Value != "two" {
		t.Fatal("Wrong tag value")
	}
}

func TestAllTags(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)

	tags, err := ps.AllTags()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Fatal("There should be no tags")
	}

	projs, _ := ps.GetAll()
	projs[0].Update(func(pp *ProtoProject) {
		pp.Tags = []tag.ProtoTag{
			{Value: "golang"},
			{Value: "python"},
		}
	})
	projs[1].Update(func(pp *ProtoProject) {
		pp.Tags = []tag.ProtoTag{
			{Value: "golang"},
			{Value: "java"},
			{Value: "svelte"},
		}
	})

	expectedTags := map[string]string{
		"golang": "white",
		"python": "white",
		"java":   "white",
		"svelte": "white",
	}

	tags, err = ps.AllTags()

	if len(tags) != len(expectedTags) {
		t.Fatalf("Expected %d tags, got %d", len(expectedTags), len(tags))
	}

	for _, tag := range tags {
		expectedColor, exists := expectedTags[tag.Value]
		if !exists {
			t.Fatalf("Unexpected tag value: %s", tag.Value)
		}
		if tag.Color != expectedColor {
			t.Fatalf("Expected color for tag %s to be %s, got %s", tag.Value, expectedColor, tag.Color)
		}
	}
}

func TestAllSharedTags(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)

	tags, err := ps.AllSharedTags("golang")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Fatal("There should be no tags")
	}

	projs, _ := ps.GetAll()
	projs[0].Update(func(pp *ProtoProject) {
		pp.Tags = []tag.ProtoTag{
			{Value: "golang"},
			{Value: "python"},
		}
	})
	projs[1].Update(func(pp *ProtoProject) {
		pp.Tags = []tag.ProtoTag{
			{Value: "golang"},
			{Value: "java"},
			{Value: "svelte"},
		}
	})
	projs[2].Update(func(pp *ProtoProject) {
		pp.Tags = []tag.ProtoTag{
			{Value: "lisp"},
			{Value: "lua"},
		}
	})

	expectedTags := map[string]string{
		"golang": "white",
		"python": "white",
		"java":   "white",
		"svelte": "white",
	}

	tags, err = ps.AllSharedTags("golang")

	if len(tags) != len(expectedTags) {
		t.Fatalf("Expected %d tags, got %d", len(expectedTags), len(tags))
	}

	for _, tag := range tags {
		expectedColor, exists := expectedTags[tag.Value]
		if !exists {
			t.Fatalf("Unexpected tag value: %s", tag.Value)
		}
		if tag.Color != expectedColor {
			t.Fatalf("Expected color for tag %s to be %s, got %s", tag.Value, expectedColor, tag.Color)
		}
	}
}

package project

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"samuellando.com/internal/db"
	"samuellando.com/internal/testutil"
	"strings"
	"testing"
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
	doc, _ := ps.GetById(1296269)
	if doc.Title() != "Hello-World" {
		t.Fatal("Ttitle does not match")
	}
}

func TestGetAllBase(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	projects, _ := ps.GetAll()
	if len(projects) != 4 {
		t.Fatal("Should have 4 projects")
	}
}

func TestFilter(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	filtered := ps.Filter(func(d *Project) bool {
		return strings.HasSuffix(d.Title(), "two")
	})
	data, _ := filtered.GetAll()
	if len(data) != 2 {
		t.Errorf("%d should contain 2 elements", len(data))
	}
}

func TestGetIdFiltered(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
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
	ps, ts, db := setup()
	defer teardown(ts, db)
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
	ps, ts, db := setup()
	defer teardown(ts, db)
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
	ps, ts, db := setup()
	defer teardown(ts, db)
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

func TestGetDesc(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	id := 1296269
	doc, _ := ps.GetById(id)
	if doc.Title() != "Hello-World" {
		t.Fatal("Wrong Title")
	}
	if doc.Description() != "This your first repo!" {
		t.Fatalf("Wrong desc %s", doc.Description())
	}
	query := `
    UPDATE project
    SET description = $2
    WHERE id = $1;`
	_, err := db.Exec(query, id, "And your Last!")
	if err != nil {
		t.Fatal(err)
	}
	doc, _ = ps.GetById(id)
	if doc.Description() != "And your Last!" {
		t.Fatalf("Wrong desc %s", doc.Description())
	}
}

func TestGetTags(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	id := 1296269
	doc, _ := ps.GetById(id)
	if doc.Title() != "Hello-World" {
		t.Fatal("Wrong Title")
	}
	if doc.Description() != "This your first repo!" {
		t.Fatalf("Wrong desc %s", doc.Description())
	}
	var id1 int
	var id2 int
	query := `
    INSERT INTO tag (value) VALUES ('one')
    RETURNING id;
    `
	row := db.QueryRow(query)
	err := row.Scan(&id1)
	if err != nil {
		t.Fatal(err)
	}
	query = `
    INSERT INTO tag (value) VALUES ('two')
    RETURNING id;
    `
	row = db.QueryRow(query)
	err = row.Scan(&id2)
	if err != nil {
		t.Fatal(err)
	}
	query = `
    INSERT INTO project_tag (tag, project) VALUES ($2, $1), ($3, $1);
    `
	_, err = db.Exec(query, id, id1, id2)
	if err != nil {
		t.Fatal(err)
	}
	doc, _ = ps.GetById(id)
	if len(doc.Tags()) != 2 {
		t.Fatal("Wrong number of tags")
	}
	if doc.Tags()[0].Value() != "one" {
		t.Fatal("Wrong tag value")
	}
	if doc.Tags()[1].Value() != "two" {
		t.Fatal("Wrong tag value")
	}
}

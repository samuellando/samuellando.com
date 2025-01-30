package document

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"samuellando.com/internal/db"
	"samuellando.com/internal/testutil"
	"github.com/lib/pq"
)

func setup() (Store, *sql.DB) {
	if err := testutil.ResetDb(); err != nil {
		panic(err)
	}
	con := db.ConnectPostgres(testutil.GetDbCredentials())
	migrations, err := testutil.GetMigrationsPath()
	if err != nil {
		panic(err)
	}
	if err := db.ApplyMigrations(con, func(o *db.Options) {
        o.MigrationsDir = migrations
        o.Logger = testutil.CreateDiscardLogger()
    }); err != nil {
		panic(err)
	}
	return CreateStore(con), con
}

func teardown(s Store) {
	s.db.Close()
	testutil.ResetDb()
}

func addDocument(db *sql.DB, title string) int {
	query := `
    INSERT INTO document (title, content) VALUES ($1, 'test')
    RETURNING id;
    `
	row := db.QueryRow(query, title)
	var id int
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func TestGetByIdBase(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	title := "terstingId"
	id := addDocument(db, title)
	doc, _ := ds.GetById(id)
	if doc.Title() != title {
		t.Fatal("Ttitle does not match")
	}
}

func TestGetAllBase(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "test")
	addDocument(db, "test")
	addDocument(db, "test")
	data, _ := ds.GetAll()
	if len(data) != 3 {
		t.Errorf("%d should contain 3 elements", len(data))
	}
}

func TestGetAllUpdates(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "test")
	addDocument(db, "test")
	addDocument(db, "test")
	data, _ := ds.GetAll()
	if len(data) != 3 {
		t.Errorf("%d should contain 3 elements", len(data))
	}
	addDocument(db, "test")
	data, _ = ds.GetAll()
	if len(data) != 4 {
		t.Errorf("%d should contain 4 elements", len(data))
	}
}

func TestAdd(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
    start := time.Now()
    ds.Add(CreateProto(func(df *DocumentFeilds) {
        df.Title = "Test doc"
        df.Content = "Test content"
        df.Tags = []string{"one", "two"}
        df.Created = start
    }))
	query := `
    SELECT 
    d.id AS document_id, 
    d.title, 
    d.content, 
    d.created, 
    array_agg(t.value) AS tags
    FROM 
        document d
    LEFT JOIN 
        document_tag dt ON d.id = dt.document
    LEFT JOIN 
        tag t ON dt.tag = t.id
    GROUP BY d.id, d.title, d.created
    ORDER BY d.created DESC;
    `
    row := db.QueryRow(query)
    var id int
    var title string
    var content string
    var created time.Time
    var tags []sql.NullString
    err := row.Scan(&id, &title, &content, &created, pq.Array(&tags))
    if err != nil {
        t.Error(err)
    }
    if title != "Test doc" {
        t.Error("Wrong title")
    }
    if content != "Test content" {
        t.Error("Wrong content")
    }
    if len(tags) != 2 {
        t.Error("Wrong tag count")
    }
    if tags[0].String != "one" || tags[1].String != "two" {
        t.Error("Wrong tags")
    }
}

func TestRemove(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
    addDocument(db, "test")
    data, err := ds.GetAll()
    if err != nil {
        t.Error(err)
    }
    err = ds.Remove(data[0])
    if err != nil {
        t.Error(err)
    }
    row := db.QueryRow("SELECT count(*) FROM document;")
    var count int
    row.Scan(&count)
    if count > 0 {
        t.Error("Should be empty")
    }
}


func TestFilter(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "abc")
	addDocument(db, "abd")
	addDocument(db, "aaa")
	filtered := ds.Filter(func(d *Document) bool {
		return strings.HasPrefix(d.Title(), "ab")
	})
	data, _ := filtered.GetAll()
	if len(data) != 2 {
		t.Errorf("%d should contain 3 elements", len(data))
	}
}

func TestFilterUpdates(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "abc")
	addDocument(db, "abd")
	addDocument(db, "aaa")
	filtered := ds.Filter(func(d *Document) bool {
		return strings.HasPrefix(d.Title(), "ab")
	})
	data, _ := filtered.GetAll()
	if len(data) != 2 {
		t.Errorf("%d should contain 3 elements", len(data))
	}
	addDocument(db, "abg")
	data, _ = filtered.GetAll()
	if len(data) != 3 {
		t.Errorf("%d should contain 3 elements", len(data))
	}
}

func TestGetIdFiltered(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	in := addDocument(db, "abc")
	addDocument(db, "abd")
	out := addDocument(db, "aaa")
	filtered := ds.Filter(func(d *Document) bool {
		return strings.HasPrefix(d.Title(), "ab")
	})
	elem, err := filtered.GetById(in)
	if err != nil {
		t.Error(err)
	}
	if elem.Title() != "abc" {
		t.Error("Wrong element")
	}
	_, err = filtered.GetById(out)
	if err == nil {
		t.Error("Should not be included")
	}
}

func TestSort(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "d")
	addDocument(db, "b")
	addDocument(db, "a")
	sorted := ds.Sort(func(a, b *Document) bool {
		return strings.Compare(a.Title(), b.Title()) < 0
	})
	data, _ := sorted.GetAll()
	if len(data) != 3 {
		t.Errorf("%d should contain 3 elements", len(data))
	}
	for i, c := range []string{"a", "b", "d"} {
		if data[i].Title() != c {
			t.Error("Out of order")
		}
	}
}

func TestSortUpdates(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "d")
	addDocument(db, "b")
	addDocument(db, "a")
	sorted := ds.Sort(func(a, b *Document) bool {
		return strings.Compare(a.Title(), b.Title()) < 0
	})
	data, _ := sorted.GetAll()
	if len(data) != 3 {
		t.Errorf("%d should contain 3 elements", len(data))
	}
	for i, c := range []string{"a", "b", "d"} {
		if data[i].Title() != c {
			t.Error("Out of order")
		}
	}
	addDocument(db, "c")
	data, _ = sorted.GetAll()
	if len(data) != 4 {
		t.Errorf("%d should contain 4 elements", len(data))
	}
	for i, c := range []string{"a", "b", "c", "d"} {
		if data[i].Title() != c {
			t.Error("Out of order")
		}
	}
}

func TestGroup(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "abb")
	addDocument(db, "acc")
	addDocument(db, "bb")
	addDocument(db, "b")
	addDocument(db, "bxx")
	addDocument(db, "haa")
	addDocument(db, "hb")
	addDocument(db, "hc")
	addDocument(db, "haa")
	addDocument(db, "hb")
	addDocument(db, "hc")
	groups := ds.Group(func(d *Document) string {
		return string(d.Title()[0])
	})
	if len(groups) != 3 {
		t.Error("Wrong number of groups")
	}
	expectedLens := map[string]int{"a": 2, "b": 3, "h": 6}
	for k, s := range groups {
		data, _ := s.GetAll()
		if len(data) != expectedLens[k] {
			t.Errorf("%d should contain %d elements", len(data), expectedLens[k])
		}
	}
}

func TestGroupUpdates(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "abb")
	addDocument(db, "acc")
	addDocument(db, "bb")
	addDocument(db, "b")
	addDocument(db, "bxx")
	addDocument(db, "haa")
	addDocument(db, "hb")
	addDocument(db, "hc")
	addDocument(db, "haa")
	addDocument(db, "hb")
	addDocument(db, "hc")
	groups := ds.Group(func(d *Document) string {
		return string(d.Title()[0])
	})
	if len(groups) != 3 {
		t.Error("Wrong number of groups")
	}
	expectedLens := map[string]int{"a": 2, "b": 3, "h": 6}
	for k, s := range groups {
		data, _ := s.GetAll()
		if len(data) != expectedLens[k] {
			t.Errorf("%d should contain %d elements", len(data), expectedLens[k])
		}
	}
	addDocument(db, "haa")
	expectedLens = map[string]int{"a": 2, "b": 3, "h": 7}
	for k, s := range groups {
		data, _ := s.GetAll()
		if len(data) != expectedLens[k] {
			t.Errorf("%d should contain %d elements", len(data), expectedLens[k])
		}
	}
}

func TestFilterGroupSortAndUpdate(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	addDocument(db, "iac")
	addDocument(db, "iab")
	addDocument(db, "iaa")
	addDocument(db, "xad")
	addDocument(db, "xab")
	addDocument(db, "ibc")
	addDocument(db, "ibb")
	addDocument(db, "iba")
	addDocument(db, "xbc")
	addDocument(db, "xbb")
	addDocument(db, "xba")
	addDocument(db, "icc")
	addDocument(db, "icb")
	addDocument(db, "ica")
	addDocument(db, "xcc")
	addDocument(db, "xcb")
	addDocument(db, "xca")
	addDocument(db, "xco")
	addDocument(db, "rja")
	addDocument(db, "rja")
	addDocument(db, "rja")
	groups := ds.Filter(func(d *Document) bool {
		return strings.Contains(d.Title(), "x")
	}).Sort(func(a, b *Document) bool {
		return strings.Compare(a.Title(), b.Title()) < 0
	}).Group(func(d *Document) string {
		return string(d.Title()[1])
	})
	if len(groups) != 3 {
		t.Error("Wrong number of groups")
	}
	expectedLens := map[string]int{"a": 2, "b": 3, "c": 4}
	for k, s := range groups {
		data, _ := s.GetAll()
		if len(data) != expectedLens[k] {
			t.Errorf("%s, %d should contain %d elements", k, len(data), expectedLens[k])
		}
	}
	ga, _ := groups["a"].GetAll()
	order := []string{"xab", "xad"}
	for i, v := range ga {
		if order[i] != v.Title() {
			t.Error("Wrong ordering")
		}
	}
	addDocument(db, "xac")
	expectedLens = map[string]int{"a": 3, "b": 3, "c": 4}
	for k, s := range groups {
		data, _ := s.GetAll()
		if len(data) != expectedLens[k] {
			t.Errorf("%s, %d should contain %d elements", k, len(data), expectedLens[k])
		}
	}
	ga, _ = groups["a"].GetAll()
	order = []string{"xab", "xac", "xad"}
	for i, v := range ga {
		if order[i] != v.Title() {
			t.Error("Wrong ordering")
		}
	}
    ga[0].Update(func(df *DocumentFeilds) {
        df.Title = "xcz"
    })
	expectedLens = map[string]int{"a": 2, "b": 3, "c": 5}
	for k, s := range groups {
		data, _ := s.GetAll()
		if len(data) != expectedLens[k] {
			t.Errorf("%s, %d should contain %d elements", k, len(data), expectedLens[k])
		}
	}
    gc, _ := groups["c"].GetAll()
    if gc[4].Title() != "xcz" {
        t.Error("Wrong ordering after doc Update")
    }
}

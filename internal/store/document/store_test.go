package document

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"samuellando.com/data"
	"samuellando.com/internal/db"
	"samuellando.com/internal/store/tag"
	"samuellando.com/internal/testutil"
)

func setup() (Store, *sql.DB) {
	con := db.ConnectPostgres(testutil.GetDbCredentials())
	if err := testutil.ResetDb(con, "documentTests"); err != nil {
		panic(err)
	}
	return CreateStore(con), con
}

func setupSampleData(con *sql.DB) {
	ctx := context.TODO()
	tx, err := con.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		panic(err)
	}
	queries := data.New(con).WithTx(tx)
	id1, err := queries.CreateDocument(ctx, data.CreateDocumentParams{
		Title:   "First Document",
		Content: "Sample",
		Created: time.Date(2007, 07, 07, 0, 0, 0, 0, &time.Location{}),
	})
	if err != nil {
		panic(err)
	}
	id2, err := queries.CreateDocument(ctx, data.CreateDocumentParams{
		Title:   "Second Document",
		Content: "Sample test 2",
		Created: time.Now(),
	})
	if err != nil {
		panic(err)
	}
	_, err = queries.SetDocumentTags(ctx, data.SetDocumentTagsParams{
		Document:  id1,
		TagValues: []string{"golang", "python"},
	})
	if err != nil {
		panic(err)
	}
	_, err = queries.SetDocumentTags(ctx, data.SetDocumentTagsParams{
		Document:  id2,
		TagValues: []string{"java", "svelte"},
	})
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func teardown(con *sql.DB) {
	con.Close()
}

func TestGetById(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	setupSampleData(con)
	doc, err := ds.GetById(2)
	if err != nil {
		t.Fatal(err)
	}
	// Check the all the values
	if doc.Id() != 2 {
		t.Fatalf("Expected id to be 2 for the second document, got [%d]", doc.Id())
	}
	if doc.Title() != "Second Document" {
		t.Fatalf("Expected title to be 'Second Document', got [%s]", doc.Title())
	}
	if doc.Content() != "Sample test 2" {
		t.Fatalf("Expected content to be 'Sample test 2', got [%s]", doc.Content())
	}
	if len(doc.Tags()) != 2 || doc.Tags()[0].Value != "java" || doc.Tags()[1].Value != "svelte" {
		t.Fatalf("Expected tags to be ['java', 'svelte'], got [%v]", doc.Tags())
	}
	if doc.Tags()[0].Color != "white" || doc.Tags()[1].Color != "white" {
		t.Fatalf("Expected tag colors to be ['white', 'white'], got [%v]", doc.Tags())
	}
}

func TestGetByIdNotFound(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	setupSampleData(con)
	_, err := ds.GetById(3)
	if err == nil {
		t.Fatal("Should fail if doc is not found")
	}
}

func TestGetAll(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	setupSampleData(con)
	docs, err := ds.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 2 {
		t.Fatalf("Expected 2 documents [%d]", len(docs))
	}
	// Now sort the docs by id
	sorted, err := ds.Sort(func(d1, d2 Document) bool {
		return d1.Id() < d2.Id()
	})
	if err != nil {
		t.Fatal(err)
	}
	docs, _ = sorted.GetAll()
	// Check the first document
	if docs[0].Id() != 1 {
		t.Fatalf("Expected id to be 1 for the first document, got [%d]", docs[0].Id())
	}
	if docs[0].Title() != "First Document" {
		t.Fatalf("Expected title to be 'First Document', got [%s]", docs[0].Title())
	}
	if docs[0].Content() != "Sample" {
		t.Fatalf("Expected content to be 'Sample', got [%s]", docs[0].Content())
	}
	if !docs[0].Created().Equal(time.Date(2007, 07, 07, 0, 0, 0, 0, &time.Location{})) {
		t.Fatalf("Expected created date to be 2007-07-07, got [%s]", docs[0].Created())
	}
	if len(docs[0].Tags()) != 2 || docs[0].Tags()[0].Value != "golang" || docs[0].Tags()[1].Value != "python" {
		t.Fatalf("Expected tags to be ['golang', 'python'], got [%v]", docs[0].Tags())
	}
	if docs[0].Tags()[0].Color != "white" || docs[0].Tags()[1].Color != "white" {
		t.Fatalf("Expected tag colors to be ['white', 'white'], got [%v]", docs[0].Tags())
	}

	// Check the second document
	if docs[1].Id() != 2 {
		t.Fatalf("Expected id to be 2 for the second document, got [%d]", docs[1].Id())
	}
	if docs[1].Title() != "Second Document" {
		t.Fatalf("Expected title to be 'Second Document', got [%s]", docs[1].Title())
	}
	if docs[1].Content() != "Sample test 2" {
		t.Fatalf("Expected content to be 'Sample test 2', got [%s]", docs[1].Content())
	}
	if len(docs[1].Tags()) != 2 || docs[1].Tags()[0].Value != "java" || docs[1].Tags()[1].Value != "svelte" {
		t.Fatalf("Expected tags to be ['java', 'svelte'], got [%v]", docs[1].Tags())
	}
	if docs[1].Tags()[0].Color != "white" || docs[1].Tags()[1].Color != "white" {
		t.Fatalf("Expected tag colors to be ['white', 'white'], got [%v]", docs[1].Tags())
	}
}

func TestGetAllNone(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	docs, err := ds.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 0 {
		t.Fatalf("Expected 0 documents [%d]", len(docs))
	}
}

func TestAdd(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	setupSampleData(con)
	c, err := ds.Add(ProtoDocument{
		Title:   "New Doc",
		Content: "Sample content",
		Created: time.Now(),
		Tags: []tag.ProtoTag{
			{Value: "c"},
			{Value: "c++"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Id() != 3 {
		t.Fatalf("Expected id to be 3 for the second cument, got [%d]", c.Id())
	}
	if c.Title() != "New Doc" {
		t.Fatalf("Expected title to be 'Second Document', got [%s]", c.Title())
	}
	if c.Content() != "Sample content" {
		t.Fatalf("Expected content to be 'Sample test 2', got [%s]", c.Content())
	}
	if len(c.Tags()) != 2 || c.Tags()[0].Value != "c" || c.Tags()[1].Value != "c++" {
		t.Fatalf("Expected tags to be ['java', 'svelte'], got [%v]", c.Tags())
	}
	if c.Tags()[0].Color != "white" || c.Tags()[1].Color != "white" {
		t.Fatalf("Expedocted tag docolors to be ['white', 'white'], got [%v]", c.Tags())
	}
	doc, err := ds.GetById(c.Id())
	if err != nil {
		t.Fatal(err)
	}
	if doc.Id() != 3 {
		t.Fatalf("Expedocted id to be 3 for the sedocond document, got [%d]", doc.Id())
	}
	if doc.Title() != "New Doc" {
		t.Fatalf("Expedocted title to be 'Sedocond Dodocument', got [%s]", doc.Title())
	}
	if doc.Content() != "Sample content" {
		t.Fatalf("Expedocted docontent to be 'Sample test 2', got [%s]", doc.Content())
	}
	if len(doc.Tags()) != 2 || doc.Tags()[0].Value != "c" || doc.Tags()[1].Value != "c++" {
		t.Fatalf("Expedocted tags to be ['c', 'c++'], got [%v]", doc.Tags())
	}
	if doc.Tags()[0].Color != "white" || doc.Tags()[1].Color != "white" {
		t.Fatalf("Expedocted tag docolors to be ['white', 'white'], got [%v]", doc.Tags())
	}
}

func TestAllTags(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	setupSampleData(con)

	tags, err := ds.AllTags()
	if err != nil {
		t.Fatal(err)
	}

	expectedTags := map[string]string{
		"golang": "white",
		"python": "white",
		"java":   "white",
		"svelte": "white",
	}

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
	ds, con := setup()
	defer teardown(con)
	setupSampleData(con)
	_, err := ds.Add(ProtoDocument{
		Title:   "New Doc",
		Content: "Sample content",
		Created: time.Now(),
		Tags: []tag.ProtoTag{
			{Value: "java"},
			{Value: "c"},
			{Value: "c++"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	sharedTags, err := ds.AllSharedTags("java")
	if err != nil {
		t.Fatal(err)
	}

	expectedSharedTags := map[string]string{
		"java":   "white",
		"c":      "white",
		"c++":    "white",
		"svelte": "white",
	}

	if len(sharedTags) != len(expectedSharedTags) {
		t.Fatalf("Expected %d shared tags, got %d", len(expectedSharedTags), len(sharedTags))
	}

	for _, tag := range sharedTags {
		expectedColor, exists := expectedSharedTags[tag.Value]
		if !exists {
			t.Fatalf("Unexpected shared tag value: %s", tag.Value)
		}
		if tag.Color != expectedColor {
			t.Fatalf("Expected color for shared tag %s to be %s, got %s", tag.Value, expectedColor, tag.Color)
		}
	}
}

package document

import (
	"testing"
	"time"

	"samuellando.com/internal/store/tag"
)

func TestUpdateProvidesDefaults(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	doc, err := ds.Add(ProtoDocument{
		Title:   "Sample",
		Content: "Content",
		Created: time.Now(),
		Tags: []tag.ProtoTag{
			{Value: "one"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	var tagsLeak []tag.ProtoTag
	err = doc.Update(func(pd *ProtoDocument) {
		if pd.Title != "Sample" {
			t.Fatal("Title no defaulted")
		}
		if pd.Content != "Content" {
			t.Fatal("Content not defaulted")
		}
		if time.Since(pd.Created) > time.Second {
			t.Fatal("Created no defaulted")
		}
		if len(pd.Tags) != 1 || pd.Tags[0].Value != "one" {
			t.Fatal("Tags no defaulted")
		}
		tagsLeak = pd.Tags
	})
	if err != nil {
		t.Fatal(err)
	}
	tagsLeak[0] = tag.ProtoTag{Value: "Two"}
	if doc.Tags()[0].Value != "one" {
		t.Fatal("tag array leaks to old document")
	}
}

func TestUpdate(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	doc, err := ds.Add(ProtoDocument{
		Title:   "Sample",
		Content: "Content",
		Created: time.Now(),
		Tags: []tag.ProtoTag{
			{Value: "one"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = doc.Update(func(pd *ProtoDocument) {
		pd.Title = "Sample2"
		pd.Content = "Content2"
		pd.Created = time.Date(2007, 07, 07, 01, 01, 01, 01, &time.Location{})
		pd.Tags = []tag.ProtoTag{
			{Value: "two"},
			{Value: "three"},
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title() != "Sample2" {
		t.Fatalf("expected updated title to be 'Sample2', got '%s'", doc.Title())
	}
	if doc.Content() != "Content2" {
		t.Fatalf("expected updated content to be 'Content2', got '%s'", doc.Content())
	}
	expectedTime := time.Date(2007, 07, 07, 01, 01, 01, 01, &time.Location{})
	if !doc.Created().Equal(expectedTime) {
		t.Fatalf("expected updated created time to be '%v', got '%v'", expectedTime, doc.Created())
	}
	if len(doc.Tags()) != 2 || doc.Tags()[0].Value != "three" || doc.Tags()[1].Value != "two" {
		t.Fatalf("expected updated tags to be '[two, three]', got '%v'", doc.Tags())
	}
	err = doc.Update(func(pd *ProtoDocument) {
		pd.Tags = []tag.ProtoTag{}
	})
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title() != "Sample2" {
		t.Fatalf("expected updated title to be 'Sample2', got '%s'", doc.Title())
	}
	if doc.Content() != "Content2" {
		t.Fatalf("expected updated content to be 'Content2', got '%s'", doc.Content())
	}
	expectedTime = time.Date(2007, 07, 07, 01, 01, 01, 01, &time.Location{})
	if !doc.Created().Equal(expectedTime) {
		t.Fatalf("expected updated created time to be '%v', got '%v'", expectedTime, doc.Created())
	}
	if len(doc.Tags()) != 0 {
		t.Fatalf("expected updated tags to be '[two, three]', got '%v'", doc.Tags())
	}
}

func TestDelete(t *testing.T) {
	ds, con := setup()
	defer teardown(con)
	setupSampleData(con)
	doc, err := ds.GetById(1)
	if err != nil {
		t.Fatal(err)
	}
	err = doc.Delete()
	if err != nil {
		t.Fatal(err)
	}
	doc, err = ds.GetById(1)
	if err == nil {
		t.Fatal("Should throw an error")
	}
}

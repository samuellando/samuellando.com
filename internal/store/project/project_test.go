package project

import (
	"samuellando.com/internal/store/tag"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)

	proj, _ := ps.GetById(1)

	if proj.Title() != "Bye-World-two" {
		t.Fatalf("Expected title 'Bye-World-two', got '%s'", proj.Title())
	}
	if proj.Description() != "This your first repo!" {
		t.Fatalf("Expected description 'This your first repo!', got '%s'", proj.Description())
	}
	if proj.Url() != "https://github.com/octocat/Hello-World" {
		t.Fatalf("Expected url %s", proj.Url())
	}
	expectedTime, _ := time.Parse(time.RFC3339, "2011-01-26T19:01:12Z")
	if !proj.Created().Equal(expectedTime) {
		t.Fatalf("Expected created %v", proj.Created())
	}
	expectedTime, _ = time.Parse(time.RFC3339, "2011-01-26T19:06:43Z")
	if !proj.Pushed().Equal(expectedTime) {
		t.Fatal("Expected pushed")
	}
	if len(proj.Tags()) != 0 {
		t.Fatalf("Expected 0 tags, got %d", len(proj.Tags()))
	}

	// Perform an update
	newDesc := "Updated description"
	newTags := []tag.ProtoTag{
		{Value: "updated-tag-1"},
		{Value: "updated-tag-2"},
	}
	err := proj.Update(func(pp *ProtoProject) {
		pp.Description = &newDesc
		pp.Tags = newTags
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check the inmemory data
	if proj.Title() != "Bye-World-two" {
		t.Fatalf("Expected title 'Bye-World-two', got '%s'", proj.Title())
	}
	if proj.Description() != newDesc {
		t.Fatalf("Expected description '%s', got '%s'", newDesc, proj.Description())
	}
	if len(proj.Tags()) != len(newTags) {
		t.Fatalf("Expected %d tags, got %d", len(newTags), len(proj.Tags()))
	}
	for i, tag := range proj.Tags() {
		if tag.Value != newTags[i].Value {
			t.Fatalf("Expected tag %d to be %v, got %v", i, newTags[i], tag)
		}
	}

	// Retrieve the project again
	proj, _ = ps.GetById(1)

	// Check updated values
	if proj.Title() != "Bye-World-two" {
		t.Fatalf("Expected title 'Bye-World-two', got '%s'", proj.Title())
	}
	if proj.Description() != newDesc {
		t.Fatalf("Expected description '%s', got '%s'", newDesc, proj.Description())
	}
	if len(proj.Tags()) != len(newTags) {
		t.Fatalf("Expected %d tags, got %d", len(newTags), len(proj.Tags()))
	}
	for i, tag := range proj.Tags() {
		if tag.Value != newTags[i].Value {
			t.Fatalf("Expected tag %d to be %v, got %v", i, newTags[i], tag)
		}
	}

	err = proj.Update(func(pp *ProtoProject) {
		pp.Description = nil
		pp.Tags = []tag.ProtoTag{}
	})
	if err != nil {
		t.Fatal(err)
	}
	if proj.Description() != "" {
		t.Fatalf("Expected description 'This your first repo!', got '%s'", proj.Description())
	}
	if len(proj.Tags()) != 0 {
		t.Fatalf("Expected 0 tags, got %d", len(proj.Tags()))
	}
	proj, _ = ps.GetById(1)
	if proj.Description() != "This your first repo!" {
		t.Fatalf("Expected description 'This your first repo!', got '%s'", proj.Description())
	}
	if len(proj.Tags()) != 0 {
		t.Fatalf("Expected 0 tags, got %d", len(proj.Tags()))
	}
}

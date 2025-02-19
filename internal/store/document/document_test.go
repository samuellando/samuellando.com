package document

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
)

func TestCreateProtoCopiesTags(t *testing.T) {
	title := "Test title"
	content := "Test Content"
	tags := []string{"one", "two"}
	created := time.Now()
	doc := CreateProto(func(df *DocumentFeilds) {
		df.Title = title
		df.Content = content
		df.Tags = tags
		df.Created = created
	})
	docTags := doc.Tags()
	docTags = append(docTags, "three")
	if len(tags) != 2 {
		t.Fatal("The original array should unaffected")
	}
}

func TestCreateProtoAndGet(t *testing.T) {
	title := "Test title"
	content := "Test Content"
	tags := []string{"one", "two"}
	created := time.Now()
	doc := CreateProto(func(df *DocumentFeilds) {
		df.Title = title
		df.Content = content
		df.Tags = tags
		df.Created = created
	})
	if doc.Id() != -1 {
		t.Fatal("Proto docs should have id -1")
	}
	if doc.Title() != title {
		t.Fatal("title is wrong")
	}
	docContent, _ := doc.Content()
	if docContent != content {
		t.Fatal("Content is wrong")
	}
	if doc.Created() != created {
		t.Fatal("created is wrong")
	}
	docTags := doc.Tags()
	if len(docTags) != 2 {
		t.Fatal("Expected 2 tags")
	}
}

func TestCreateProtoDefaults(t *testing.T) {
	doc := CreateProto()
	if doc.Id() != -1 {
		t.Fatal("Proto docs should have id -1")
	}
	if doc.Title() != "" {
		t.Fatal("title is wrong")
	}
	docContent, _ := doc.Content()
	if docContent != "" {
		t.Fatal("Content is wrong")
	}
	if time.Since(doc.Created()) < 0 || time.Since(doc.Created()) > time.Second {
		t.Fatal("created is wrong")
	}
	docTags := doc.Tags()
	if len(docTags) != 0 {
		t.Fatal("Expected 0 tags")
	}
}

func TestCreateProtoAndAdd(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	doc := CreateProto()
	ds.Add(doc)
	if doc.Id() == -1 {
		t.Fatal("The doc should have an Id now")
	}
	query := "SELECT count(*) FROM document;"
	row := db.QueryRow(query)
	var count int
	row.Scan(&count)
	if count != 1 {
		t.Fatal("The doc should be in the DB")
	}
}

func TestCreateProtoAndUpdate(t *testing.T) {
	title := "Test title"
	doc := CreateProto()
	err := doc.Update(func(df *DocumentFeilds) {
		df.Title = title
	})
	if err == nil {
		t.Fatal("Should not be possible to update a proto item")
	}
}

func TestCreateAddAndUpdate(t *testing.T) {
	ds, db := setup()
	defer teardown(ds)
	title := "Test title"
	content := "Test Content"
	tags := []string{"one", "two"}
	created := time.UnixMilli(100)
	doc := CreateProto()
	ds.Add(doc)
	err := doc.Update(func(df *DocumentFeilds) {
		df.Title = title
		df.Content = content
		df.Tags = tags
		df.Created = created
	})
	if err != nil {
		t.Fatal("Should be possible to update a addded item")
	}
	if doc.Title() != title {
		t.Fatal("The title should be updated")
	}
	docContent, _ := doc.Content()
	if docContent != content {
		t.Fatal("Content should be udated")
	}
	if doc.Created() != created {
		t.Fatal("created should be updated")
	}
	docTags := doc.Tags()
	if len(docTags) != 2 {
		t.Fatal("Expected 2 tags")
	}
	query := `
    SELECT 
        d.id,
        d.title, 
        d.content,
        d.created, 
        array_agg(t.value)
    FROM document d
    LEFT JOIN document_tag dt ON d.id = dt.document 
    LEFT JOIN tag t ON dt.tag = t.id
    GROUP BY d.id, d.title, d.content, d.created;
    `
	row := db.QueryRow(query)
	var dbId int
	var dbTitle string
	var dbContent string
	var dbCreated time.Time
	var dbTags []sql.NullString
	err = row.Scan(&dbId, &dbTitle, &dbContent, &dbCreated, pq.Array(&dbTags))
	if err != nil {
		t.Fatal(err)
	}
	if dbTitle != title {
		t.Fatal("db The title should be updated")
	}
	if dbContent != content {
		t.Fatal("db Content should be udated")
	}
	if !dbCreated.Equal(created) {
		t.Fatal("db created should be updated")
	}
	if len(dbTags) != 2 {
		t.Fatal("db Expected 2 tags")
	}
}

func TestUpdateCopies(t *testing.T) {
	ds, _ := setup()
	defer teardown(ds)
	title := "Test title"
	tags := []string{"one", "two"}
	doc := CreateProto()
	ds.Add(doc)
	var stolen *DocumentFeilds
	err := doc.Update(func(df *DocumentFeilds) {
		df.Title = title
		df.Tags = tags
		stolen = df
	})
	if err != nil {
		t.Fatal(err)
	}
	tags[0] = "Zero"
	docTags := doc.Tags()
	if docTags[0] != "one" {
		t.Fatal("The array should unaffected")
	}
	stolen.Title = "New title"
	if doc.Title() != title {
		t.Fatal("The stoen fields object leaks")
	}
}

func TestUpdateRollsback(t *testing.T) {
	ds, _ := setup()
	defer teardown(ds)
	title := "Test title"
	tags := []string{"one", "two"}
	doc := CreateProto()
	ds.Add(doc)
	var stolen *DocumentFeilds
	doc.id = -1 // Will cause a query error
	err := doc.Update(func(df *DocumentFeilds) {
		df.Title = title
		df.Tags = tags
		stolen = df
	})
	if err == nil {
		t.Fatal("Should fail")
	}
	docTags := doc.Tags()
	stolen.Tags[0] = "Zero"
	if len(docTags) != 0 {
		t.Fatal("The array should unaffected")
	}
	stolen.Title = "New title"
	if doc.Title() != "" {
		t.Fatal("The stolen fields object leaks")
	}
}

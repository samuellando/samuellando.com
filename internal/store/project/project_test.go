package project

import (
	"database/sql"
	"testing"

	"github.com/lib/pq"
	"samuellando.com/internal/store/tag"
)

func TestUpdate(t *testing.T) {
	ps, ts, db := setup()
	defer teardown(ts, db)
	doc, _ := ps.GetById(1296269)
	arr := []tag.Tag{
		tag.CreateProto(func(tf *tag.TagFields) { tf.Value = "one" }),
		tag.CreateProto(func(tf *tag.TagFields) { tf.Value = "two" }),
	}
	err := doc.Update(func(pf *ProjectFields) {
		pf.Description = "Testing"
		pf.Tags = arr
	})
	if err != nil {
		t.Fatal(err)
	}
	if doc.Description() != "Testing" {
		t.Fatal("Desc should be updated in the object")
	}
	if len(doc.Tags()) != 2 {
		t.Fatal("Tag should be in the object")
	}
	query := `
    SELECT 
        p.description,  
        array_agg(t.value) AS tags
    FROM 
        project p
    LEFT JOIN 
        project_tag pt ON p.id = pt.project
    LEFT JOIN 
        tag t ON t.id = pt.tag
    WHERE p.id = $1
    GROUP BY p.description;
    `
	var desc sql.NullString
	var tags []sql.NullString
	row := db.QueryRow(query, doc.Id())
	err = row.Scan(&desc, pq.Array(&tags))
	if err != nil {
		t.Fatal(err)
	}
	if desc.String != "Testing" {
		t.Fatalf("Desc shold be updated in the db %s", desc.String)
	}
	for i, tag := range tags {
		if tag.String != arr[i].Value() {
			t.Fatalf("Tag shold be updated in db %d %s", i, tag.String)
		}
	}
}

// This package provides the Document type for use with the DocumentStore.
//
// This type ensures consitency with the database by making the type immutable
// And only modifiable with the Update method.
package document

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"samuellando.com/internal/markdown"
)

// The fields a document contains
type DocumentFeilds struct {
	Title   string
	Content string
	Tags    []string
	Created time.Time
}

// An actual document in the database, or a proto document (not cerated or deleted)
type Document struct {
	db     *sql.DB
	loaded bool
	id     int
	fields DocumentFeilds
}

func (d *Document) Id() int {
	return d.id
}

func (d *Document) Title() string {
	return d.fields.Title
}

// Dynamically loads the content from the database if it's not loaded yet
func (d *Document) Content() (string, error) {
	if !d.loaded {
		if err := d.loadContent(); err != nil {
			return "", err
		}
	}
	return d.fields.Content, nil
}

func (d *Document) Html() (template.HTML, error) {
    content, err := d.Content() 
    if err != nil {
        return "", err
    }
    return markdown.ToHtml(content)
} 

func (d *Document) Tags() []string {
	return copyOf(d.fields.Tags)
}

func (d *Document) Created() time.Time {
	return d.fields.Created
}

// Creates a document that does not yet exist in a DocumentStore and must be added with
// DocumentStore.Add
//
// Once the protoDoc is created is is immutable (everything is deep copied) until
// it is added to the Document Store and updated with the Update method.
func CreateProto(setters ...func(*DocumentFeilds)) *Document {
	docFields := DocumentFeilds{
		Title:   "",
		Content: "",
		Tags:    []string{},
		Created: time.Now(),
	}
	for _, set := range setters {
		set(&docFields)
	}
	docFields.Tags = copyOf(docFields.Tags)
	doc := Document{db: nil, id: -1, loaded: true, fields: docFields}
	return &doc
}

// Update a document
//
// This does not work for proto documents
//
// everything is deep copied, and rolled back in case of an error
func (d *Document) Update(setters ...func(*DocumentFeilds)) error {
	if d.db == nil {
		return fmt.Errorf("Cannot update a proto document")
	}
    original := d.fields
	c := d.fields
    c.Tags = copyOf(c.Tags)
	for _, set := range setters {
		set(&c)
	}
	// Make a copy of everything
	d.fields = c
	d.fields.Tags = copyOf(c.Tags)
    if d.fields.Content != "" {
        d.loaded = true
    }
	// and update in the db
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		tx.Rollback()
		if err != nil {
			d.fields = original
		}
	}()
	query := `
    UPDATE document SET 
        title = $1,
        content = $2,
        created = $3
    WHERE 
        id = $4;`
	content, err := d.Content()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, d.Title(), content, d.Created(), d.Id())
	err = d.setTags(tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (d *Document) setTags(tx *sql.Tx) error {
	// Create the tags if they are missingue
	tagIds := make([]int, len(d.Tags()))
	query := `
    INSERT INTO tag (value) 
    VALUES ($1) 
    ON CONFLICT (value) DO UPDATE 
    SET value = tag.value
    RETURNING id;
    `
	for i, tag := range d.Tags() {
		row := tx.QueryRow(query, tag)
		err := row.Scan(&tagIds[i])
		if err != nil {
			return fmt.Errorf("Failed to create tag: %w", err)
		}
	}
	// Reset document - tag associations
	query = "DELETE FROM document_tag WHERE document = $1;"
	_, err := tx.Exec(query, d.Id())
	if err != nil {
		return err
	}
	query = "INSERT INTO document_tag (document, tag) VALUES ($1, $2);"
	for i := range d.Tags() {
		_, err := tx.Exec(query, d.Id(), tagIds[i])
		if err != nil {
			return fmt.Errorf("Failed to apply tag: %w", err)
		}
	}
	return nil
}

func (d *Document) loadContent() error {
	if d.db == nil {
		return fmt.Errorf("Can't load contant of a proto document")
	}
	query := "SELECT content FROM document WHERE id = $1"
	row := d.db.QueryRow(query, d.id)
	err := row.Scan(&d.fields.Content)
	if err != nil {
		return err
	}
	d.loaded = true
	return nil
}

func copyOf(src []string) []string {
	tagsCopy := make([]string, len(src))
	copy(tagsCopy, src)
	return tagsCopy
}

func (d *Document) ToString() string {
    content, err := d.Content()
    if err != nil {
        content = ""
    }
    s := d.Title() + "\n" +  content
    return s
}


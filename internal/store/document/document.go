// This package provides the Document type for use with the DocumentStore.
//
// This type ensures consitency with the database by making the type immutable
// And only modifiable with the Update method.
package document

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"strings"
	"time"

	"samuellando.com/data"
	"samuellando.com/internal/markdown"
	"samuellando.com/internal/store/tag"
)

// A prototype document used for creating documents.
type ProtoDocument struct {
	Title   string
	Content string
	Tags    []tag.ProtoTag
	Created time.Time
}

// An actual document in the database.
type Document struct {
	db      *sql.DB
	id      int64
	title   string
	content string
	tags    []tag.ProtoTag
	created time.Time
}

func (d Document) Id() int64 {
	return d.id
}

func (d Document) Title() string {
	return d.title
}

// Dynamically loads the content from the database if it's not loaded yet
func (d Document) Content() string {
	return d.content
}

func (d Document) Html() (template.HTML, error) {
	content := d.Content()
	return markdown.ToHtml(content)
}

func (d Document) Tags() []tag.ProtoTag {
	return copyTags(d.tags)
}

func (d Document) Created() time.Time {
	return d.created
}

func (d Document) ToString() string {
	content := d.Content()
	s := fmt.Sprintf("%s\n%s\n%s", d.Title(), content, strings.Join(tagValues(d.Tags()), " "))
	return s
}

// Update a document
//
// everything is deep copied, and rolled back in case of an error
func (d *Document) Update(setters ...func(*ProtoDocument)) error {
	p := ProtoDocument{
		Title:   d.Title(),
		Content: d.Content(),
		Created: d.Created(),
		Tags:    d.Tags(),
	}
	for _, setter := range setters {
		setter(&p)
	}

	ctx := context.TODO()
	tx, err := d.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	queries := data.New(d.db).WithTx(tx)
	err = queries.UpdateDocument(ctx, data.UpdateDocumentParams{
		ID:      d.Id(),
		Title:   p.Title,
		Content: p.Content,
		Created: p.Created,
	})
	if err != nil {
		return err
	}
	tagRows, err := queries.SetDocumentTags(ctx, data.SetDocumentTagsParams{
		Document:  d.id,
		TagValues: tagValues(p.Tags),
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	tags := make([]tag.ProtoTag, len(tagRows))
	for i, tagRow := range tagRows {
		tags[i] = tag.ProtoTag{
			Value: tagRow.Value,
			Color: tagRow.Color,
		}
	}
	d.title = p.Title
	d.content = p.Content
	d.created = p.Created
	d.tags = tags
	return nil
}

func (d Document) Delete() error {
	ctx := context.TODO()
	queries := data.New(d.db)
	err := queries.DeleteDocument(ctx, d.id)
	return err
}

func tagValues(src []tag.ProtoTag) []string {
	s := make([]string, len(src))
	for i, tag := range src {
		s[i] = tag.Value
	}
	return s
}

func copyTags(in []tag.ProtoTag) []tag.ProtoTag {
	tagsCopy := make([]tag.ProtoTag, len(in))
	copy(tagsCopy, in)
	return tagsCopy
}

package project

import (
	"database/sql"
	"fmt"
	"samuellando.com/internal/store/tag"
	"strings"
	"time"
)

type Project struct {
	db             *sql.DB
	id             int
	name           string
	gitDescription string
	created        time.Time
	pushed         time.Time
	url            string
	fields         ProjectFields
}

type ProjectFields struct {
	Description string
	Tags        []tag.Tag
}

func (p *Project) Id() int {
	return p.id
}

func (p *Project) Title() string {
	return p.name
}

func (p *Project) Description() string {
	d := p.fields.Description
	if d == "" {
		return p.gitDescription
	} else {
		return d
	}
}

func (p *Project) Created() time.Time {
	return p.created
}

func (p *Project) Pushed() time.Time {
	return p.pushed
}

func (p *Project) Url() string {
	return p.url
}

func (p *Project) Tags() []tag.Tag {
	return p.fields.Tags
}

func (p *Project) Update(opts ...func(*ProjectFields)) error {
	if p.db == nil {
		return fmt.Errorf("Cannot update a proto document")
	}
	original := p.fields
	c := p.fields
	c.Tags = copyOf(c.Tags)
	for _, set := range opts {
		set(&c)
	}
	// Make a copy of everything
	p.fields = c
	p.fields.Tags = copyOf(c.Tags)
	// and update in the db
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		tx.Rollback()
		if err != nil {
			p.fields = original
		}
	}()
	query := `
    UPDATE project SET 
        description = $1
    WHERE 
        id = $2;`
	_, err = tx.Exec(query, p.Description(), p.Id())
	if err != nil {
		return err
	}
	err = p.setTags(tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (p *Project) setTags(tx *sql.Tx) error {
	// Create the tags if they are missingue
	tagIds := make([]int, len(p.Tags()))
	query := `
    INSERT INTO tag (value) 
    VALUES ($1) 
    ON CONFLICT (value) DO UPDATE 
    SET value = tag.value
    RETURNING id;
    `
	for i, tag := range p.Tags() {
		row := tx.QueryRow(query, tag.Value())
		err := row.Scan(&tagIds[i])
		if err != nil {
			return fmt.Errorf("Failed to create tag: %w", err)
		}
	}
	// Reset document - tag associations
	query = "DELETE FROM project_tag WHERE project = $1;"
	_, err := tx.Exec(query, p.Id())
	if err != nil {
		return err
	}
	query = "INSERT INTO project_tag (project, tag) VALUES ($1, $2);"
	for i := range p.Tags() {
		_, err := tx.Exec(query, p.Id(), tagIds[i])
		if err != nil {
			return fmt.Errorf("Failed to apply tag: %w", err)
		}
	}
	return nil
}

func loadProject(db *sql.DB, s *schema) (*Project, error) {
	p := &Project{
		db:             db,
		id:             s.ID,
		name:           s.Name,
		gitDescription: s.Description,
		created:        s.CreatedAt,
		pushed:         s.PushedAt,
		url:            s.HTMLURL,
	}
	// Grab local data from the database.
	query := `
    SELECT description
    FROM project
    WHERE id = $1
    `
	row := db.QueryRow(query, p.Id())
	var desc sql.NullString
	err := row.Scan(&desc)
	if err != nil {
		desc = sql.NullString{Valid: false}
		// And add it to the table
		query := `
        INSERT INTO project (id) 
        VALUES ($1) 
        ON CONFLICT (id) DO NOTHING;
        `
		_, err := db.Exec(query, p.Id())
		if err != nil {
			return nil, err
		}
	}
	pDesc := ""
	if desc.Valid {
		pDesc = desc.String
	}
	query = `
    SELECT value, color
    FROM tag t
    LEFT JOIN project_tag pt ON pt.tag = t.id
    LEFT JOIN project p ON pt.project = p.id
    WHERE p.id = $1
    `
	tagRows, err := db.Query(query, p.Id())
	pTags := make([]tag.Tag, 0)
	for tagRows.Next() {
		tag := tag.CreateProto(func(tf *tag.TagFields) {
			tagRows.Scan(&tf.Value, &tf.Color)
		})
		pTags = append(pTags, tag)
	}
	fields := ProjectFields{
		Description: pDesc,
		Tags:        pTags,
	}
	p.fields = fields
	return p, nil
}

func copyOf(src []tag.Tag) []tag.Tag {
	tagsCopy := make([]tag.Tag, len(src))
	copy(tagsCopy, src)
	return tagsCopy
}

func values(src []tag.Tag) []string {
	s := make([]string, len(src))
	for i, tag := range src {
		s[i] = tag.Value()
	}
	return s
}

func (p *Project) ToString() string {
	s := fmt.Sprintf("%s\n%s\n%s", p.Title(), p.Description(), strings.Join(values(p.Tags()), " "))
	return s
}

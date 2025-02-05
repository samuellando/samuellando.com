package project

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
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
	Tags        []string
}

func (p *Project) Id() int {
	return p.id
}

func (p *Project) Title() string {
	return p.name
}

func (p *Project) GitDescription() string {
	return p.gitDescription
}

func (p *Project) Description() string {
	d := p.fields.Description
	if d == "" {
		return p.GitDescription()
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

func (p *Project) Tags() []string {
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
        Description = $1
    WHERE 
        id = $2;`
	_, err = tx.Exec(query, p.Description(), p.Id())
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
		row := tx.QueryRow(query, tag)
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

func createProject(db *sql.DB, s *schema) (*Project, error) {
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
	row := db.QueryRow(query, p.Id())
	var desc sql.NullString
	var tags []sql.NullString
	err := row.Scan(&desc, pq.Array(&tags))
	if err != nil {
		desc = sql.NullString{Valid: false}
		tags = []sql.NullString{}
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
	pTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag.Valid {
			pTags = append(pTags, tag.String)
		}
	}
	fields := ProjectFields{
		Description: pDesc,
		Tags:        pTags,
	}
	p.fields = fields
	return p, nil
}

func copyOf(src []string) []string {
	tagsCopy := make([]string, len(src))
	copy(tagsCopy, src)
	return tagsCopy
}

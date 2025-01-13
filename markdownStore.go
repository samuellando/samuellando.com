package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "postgres"
)

type Document interface {
	Id() int
	Tags() []string
	SetTags([]string) error
	Title() string
	SetTitle(string) error
	Content() string
	SetContent(string) error
	Created() time.Time
	Published() bool
	SetPublished(bool) error
	Delete() error
}

type MarkdownStore interface {
	GetDocumentById(int) (Document, error)
	GetDocumentsByTag(string) ([]Document, error)
	Documents() ([]Document, error)
	Drafts() ([]Document, error)
	CreateDocument(string, string, []string) (Document, error)
	Close()
}

func createDbConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func initializeModel(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw the panic
		} else if err != nil {
			tx.Rollback()
		}
	}()
	query := `
    CREATE TABLE IF NOT EXISTS document (
        id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        title text NOT NULL,
        content text NOT NULL,
        published BOOL NOT NULL DEFAULT false,
        created timestamp with time zone DEFAULT now()
    );
    CREATE TABLE IF NOT EXISTS tag (
        id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        value text UNIQUE NOT NULL
    );
    CREATE TABLE IF NOT EXISTS document_tag (
        document bigint NOT NULL REFERENCES document (id) ON DELETE CASCADE,
        tag bigint NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
        PRIMARY KEY (document, tag)
    );
    `
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	return nil
}

type basicMds struct {
	db *sql.DB
}

func initializeMarkdownStore() MarkdownStore {
	db := createDbConnection()
	initializeModel(db)
	return &basicMds{db: db}
}

func (ms *basicMds) Close() {
	ms.db.Close()
}

func (ms *basicMds) GetDocumentById(id int) (Document, error) {
	docs, err := ms.queryDocuments("d.id = " + strconv.Itoa(id))
	if err != nil {
		return nil, err
	} else if len(docs) == 0 {
		return nil, fmt.Errorf("Document %d does not exist", id)
	} else {
		return docs[0], nil
	}
}

func (ms *basicMds) GetDocumentsByTag(tag string) ([]Document, error) {
	return ms.queryDocuments("published is true and t.value = '" + tag + "'")
}

func (ms *basicMds) Documents() ([]Document, error) {
	return ms.queryDocuments("published is true")
}

func (ms *basicMds) Drafts() ([]Document, error) {
	return ms.queryDocuments("published is false")
}

func (ms *basicMds) queryDocuments(filter string) ([]Document, error) {
	documents := make([]Document, 0)
	query := `
    SELECT 
    d.id AS document_id, 
    d.title, 
    d.published,
    d.created, 
    array_agg(t.value) AS tags
    FROM 
        document d
    LEFT JOIN 
        document_tag dt ON d.id = dt.document
    LEFT JOIN 
        tag t ON dt.tag = t.id
    `
	if filter != "" {
		query += `
        WHERE ` + filter
	}
	query += `
    GROUP BY d.id, d.title, d.created
    ORDER BY d.created DESC;
    `
	rows, err := ms.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to query documents: %w", err)
	}
	for rows.Next() {
		var id int
		var title string
		var published bool
		var created time.Time
		var tags []sql.NullString
		err := rows.Scan(&id, &title, &published, &created, pq.Array(&tags))
		if err != nil {
			panic(err)
		}
		docTags := make([]string, 0, len(tags))
		for _, tag := range tags {
			if tag.Valid {
				docTags = append(docTags, tag.String)
			}
		}
		documents = append(documents, &basicDocument{
			db:        ms.db,
			id:        id,
			title:     title,
			created:   created,
			tags:      docTags,
			published: published,
			loaded:    false,
		})
	}
	return documents, nil
}

type basicDocument struct {
	db        *sql.DB
	id        int
	title     string
	content   string
	published bool
	tags      []string
	created   time.Time
	loaded    bool
}

func (ms *basicMds) CreateDocument(title string, content string, tags []string) (Document, error) {
	tx, err := ms.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw the panic
		} else if err != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("Failed to begin transaction: %w", err)
	}
	// Create the document
	query := `
    INSERT INTO document (title, content) VALUES ($1, $2) 
    RETURNING id, created;
    `
	row := ms.db.QueryRow(query, title, content)
	var docId int
	var created time.Time
	err = row.Scan(&docId, &created)
	if err != nil {
		return nil, fmt.Errorf("Failed to create document: %w", err)
	}
	// Create the tags if they are missing
	tagIds := make([]int, len(tags))
	query = `
    INSERT INTO tag (value) 
    VALUES ($1) 
    ON CONFLICT (value) DO UPDATE 
    SET value = tag.value
    RETURNING id;
    `
	for i, tag := range tags {
		row := ms.db.QueryRow(query, tag)
		err = row.Scan(&tagIds[i])
		if err != nil {
			return nil, fmt.Errorf("Failed to create tag: %w", err)
		}
	}
	// Insert document - tag associations
	query = `
    INSERT INTO document_tag (document, tag) 
    VALUES ($1, $2); 
    `
	for _, tagId := range tagIds {
		_, err := ms.db.Exec(query, docId, tagId)
		if err != nil {
			return nil, fmt.Errorf("Failed to assign tag to doc: %w", err)
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	// Return the created document
	return &basicDocument{
		db:        ms.db,
		id:        docId,
		title:     title,
		content:   content,
		tags:      tags,
		published: false,
		created:   created,
		loaded:    true,
	}, nil
}

func (d *basicDocument) Id() int {
	return d.id
}

func (d *basicDocument) Tags() []string {
	return d.tags
}

func (d *basicDocument) SetTags(tags []string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw the panic
		} else if err != nil {
			tx.Rollback()
		}
	}()
	// Create the tags if they are missing
	tagIds := make([]int, len(tags))
	query := `
    INSERT INTO tag (value) 
    VALUES ($1) 
    ON CONFLICT (value) DO UPDATE 
    SET value = tag.value
    RETURNING id;
    `
	for i, tag := range tags {
		row := tx.QueryRow(query, tag)
		err = row.Scan(&tagIds[i])
		if err != nil {
			return fmt.Errorf("Failed to create tag: %w", err)
		}
	}
	// Reset document - tag associations
	query = fmt.Sprintf("DELETE FROM document_tag WHERE document = %d;", d.id)
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}
	query = `
    INSERT INTO document_tag (document, tag) 
    VALUES ($1, $2); 
    `
	for _, tagId := range tagIds {
		_, err := tx.Exec(query, d.id, tagId)
		if err != nil {
			return fmt.Errorf("Failed to assign tag to doc: %w", err)
		}
	}
	if err == nil {
		tx.Commit()
	}
	return err
}

func (d *basicDocument) Title() string {
	return d.title
}

func (d *basicDocument) SetTitle(title string) error {
	err := d.update(fmt.Sprintf("title = '%s'", title))
	if err != nil {
		d.title = title
	}
	return err
}

func (d *basicDocument) Content() string {
	if d.loaded {
		return d.content
	} else {
		err := d.loadContent()
		if err != nil {
			return fmt.Sprint(err)
		}
		return d.content
	}
}

func (d *basicDocument) loadContent() error {
	query := fmt.Sprintf("SELECT content FROM document WHERE id = %d", d.Id())
	row := d.db.QueryRow(query)
	var content string
	err := row.Scan(&content)
	if err != nil {
		return err
	}
	d.content = content
	d.loaded = true
	return nil
}

func (d *basicDocument) SetContent(content string) error {
	err := d.update(fmt.Sprintf("content = '%s'", content))
	if err != nil {
		d.content = content
	}
	return err
}

func (d *basicDocument) Created() time.Time {
	return d.created
}

func (d *basicDocument) Published() bool {
	return d.published
}

func (d *basicDocument) SetPublished(published bool) error {
	err := d.update(fmt.Sprintf("published = %t", published))
	if err != nil {
		d.published = published
	}
	return err
}

func (d *basicDocument) Delete() error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw the panic
		} else if err != nil {
			tx.Rollback()
		}
	}()
	query := fmt.Sprintf("DELETE FROM document WHERE id = %d;", d.id)
	_, err = tx.Exec(query)
	if err == nil {
		tx.Commit()
	}
	return err
}

func (d *basicDocument) update(set string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw the panic
		} else if err != nil {
			tx.Rollback()
		}
	}()
	query := fmt.Sprintf("UPDATE document SET %s WHERE id = %d;", set, d.id)
	_, err = tx.Exec(query)
	if err == nil {
		tx.Commit()
	}
	return err
}

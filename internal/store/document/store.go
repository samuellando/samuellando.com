package document

import (
	"database/sql"
	"fmt"

	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/store"

	"github.com/lib/pq"
)

type Store struct {
	db  *sql.DB
	run func() ([]*Document, error)
}

func CreateStore(db *sql.DB) Store {
	return Store{db: db, run: func() ([]*Document, error) {
		return queryDocuments(db, "")
	}}
}

func createErrorStore(err error) *Store {
	return &Store{db: nil, run: func() ([]*Document, error) {
		return nil, err
	}}
}

func (ds *Store) New(data []*Document) store.Store[*Document] {
	return &Store{db: ds.db, run: func() ([]*Document, error) {
        return data, nil
	}}
}

func (ds *Store) GetById(id int) (*Document, error) {
	docs, err := ds.run()
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		if doc.Id() == id {
			return doc, nil
		}
	}
	return nil, fmt.Errorf("Document %d does not exist", id)
}

func (ds *Store) GetAll() ([]*Document, error) {
	return ds.run()
}

func (ds *Store) Add(d *Document) error {
	tx, err := ds.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	// Create the document
	query := `
    INSERT INTO document (title, content, created) VALUES ($1, $2, $3) 
    RETURNING id, created;
    `
	content, err := d.Content()
	if err != nil {
		return err
	}
	row := tx.QueryRow(query, d.Title(), content, d.Created())
	err = row.Scan(&d.id, &d.fields.Created)
	if err != nil {
		return fmt.Errorf("Failed to create document: %w", err)
	}
	d.setTags(tx)
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	d.db = ds.db
	return nil
}

func (ds *Store) Remove(d *Document) error {
	if d.db == nil {
		return fmt.Errorf("Cannot delete a proto document")
	}
	docs, err := ds.GetAll()
	if err != nil {
		return err
	}
	exists := false
	for _, item := range docs {
		if item.Id() == d.Id() {
			exists = true
		}
	}
	if !exists {
		return fmt.Errorf("This document is not in this store")
	}
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := "DELETE FROM document WHERE id = $1;"
	_, err = tx.Exec(query, d.id)
	if err != nil {
		return fmt.Errorf("Failed to delete document: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	// Make it a proto
	d.id = -1
	d.db = nil
	return nil
}

func (ds *Store) Filter(f func(*Document) bool) store.Store[*Document] {
    n, err := store.Filter(ds, f)
    if err != nil {
        return createErrorStore(err)
    }
    return n
}

func (ds *Store) Group(f func(*Document) string) *datatypes.OrderedMap[string, store.Store[*Document]] {
    m, err := store.Group(ds, f)
    if err != nil {
        m := datatypes.NewOrderedMap[string, store.Store[*Document]]()
        m.Set("", createErrorStore(err))
        return m
    }
    return m
}

func (ds *Store) Sort(f func(*Document, *Document) bool) store.Store[*Document] {
    n, err := store.Sort(ds, f)
    if err != nil {
        return createErrorStore(err)
    }
    return n
}

func queryDocuments(db *sql.DB, filter string, args ...any) ([]*Document, error) {
	query := `
    SELECT 
    d.id AS document_id, 
    d.title, 
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
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to query documents: %w", err)
	}
	documents := make([]*Document, 0)
	for rows.Next() {
		var id int
		docFields := DocumentFeilds{}
		var tags []sql.NullString
		err := rows.Scan(&id, &docFields.Title, &docFields.Created, pq.Array(&tags))
		if err != nil {
			panic(err)
		}
		docTags := make([]string, 0, len(tags))
		for _, tag := range tags {
			if tag.Valid {
				docTags = append(docTags, tag.String)
			}
		}
		docFields.Tags = docTags
		documents = append(documents, &Document{
			db:     db,
			loaded: false,
			id:     id,
			fields: docFields,
		})
	}
	return documents, nil
}

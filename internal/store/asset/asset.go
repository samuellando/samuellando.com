package asset

import (
	"database/sql"
	"fmt"
	"time"
)

type Asset struct {
	db      *sql.DB
	id      int
	name    string
	created time.Time
	content []byte
	loaded  bool
}

type AssetFields struct {
	Name    string
	Content []byte
}

func CreateProto(opts ...func(*AssetFields)) Asset {
	fields := AssetFields{
		Content: []byte{},
	}
	for _, opt := range opts {
		opt(&fields)
	}
    content := make([]byte, len(fields.Content))
    copy(content, fields.Content)
	return Asset{name: fields.Name, content: content, loaded: true}
}

func (a *Asset) Id() int {
	return a.id
}

func (a *Asset) Name() string {
	return a.name
}

func (a *Asset) Created() time.Time {
	return a.created
}

func (a *Asset) Content() ([]byte, error) {
	query := `
        SELECT content
        FROM asset
        WHERE id = $1
    `
	if !a.loaded {
		row := a.db.QueryRow(query, a.Id())
		err := row.Scan(&a.content)
		if err != nil {
			return nil, err
		}
		a.loaded = true
	}
	return a.content, nil
}

func (a *Asset) Delete() error {
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
        DELETE FROM asset
        WHERE id = $1
    `
	_, err = tx.Exec(query, a.Id())
	if err != nil {
		return err
	}
    err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
    return nil
}

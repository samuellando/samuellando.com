package tag

import (
	"database/sql"
	"fmt"
	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/store"
)

type Store struct {
	db  *sql.DB
	run func() ([]Tag, error)
}

func CreateStore(db *sql.DB) Store {
	return Store{db: db, run: func() ([]Tag, error) {
		return loadTags(db)
	}}
}

func createErrorStore(err error) *Store {
	return &Store{run: func() ([]Tag, error) {
		return nil, err
	}}
}

func (as *Store) Add(a Tag) error {
	tx, err := as.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	// Create the asset
	query := `
    INSERT INTO tag (value, color) VALUES ($1, $2) 
    RETURNING id;
    `
	row := tx.QueryRow(query, a.Value(), a.Color())
	err = row.Scan(&a.id)
	if err != nil {
		return fmt.Errorf("Failed to create document: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	a.db = as.db
	return nil
}

func (as *Store) GetById(id int) (Tag, error) {
	assets, err := as.run()
	if err != nil {
		return Tag{}, err
	}
	for _, a := range assets {
		if a.Id() == id {
			return a, nil
		}
	}
	return Tag{}, fmt.Errorf("Tag with id: %d does not exist", id)
}

func (as *Store) GetByValue(value string) (Tag, error) {
	assets, err := as.run()
	if err != nil {
		return Tag{}, err
	}
	for _, a := range assets {
		if a.Value() == value {
			return a, nil
		}
	}
	return Tag{}, fmt.Errorf("Tag with value: %s does not exist", value)
}

func (as *Store) GetAll() ([]Tag, error) {
	return as.run()
}

func (as *Store) Filter(f func(Tag) bool) store.Store[Tag] {
	n, err := store.Filter(as, f)
	if err != nil {
		return createErrorStore(err)
	}
	return n
}

func (as *Store) Group(f func(Tag) string) *datatypes.OrderedMap[string, store.Store[Tag]] {
	n, err := store.Group(as, f)
	if err != nil {
		m := datatypes.NewOrderedMap[string, store.Store[Tag]]()
		m.Set("", createErrorStore(err))
		return m
	}
	return n
}

func (as *Store) Sort(f func(Tag, Tag) bool) store.Store[Tag] {
	n, err := store.Sort(as, f)
	if err != nil {
		return createErrorStore(err)
	}
	return n
}

func (as *Store) New(d []Tag) store.Store[Tag] {
	return &Store{db: as.db, run: func() ([]Tag, error) {
		return d, nil
	}}
}

func loadTags(db *sql.DB) ([]Tag, error) {
	rows, err := db.Query(`
    SELECT 
        id, value, color
    FROM
        tag
    `)
	if err != nil {
		return nil, err
	}
	tags := make([]Tag, 0)
	for rows.Next() {
		tag := Tag{db: db}
		rows.Scan(&tag.id, &tag.value, &tag.color)
		tags = append(tags, tag)
	}
	return tags, nil
}

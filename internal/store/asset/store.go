package asset

import (
	"database/sql"
	"fmt"
    "samuellando.com/internal/store"
    "samuellando.com/internal/datatypes"
)

type Store struct {
	db  *sql.DB
	run func() ([]*Asset, error)
}

func CreateStore(db *sql.DB) Store {
    return Store{db: db, run: func() ([]*Asset, error) {
        return loadAssets(db)
    }}
}

func createErrorStore(err error) *Store {
	return &Store{run: func() ([]*Asset, error) {
		return nil, err
	}}
}

func (as *Store) Add(a *Asset) error {
	tx, err := as.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	// Create the asset
	query := `
    INSERT INTO asset (name, content) VALUES ($1, $2) 
    RETURNING id, created;
    `
	content, err := a.Content()
	if err != nil {
		return err
	}
	row := tx.QueryRow(query, a.Name(), content)
	err = row.Scan(&a.id, &a.created)
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

func (as *Store) GetById(id int) (*Asset, error) {
    assets, err := as.run()
    if err != nil {
        return nil, err
    }
    for _, a := range assets {
        if a.Id() == id {
            return a, nil
        }
    }
    return nil, fmt.Errorf("Asset with id: %d does not exist", id)
}

func (as *Store) GetByName(name string) (*Asset, error) {
    assets, err := as.run()
    if err != nil {
        return nil, err
    }
    for _, a := range assets {
        if a.Name() == name {
            return a, nil
        }
    }
    return nil, fmt.Errorf("Asset with name: %s does not exist", name)
}

func (as *Store) GetAll() ([]*Asset, error) {
	return as.run()
}

func (as *Store) Filter(f func(*Asset) bool) store.Store[*Asset] {
    n, err := store.Filter(as, f)
    if err != nil {
        return createErrorStore(err)
    }
    return n
}

func (as *Store) Group(f func(*Asset) string) *datatypes.OrderedMap[string, store.Store[*Asset]] {
    n, err := store.Group(as, f)
    if err != nil {
        m := datatypes.NewOrderedMap[string, store.Store[*Asset]]()
        m.Set("", createErrorStore(err))
        return m
    }
    return n
}

func (as *Store) Sort(f func(*Asset, *Asset) bool) store.Store[*Asset] {
    n, err := store.Sort(as, f)
    if err != nil {
        return createErrorStore(err)
    }
    return n
}

func (as *Store) New(d []*Asset) store.Store[*Asset] {
    return &Store{db: as.db, run: func() ([]*Asset, error) {
		return d, nil
	}}
}

func loadAssets(db *sql.DB) ([]*Asset, error) {
    rows, err := db.Query(`
    SELECT 
        id, name, created
    FROM
        asset
    `)
    if err != nil {
        return nil, err
    }
    assets := make([]*Asset, 0)
    for rows.Next() {
        asset := Asset{db: db, loaded: false}
        rows.Scan(&asset.id, &asset.name, &asset.created)
        assets = append(assets, &asset)
    }
    return assets, nil
}

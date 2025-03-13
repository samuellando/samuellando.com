package asset

import (
	"context"
	"database/sql"

	"samuellando.com/data"
	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/store"
)

type Store struct {
	db *sql.DB
}

func CreateStore(db *sql.DB) Store {
	return Store{db: db}
}

func (as Store) Add(a ProtoAsset) (Asset, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	row, err := queries.CreateAsset(ctx, data.CreateAssetParams{
		Name:    a.Name,
		Content: a.Content,
	})
	if err != nil {
		return Asset{}, err
	}
	return Asset{
		db:      as.db,
		id:      row.Asset.ID,
		name:    row.Asset.Name,
		created: row.Asset.Created,
		content: row.Asset.Content,
		loaded:  true,
	}, nil
}

func (as Store) GetById(id int64) (Asset, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	row, err := queries.GetAsset(ctx, id)
	if err != nil {
		return Asset{}, err
	}
	return Asset{
		db:      as.db,
		id:      row.Asset.ID,
		name:    row.Asset.Name,
		created: row.Asset.Created,
		content: row.Asset.Content,
		loaded:  true,
	}, nil
}

func (as Store) GetByName(name string) (Asset, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	row, err := queries.GetAssetByName(ctx, name)
	if err != nil {
		return Asset{}, err
	}
	return Asset{
		db:      as.db,
		id:      row.Asset.ID,
		name:    row.Asset.Name,
		created: row.Asset.Created,
		content: row.Asset.Content,
		loaded:  true,
	}, nil
}

func (as Store) GetAll() ([]Asset, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	rows, err := queries.GetAssets(ctx)
	if err != nil {
		return nil, err
	}
	assets := make([]Asset, len(rows))
	for i, row := range rows {
		assets[i] = Asset{
			db:      as.db,
			id:      row.ID,
			name:    row.Name,
			created: row.Created,
			loaded:  false,
		}
	}
	return assets, nil
}

func (as Store) Filter(f func(Asset) bool) (store.Store[Asset], error) {
	return store.Filter(as, f)
}

func (as Store) Group(f func(Asset) string) (datatypes.OrderedMap[string, store.Store[Asset]], error) {
	return store.Group(as, f)
}

func (as Store) Sort(f func(Asset, Asset) bool) (store.Store[Asset], error) {
	return store.Sort(as, f)
}

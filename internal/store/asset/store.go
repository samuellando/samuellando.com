package asset

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"strconv"
	"time"

	"samuellando.com/data"
	"samuellando.com/internal/cache"
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

func encode(o any) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(o); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decode(data []byte, p any) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	if err := dec.Decode(p); err != nil {
		return err
	}
	return nil
}

func (as Store) GetById(id int64) (Asset, error) {
	c := cache.ParamCached(func() ([]byte, error) {
		ctx := context.TODO()
		queries := data.New(as.db)
		row, err := queries.GetAsset(ctx, id)
		if err != nil {
			return nil, err
		}
		b, err := encode(row)
		if err != nil {
			return nil, err
		}
		return b, nil
	}, strconv.Itoa(int(id)), func(co *cache.CacheOptions) {
		co.MaxAge = time.Minute * 5
	})
	d, err := c()
	if err != nil {
		return Asset{}, err
	}
	var row data.GetAssetRow
	err = decode(d, &row)
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
	c := cache.ParamCached(func() ([]byte, error) {
		ctx := context.TODO()
		queries := data.New(as.db)
		row, err := queries.GetAssetByName(ctx, name)
		if err != nil {
			return nil, err
		}
		b, err := encode(row)
		if err != nil {
			return nil, err
		}
		return b, nil
	}, name, func(co *cache.CacheOptions) {
		co.MaxAge = time.Minute * 5
	})
	d, err := c()
	if err != nil {
		return Asset{}, err
	}
	var row data.GetAssetByNameRow
	err = decode(d, &row)
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

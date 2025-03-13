package tag

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

func (as Store) Add(p ProtoTag) (Tag, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	row, err := queries.CreateOrUpdateTag(ctx, data.CreateOrUpdateTagParams{
		Value: p.Value,
		Color: p.Color,
	})
	if err != nil {
		return Tag{}, err
	}
	return Tag{
		db:    as.db,
		id:    row.Tag.ID,
		value: row.Tag.Value,
		color: row.Tag.Color,
	}, nil
}

func (as Store) GetById(id int64) (Tag, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	row, err := queries.GetTag(ctx, id)
	if err != nil {
		return Tag{}, err
	}
	return Tag{
		db:    as.db,
		id:    row.Tag.ID,
		value: row.Tag.Value,
		color: row.Tag.Color,
	}, nil
}

func (as Store) GetByValue(value string) (Tag, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	row, err := queries.GetTagByValue(ctx, value)
	if err != nil {
		return Tag{}, err
	}
	return Tag{
		db:    as.db,
		id:    row.Tag.ID,
		value: row.Tag.Value,
		color: row.Tag.Color,
	}, nil
}

func (as Store) GetAll() ([]Tag, error) {
	ctx := context.TODO()
	queries := data.New(as.db)
	rows, err := queries.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	tags := make([]Tag, len(rows))
	for i, row := range rows {
		tags[i] = Tag{
			id:    row.Tag.ID,
			value: row.Tag.Value,
			color: row.Tag.Color,
		}
	}
	return tags, nil
}

func (as Store) Filter(f func(Tag) bool) (store.Store[Tag], error) {
	return store.Filter(as, f)
}

func (as Store) Group(f func(Tag) string) (datatypes.OrderedMap[string, store.Store[Tag]], error) {
	return store.Group(as, f)
}

func (as Store) Sort(f func(Tag, Tag) bool) (store.Store[Tag], error) {
	return store.Sort(as, f)
}

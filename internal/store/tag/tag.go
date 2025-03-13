package tag

import (
	"context"
	"database/sql"
	"fmt"
	"samuellando.com/data"
)

type Tag struct {
	db    *sql.DB
	id    int64
	value string
	color string
}

type ProtoTag struct {
	Value string
	Color string
}

func (a Tag) Id() int64 {
	return a.id
}

func (a Tag) Value() string {
	return a.value
}

func (a Tag) Color() string {
	return a.color
}

func (a *Tag) Update(opts ...func(*ProtoTag)) error {
	proto := ProtoTag{
		Value: a.Value(),
		Color: a.Color(),
	}
	for _, opt := range opts {
		opt(&proto)
	}
	if proto.Value != a.Value() {
		return fmt.Errorf("Cannot change the value of a tag")
	}

	ctx := context.TODO()
	queries := data.New(a.db)
	row, err := queries.CreateOrUpdateTag(ctx, data.CreateOrUpdateTagParams{
		Value: a.Value(),
		Color: proto.Color,
	})
	if err != nil {
		return err
	}
	a.value = row.Tag.Value
	a.color = row.Tag.Color
	return nil
}

func (a *Tag) Delete() error {
	ctx := context.TODO()
	queries := data.New(a.db)
	err := queries.DeleteTag(ctx, a.Id())
	return err
}

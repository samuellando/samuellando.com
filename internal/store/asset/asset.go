package asset

import (
	"context"
	"database/sql"
	"time"

	"samuellando.com/data"
)

type Asset struct {
	db      *sql.DB
	id      int64
	name    string
	created time.Time
	content []byte
	loaded  bool
}

type ProtoAsset struct {
	Name    string
	Content []byte
}

func (a Asset) Id() int64 {
	return a.id
}

func (a Asset) Name() string {
	return a.name
}

func (a Asset) Created() time.Time {
	return a.created
}

func (a *Asset) Content() ([]byte, error) {
	if a.loaded {
		return a.content, nil
	}
	ctx := context.TODO()
	queries := data.New(a.db)
	content, err := queries.GetAssetContent(ctx, a.id)
	if err != nil {
		return nil, err
	}
	a.content = content
	a.loaded = true
	return content, nil
}

func (a *Asset) Delete() error {
	ctx := context.TODO()
	queries := data.New(a.db)
	err := queries.DeleteAsset(ctx, a.id)
	return err
}

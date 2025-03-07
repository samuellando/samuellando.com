package document

import (
	"database/sql"
	"fmt"

	"context"

	"samuellando.com/data"
	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/store"
	"samuellando.com/internal/store/tag"
)

type Store struct {
	db           *sql.DB
	materialized *store.MaterializedStore[Document]
}

func CreateStore(db *sql.DB) Store {
	return Store{db: db, materialized: nil}
}

func (ds Store) GetById(id int64) (Document, error) {
	if ds.materialized != nil {
		return ds.materialized.GetById(id)
	}
	ctx := context.TODO()
	queires := data.New(ds.db)
	rows, err := queires.GetDocument(ctx, id)
	if err != nil {
		return Document{}, err
	}
	if len(rows) == 0 {
		return Document{}, fmt.Errorf("Document not found")
	}
	tags := make([]tag.ProtoTag, 0)
	for _, row := range rows {
		if row.TagID.Valid {
			tags = append(tags, tag.ProtoTag{
				Value: row.TagValue.String,
				Color: row.TagColor.String,
			})
		}
	}
	return Document{
		db:      ds.db,
		id:      rows[0].Document.ID,
		title:   rows[0].Document.Title,
		content: rows[0].Document.Content,
		created: rows[0].Document.Created,
		tags:    tags,
	}, nil
}

func (ds Store) GetAll() ([]Document, error) {
	if ds.materialized != nil {
		return ds.materialized.GetAll()
	}
	ctx := context.TODO()
	queires := data.New(ds.db)
	docRows, err := queires.GetDocuments(ctx)
	if err != nil {
		return nil, err
	}
	docs := make(map[int64]*Document)
	for _, row := range docRows {
		if _, ok := docs[row.Document.ID]; !ok {
			docs[row.Document.ID] = &Document{
				db:      ds.db,
				id:      row.Document.ID,
				title:   row.Document.Title,
				content: row.Document.Content,
				created: row.Document.Created,
				tags:    make([]tag.ProtoTag, 0),
			}
		}
		if row.TagID.Valid {
			tag := tag.ProtoTag{
				Value: row.TagValue.String,
				Color: row.TagColor.String,
			}
			docs[row.Document.ID].tags = append(docs[row.Document.ID].tags, tag)
		}
	}
	res := make([]Document, 0)
	for _, doc := range docs {
		res = append(res, *doc)
	}
	return res, nil
}

func (ds Store) Add(p ProtoDocument) (Document, error) {
	ctx := context.TODO()
	tx, err := ds.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return Document{}, err
	}
	queries := data.New(ds.db).WithTx(tx)
	id, err := queries.CreateDocument(ctx, data.CreateDocumentParams{
		Title:   p.Title,
		Content: p.Content,
		Created: p.Created,
	})
	if err != nil {
		return Document{}, err
	}
	tagRows, err := queries.SetDocumentTags(ctx, data.SetDocumentTagsParams{
		Document:  id,
		TagValues: tagValues(p.Tags),
	})
	if err != nil {
		return Document{}, err
	}
	err = tx.Commit()
	if err != nil {
		return Document{}, err
	}
	tags := make([]tag.ProtoTag, len(tagRows))
	for i, tagRow := range tagRows {
		tags[i] = tag.ProtoTag{
			Value: tagRow.Value,
			Color: tagRow.Color,
		}
	}
	return Document{
		db:      ds.db,
		id:      id,
		title:   p.Title,
		content: p.Content,
		created: p.Created,
		tags:    tags,
	}, nil
}

func (ds Store) Filter(f func(Document) bool) (store.Store[Document], error) {
	var filtered store.Store[Document]
	var err error
	if ds.materialized != nil {
		filtered, err = ds.materialized.Filter(f)
	} else {
		filtered, err = store.Filter(ds, f)
	}
	if err != nil {
		return ds, err
	}
	if ms, ok := filtered.(store.MaterializedStore[Document]); ok {
		return Store{db: ds.db, materialized: &ms}, nil
	} else {
		panic("Could not type cast to MaterializedStore!")
	}
}

func (ds Store) Group(f func(Document) string) (datatypes.OrderedMap[string, store.Store[Document]], error) {
	return store.Group(ds, f)
}

func (ds Store) Sort(f func(Document, Document) bool) (store.Store[Document], error) {
	var sorted store.Store[Document]
	var err error
	if ds.materialized != nil {
		sorted, err = ds.materialized.Sort(f)
	} else {
		sorted, err = store.Sort(ds, f)
	}
	if err != nil {
		return ds, err
	}
	if ms, ok := sorted.(store.MaterializedStore[Document]); ok {
		return Store{db: ds.db, materialized: &ms}, nil
	} else {
		panic("Could not type cast to MaterializedStore!")
	}
}

func (ds Store) AllTags() ([]tag.ProtoTag, error) {
	ctx := context.TODO()
	queries := data.New(ds.db)
	tagRows, err := queries.GetAllDocumentTags(ctx)
	if err != nil {
		return nil, err
	}
	tags := make([]tag.ProtoTag, len(tagRows))
	for i, tagRow := range tagRows {
		tags[i] = tag.ProtoTag{
			Value: tagRow.Tag.Value,
			Color: tagRow.Tag.Color,
		}
	}
	return tags, nil
}

func (ds Store) AllSharedTags(tagValue string) ([]tag.ProtoTag, error) {
	ctx := context.TODO()
	queries := data.New(ds.db)
	tagRows, err := queries.GetSharedDocumentTags(ctx, tagValue)
	if err != nil {
		return nil, err
	}
	tags := make([]tag.ProtoTag, len(tagRows))
	for i, tagRow := range tagRows {
		tags[i] = tag.ProtoTag{
			Value: tagRow.Tag.Value,
			Color: tagRow.Tag.Color,
		}
	}
	return tags, nil
}

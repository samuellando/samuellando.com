package project

import (
	"context"
	"database/sql"
	"fmt"
	"samuellando.com/data"
	"samuellando.com/internal/store/tag"
	"strings"
	"time"
)

type Project struct {
	db          *sql.DB
	id          int64
	name        string
	created     time.Time
	pushed      time.Time
	url         string
	description *string
	tags        []tag.ProtoTag
}

type ProtoProject struct {
	Description *string
	Tags        []tag.ProtoTag
}

func (p Project) Id() int64 {
	return p.id
}

func (p Project) Title() string {
	return p.name
}

func (p Project) Description() string {
	if p.description == nil {
		return ""
	}
	return *p.description
}

func (p Project) Created() time.Time {
	return p.created
}

func (p Project) Pushed() time.Time {
	return p.pushed
}

func (p Project) Url() string {
	return p.url
}

func (p Project) Tags() []tag.ProtoTag {
	return copyOf(p.tags)
}

func (p *Project) Update(setters ...func(*ProtoProject)) error {
	desc := p.Description()
	proto := ProtoProject{
		Description: &desc,
		Tags:        p.Tags(),
	}
	for _, setter := range setters {
		setter(&proto)
	}

	ctx := context.TODO()
	tx, err := p.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	queries := data.New(p.db).WithTx(tx)
	sqldesc := sql.NullString{Valid: false}
	if proto.Description != nil {
		sqldesc = sql.NullString{Valid: true, String: *proto.Description}
	}
	err = queries.UpdateProject(ctx, data.UpdateProjectParams{
		ID:          p.id,
		Description: sqldesc,
	})
	if err != nil {
		return err
	}
	tagRows, err := queries.SetProjectTags(ctx, data.SetProjectTagsParams{
		Project:   p.id,
		TagValues: tagValues(proto.Tags),
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	tags := make([]tag.ProtoTag, len(tagRows))
	for i, tagRow := range tagRows {
		tags[i] = tag.ProtoTag{
			Value: tagRow.Value,
			Color: tagRow.Color,
		}
	}
	p.description = proto.Description
	p.tags = tags
	return nil
}

func (p Project) ToString() string {
	s := fmt.Sprintf("%s\n%s\n%s", p.Title(), p.Description(), strings.Join(tagValues(p.Tags()), " "))
	return s
}

func copyOf(src []tag.ProtoTag) []tag.ProtoTag {
	tagsCopy := make([]tag.ProtoTag, len(src))
	copy(tagsCopy, src)
	return tagsCopy
}

func tagValues(src []tag.ProtoTag) []string {
	s := make([]string, len(src))
	for i, tag := range src {
		s[i] = tag.Value
	}
	return s
}

package tag

import (
	"database/sql"
	"fmt"
)

type Tag struct {
	db    *sql.DB
	id    int
	value string
	color *string
}

type TagFields struct {
	Value string
	Color string
}

func CreateProto(opts ...func(*TagFields)) Tag {
	fields := TagFields{
		Value: "",
		Color: "",
	}
	for _, opt := range opts {
		opt(&fields)
	}
	return Tag{value: fields.Value, color: &fields.Color}
}

func (a Tag) Id() int {
	return a.id
}

func (a Tag) Value() string {
	return a.value
}

func (a Tag) Color() *string {
	if a.color == nil || *a.color == "" {
		return nil
	}
	return a.color
}

func (a Tag) Update(opts ...func(*TagFields)) error {
	color := ""
	if a.color != nil {
		color = *a.color
	}
	fields := TagFields{
		Value: a.value,
		Color: color,
	}
	for _, opt := range opts {
		opt(&fields)
	}
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
        UPDATE tag
        SET color = $1
        WHERE id = $2
    `
	_, err = tx.Exec(query, fields.Color, a.Id())
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	a.color = &fields.Color
	return nil
}

func (a Tag) Delete() error {
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
        DELETE FROM tag
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

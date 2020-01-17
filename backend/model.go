package main

import (
	"database/sql"
)

type page struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (p *page) getPage(db *sql.DB) error {
	return db.QueryRow("SELECT title, text FROM pages WHERE id=$1",
		p.ID).Scan(&p.Title, &p.Text)
}

func (p *page) updatePage(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE pages SET title=$1, text=$2 WHERE id=$3",
			p.Title, p.Text, p.ID)

	return err
}

func (p *page) deletePage(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM pages WHERE id=$1", p.ID)

	return err
}

func (p *page) createPage(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO pages(title, text) VALUES($1, $2) RETURNING id",
		p.Title, p.Text).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getPages(db *sql.DB, start, count int) ([]page, error) {
	rows, err := db.Query(
		"SELECT id, title,  text FROM pages LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pages := []page{}

	for rows.Next() {
		var p page
		if err := rows.Scan(&p.ID, &p.Title, &p.Text); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}

	return pages, nil
}

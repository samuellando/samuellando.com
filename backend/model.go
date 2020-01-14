
package main


import (
  "database/sql"
  "errors"
)

type page struct {
  ID    int     `json:"id"`
  Title string  `json:"title"`
  Text  string  `json:"text"`
}

func (p *page) getPage(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (p *page) updatePage(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (p *page) deletePage(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (p *page) createPage(db *sql.DB) error {
  return errors.New("Not implemented")
}

func getPages(db *sql.DB, start, count int) ([]page, error) {
  return nil, errors.New("Not implemented")
}

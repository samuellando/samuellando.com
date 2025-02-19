package search

import (
	"fmt"
	"samuellando.com/internal/store"
)

// Interface for stored items that are Searchable
type Searchable interface {
	Id() int
	Title() string
	ToString() string
}

type IndexItem struct {
	Path string
	Type string
	Text string
	Item Searchable
}

func GenerateIndex[T Searchable](typ, basePath string, s store.Store[T]) func() []IndexItem {
	return func() []IndexItem {
		all, err := s.GetAll()
		if err != nil {
			return []IndexItem{}
		}
		items := make([]IndexItem, 0, len(all))
		for _, elem := range all {
			path := fmt.Sprintf("%s/%d", basePath, elem.Id())
			items = append(items, IndexItem{Type: typ, Path: path, Text: elem.ToString(), Item: elem})
		}
		return items
	}
}

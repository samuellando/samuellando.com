package store

import (
	"sort"
	"strings"

	"samuellando.com/internal/datatypes"
)

// A egneral interface for anything that stores anything.
//
// These support adding and removeing items, and should handle those operations
// To the underlying data structures.
//
// By implementing this interface you get access to a the helper methods below
// Wich make implementing Filter, Sort and Group a lot easier.
//
// Notes:
// Clone() Is assumed to return a deep copy of the structure. The elements can
// be shared, but the store should be new.
//
// Reset() is assumed to have no side effects on the original underlying data
// Ie, a new store is created an filled with the new data.
type Store[T any] interface {
	GetById(int) (T, error)
	GetAll() ([]T, error)
	Filter(func(T) bool) Store[T]
	Group(func(T) string) *datatypes.OrderedMap[string, Store[T]]
	Sort(func(T, T) bool) Store[T]
	New([]T) Store[T]
}

func Filter[T any](s Store[T], f func(T) bool) (Store[T], error) {
	filtered := make([]T, 0)
	all, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, elem := range all {
		if f(elem) {
			filtered = append(filtered, elem)
		}
	}
	return s.New(filtered), nil
}

func Group[T any](s Store[T], f func(T) string) (*datatypes.OrderedMap[string, Store[T]], error) {
	all, err := s.GetAll()
	var zero datatypes.OrderedMap[string, Store[T]]
	if err != nil {
		return &zero, err
	}
	// Classify the elements
	groups := datatypes.NewOrderedMap[string, []T]()
	groupNames := []string{}
	for _, elem := range all {
		group := f(elem)
		if v, ok := groups.Get(group); ok {
			groups.Set(group, append(v, elem))
		} else {
			groups.Set(group, []T{elem})
			groupNames = append(groupNames, group)
		}
	}
	// Sorting the group names leads to better UX
	sort.Slice(groupNames, func(i, j int) bool {
		return strings.Compare(groupNames[i], groupNames[j]) > 0
	})
	// Generate the new stores
	stores := datatypes.NewOrderedMap[string, Store[T]]()
	for _, group := range groupNames {
		a, _ := groups.Get(group)
		stores.Set(group, s.New(a))
	}
	return stores, nil
}

type by[T any] struct {
	elems    []T
	lessFunc func(T, T) bool
}

func (a *by[T]) Len() int           { return len(a.elems) }
func (a *by[T]) Swap(i, j int)      { a.elems[i], a.elems[j] = a.elems[j], a.elems[i] }
func (a *by[T]) Less(i, j int) bool { return a.lessFunc(a.elems[i], a.elems[j]) }

func Sort[T any](s Store[T], less func(T, T) bool) (Store[T], error) {
	all, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	b := by[T]{elems: all, lessFunc: less}
	sort.Sort(&b)
	return s.New(b.elems), nil
}

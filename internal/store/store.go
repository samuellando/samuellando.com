package store

import (
	"fmt"
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
type Store[T Indexable] interface {
	GetById(int64) (T, error)
	GetAll() ([]T, error)
	Filter(func(T) bool) (Store[T], error)
	Group(func(T) string) (datatypes.OrderedMap[string, Store[T]], error)
	Sort(func(T, T) bool) (Store[T], error)
}

type Indexable interface {
	Id() int64
}

type MaterializedStore[T Indexable] struct {
	data []T
}

func create[T Indexable](arr []T) Store[T] {
	return MaterializedStore[T]{data: arr}
}

func (ms MaterializedStore[T]) GetById(id int64) (T, error) {
	for _, v := range ms.data {
		if v.Id() == id {
			return v, nil
		}
	}
	var zero T
	return zero, fmt.Errorf("No value with id %d found in store", id)
}

func (ms MaterializedStore[T]) GetAll() ([]T, error) {
	c := make([]T, len(ms.data))
	copy(c, ms.data)
	return c, nil
}

func (ms MaterializedStore[T]) Filter(f func(T) bool) (Store[T], error) {
	return Filter(ms, f)
}

func (ms MaterializedStore[T]) Group(f func(T) string) (datatypes.OrderedMap[string, Store[T]], error) {
	return Group(ms, f)
}

func (ms MaterializedStore[T]) Sort(f func(T, T) bool) (Store[T], error) {
	return Sort(ms, f)
}

func Filter[T Indexable](s Store[T], f func(T) bool) (Store[T], error) {
	filtered := make([]T, 0)
	all, err := s.GetAll()
	if err != nil {
		return MaterializedStore[T]{}, err
	}
	for _, elem := range all {
		if f(elem) {
			filtered = append(filtered, elem)
		}
	}
	return create(filtered), nil
}

func Group[T Indexable](s Store[T], f func(T) string) (datatypes.OrderedMap[string, Store[T]], error) {
	all, err := s.GetAll()
	var zero datatypes.OrderedMap[string, Store[T]]
	if err != nil {
		return zero, err
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
		stores.Set(group, create(a))
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

func Sort[T Indexable](s Store[T], less func(T, T) bool) (Store[T], error) {
	all, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	b := by[T]{elems: all, lessFunc: less}
	sort.Sort(&b)
	return create(b.elems), nil
}

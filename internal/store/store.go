package store

import(
    "sort"
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
    Add(T) error
    Remove(T) error
    Filter(func(T) bool) Store[T]
    Group(func(T) string) map[string]Store[T]
    Sort(func(T, T) bool) Store[T]
}

func Filter[T any](all []T, f func(T) bool) []T {
    filtered := make([]T, 0)
    for _, elem := range all {
        if f(elem) {
            filtered = append(filtered, elem)
        }
    }
    return filtered
}

func Group[T any](all []T, f func(T) string) map[string][]T {
    groups := make(map[string][]T)
    for _, elem := range all {
        group := f(elem)
        if _, ok := groups[group]; ok {
            groups[group] = append(groups[group], elem)
        } else {
            groups[group] = []T{elem}
        }
    }
    return groups
}

type by[T any] struct {
    elems []T
    lessFunc func(T, T) bool
}

func (a *by[T]) Len() int           { return len(a.elems) }
func (a *by[T]) Swap(i, j int)      { a.elems[i], a.elems[j] = a.elems[j], a.elems[i] }
func (a *by[T]) Less(i, j int) bool { return a.lessFunc(a.elems[i], a.elems[j]) }

func Sort[T any](all []T, less func(T, T) bool) []T {
    b := by[T]{elems: all, lessFunc: less}
    sort.Sort(&b)
    return b.elems
}

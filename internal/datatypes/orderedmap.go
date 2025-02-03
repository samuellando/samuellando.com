package datatypes

import "iter"

type OrderedMap[K comparable, V any] struct {
    keys []K
    values []V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
    keys := make([]K, 0)
    values := make([]V, 0)
    return &OrderedMap[K, V]{keys: keys, values: values}
}

func (o *OrderedMap[K, V]) Set(k K, v V) {
    for i, key := range o.keys {
        if key == k {
            o.values[i] = v
            return
        }
    }
    o.keys = append(o.keys, k)
    o.values = append(o.values, v)
}

func (o *OrderedMap[K, V]) Get(k K) (V, bool) {
    for i, key := range o.keys {
        if key == k {
            return o.values[i], true
        }
    }
    var zero V
    return zero, false
}

func (o *OrderedMap[K, V]) All() iter.Seq2[K, V] {
    return func(yield func(K, V) bool) {
        for i, k := range o.keys {
            if !yield(k, o.values[i]) {
                return
            }
        }
    }
}


type elem[K comparable, V any] struct {
    Key K
    Value V
}
// Convert to a slice for use in html templates, since they don't support iter.Seq
func (o *OrderedMap[K, V]) ToSlice() []elem[K, V] {
    elems := make([]elem[K, V], 0, o.Len())
    for k, v := range o.All() {
        elems = append(elems, elem[K, V]{Key: k, Value: v })
    }
    return elems
}

func (o *OrderedMap[K, V]) Len() int {
    return len(o.keys)
}

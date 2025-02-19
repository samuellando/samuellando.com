package datatypes

import "testing"

func TestSetGet(t *testing.T) {
	m := NewOrderedMap[int, int]()
	m.Set(5, 7)
	m.Set(6, 8)
	m.Set(9, 10)
	if v, ok := m.Get(5); !ok || v != 7 {
		t.Fail()
	}
	if v, ok := m.Get(6); !ok || v != 8 {
		t.Fail()
	}
	if v, ok := m.Get(9); !ok || v != 10 {
		t.Fail()
	}
	if _, ok := m.Get(10); ok {
		t.Fail()
	}
}

func TestMultiSet(t *testing.T) {
	m := NewOrderedMap[int, int]()
	m.Set(5, 7)
	m.Set(5, 8)
	m.Set(5, 10)
	if v, ok := m.Get(5); !ok || v != 10 {
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {
	m := NewOrderedMap[int, int]()
	m.Set(5, 7)
	m.Set(6, 8)
	m.Set(9, 10)
	keyOrder := []int{5, 6, 9}
	order := []int{7, 8, 10}
	i := 0
	for k, v := range m.All() {
		if keyOrder[i] != k || order[i] != v {
			t.Fatalf("%d != %d || %d != %d", keyOrder[i], k, order[i], v)
		}
		i++
	}
}

func TestToSlice(t *testing.T) {
	m := NewOrderedMap[int, int]()
	m.Set(5, 7)
	m.Set(6, 8)
	m.Set(9, 10)
	keyOrder := []int{5, 6, 9}
	order := []int{7, 8, 10}
	for i, e := range m.ToSlice() {
		k := e.Key
		v := e.Value
		if keyOrder[i] != k || order[i] != v {
			t.Fatalf("%d != %d || %d != %d", keyOrder[i], k, order[i], v)
		}
		i++
	}
}

package store

import (
	"strings"
	"testing"

	"samuellando.com/internal/datatypes"
)

type elem struct {
	string
	id int64
}

var id int64 = 0

func new(s string) elem {
	id++
	return elem{string: s, id: id}
}

func (e elem) Id() int64 {
	return id
}

type testStore struct {
	items []elem
}

func (ts *testStore) GetById(id int64) (elem, error) {
	return ts.items[id], nil
}

func (ts *testStore) GetAll() ([]elem, error) {
	return ts.items, nil
}

func (ts *testStore) Filter(f func(elem) bool) (Store[elem], error) {
	return Filter(ts, f)
}

func (ts *testStore) Group(f func(elem) string) (datatypes.OrderedMap[string, Store[elem]], error) {
	return Group(ts, f)
}

func (ts *testStore) Sort(f func(elem, elem) bool) (Store[elem], error) {
	return Sort(ts, f)
}

func setup() Store[elem] {
	items := []elem{
		new("Monday"),
		new("Tuesday"),
		new("Wednesday"),
		new("Thursday"),
		new("Friday"),
		new("Saturday"),
		new("Sunday"),
	}
	return &testStore{items}
}

func TestFilter(t *testing.T) {
	ts := setup()
	res, _ := ts.Filter(func(s elem) bool {
		return strings.HasPrefix(s.string, "S")
	})
	data, _ := res.GetAll()
	if len(data) != 2 {
		t.Errorf("%v should contain 2 items", data)
	}
}

func TestGroup(t *testing.T) {
	ts := setup()
	res, _ := ts.Group(func(s elem) string {
		return string(s.string[0])
	})
	if res.Len() != 5 {
		t.Errorf("%s should contain 6 groups", res)
	}
	g, _ := res.Get("S")
	data, _ := g.GetAll()
	if len(data) != 2 {
		t.Errorf("%v should contain 2 items", data)
	}
}

func TestSort(t *testing.T) {
	ts := setup()
	res, _ := ts.Sort(func(a, b elem) bool {
		return strings.Compare(a.string, b.string) < 0
	})
	data, _ := res.GetAll()
	if data[0].string != "Friday" {
		t.Error("Order is wrong")
	}
	if data[1].string != "Monday" {
		t.Error("Order is wrong")
	}
	if data[2].string != "Saturday" {
		t.Error("Order is wrong")
	}
	if data[3].string != "Sunday" {
		t.Error("Order is wrong")
	}
	if data[4].string != "Thursday" {
		t.Error("Order is wrong")
	}
	if data[5].string != "Tuesday" {
		t.Error("Order is wrong")
	}
	if data[6].string != "Wednesday" {
		t.Error("Order is wrong")
	}
}

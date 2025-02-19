package store

import (
	"strings"
	"testing"

	"samuellando.com/internal/datatypes"
)

type testStore struct {
	items []string
}

func (ts *testStore) New(d []string) Store[string] {
	return &testStore{d}
}

func (ts *testStore) GetById(id int) (string, error) {
	return ts.items[id], nil
}

func (ts *testStore) GetAll() ([]string, error) {
	return ts.items, nil
}

func (ts *testStore) Filter(f func(string) bool) Store[string] {
	n, _ := Filter(ts, f)
	return n
}

func (ts *testStore) Group(f func(string) string) *datatypes.OrderedMap[string, Store[string]] {
	n, _ := Group(ts, f)
	return n
}

func (ts *testStore) Sort(f func(string, string) bool) Store[string] {
	n, _ := Sort(ts, f)
	return n
}

func setup() Store[string] {
	items := []string{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}
	return &testStore{items}
}

func TestFilter(t *testing.T) {
	ts := setup()
	res := ts.Filter(func(s string) bool {
		return strings.HasPrefix(s, "S")
	})
	data, _ := res.GetAll()
	if len(data) != 2 {
		t.Errorf("%s should contain 2 items", data)
	}
}

func TestGroup(t *testing.T) {
	ts := setup()
	res := ts.Group(func(s string) string {
		return string(s[0])
	})
	if res.Len() != 5 {
		t.Errorf("%s should contain 6 groups", res)
	}
	g, _ := res.Get("S")
	data, _ := g.GetAll()
	if len(data) != 2 {
		t.Errorf("%s should contain 2 items", data)
	}
}

func TestSort(t *testing.T) {
	ts := setup()
	res := ts.Sort(func(a, b string) bool {
		return strings.Compare(a, b) < 0
	})
	data, _ := res.GetAll()
	if data[0] != "Friday" {
		t.Error("Order is wrong")
	}
	if data[1] != "Monday" {
		t.Error("Order is wrong")
	}
	if data[2] != "Saturday" {
		t.Error("Order is wrong")
	}
	if data[3] != "Sunday" {
		t.Error("Order is wrong")
	}
	if data[4] != "Thursday" {
		t.Error("Order is wrong")
	}
	if data[5] != "Tuesday" {
		t.Error("Order is wrong")
	}
	if data[6] != "Wednesday" {
		t.Error("Order is wrong")
	}
}

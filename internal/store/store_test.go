package store

import (
	"strings"
	"testing"
)

type testStore struct {
    items []string
}

func (ts *testStore) GetById(id int) (string, error) {
    return ts.items[id], nil
}

func (ts *testStore) GetAll() ([]string, error) {
    return ts.items, nil
}

func (ts *testStore) Add(s string) error {
    ts.items = append(ts.items, s)
    return nil
}

func (ts *testStore) Remove(s string) error {
    for i, v := range ts.items {
        if v == s {
            ts.items = append(ts.items[:i], ts.items[i+1:]...)
        }
    }
    ts.items = append(ts.items, s)
    return nil
}

func (ts *testStore) Filter(f func(string) bool) Store[string] {
    data, _ := ts.GetAll()
    return &testStore{Filter(data, f)}
}

func (ts *testStore) Group(f func(string) string) map[string]Store[string] {
    data, _ := ts.GetAll()
    groups := Group(data, f)
    res := make(map[string]Store[string])
    for k, elems := range groups {
        res[k] = &testStore{elems}
    }
    return res
}

func (ts *testStore) Sort(f func(string, string) bool) Store[string] {
    data, _ := ts.GetAll()
    return &testStore{Sort(data, f)}
}

func setup() Store[string] {
    items := []string {
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
    if len(res) != 5 {
        t.Errorf("%s should contain 6 groups", res)
    }
    data, _ := res["S"].GetAll()
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

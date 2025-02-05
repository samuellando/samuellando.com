package project

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
    "database/sql"

	"samuellando.com/internal/datatypes"
	"samuellando.com/internal/store"
)

const URL = "https://api.github.com/users/samuellando/repos?per_page=100"
const API_VERSION = "2022-11-28"

type Store struct {
	run func() ([]*Project, error)
}

type Options struct {
    Url string
}

func CreateStore(db *sql.DB, opts... func(*Options)) Store {
    o := Options{Url: URL}
    for _, opt := range opts {
        opt(&o)
    }
    return Store{run: func() ([]*Project, error) {
		return loadProjects(db, o.Url)
	}}
}

func (ps *Store) New(d []*Project) store.Store[*Project] {
    return &Store{run: func() ([]*Project, error) {
		return d, nil
	}}
}

func createErrorStore(err error) *Store {
	return &Store{run: func() ([]*Project, error) {
		return nil, err
	}}
}

func (ps *Store) GetById(id int) (*Project, error) {
	projects, err := ps.run()
	if err != nil {
		return nil, err
	}
	for _, p := range projects {
		if p.Id() == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Project %d does not exist", id)
}

func (ps *Store) GetAll() ([]*Project, error) {
	return ps.run()
}

func (ps *Store) Filter(f func(*Project) bool) store.Store[*Project] {
    n, err := store.Filter(ps, f)
    if err != nil {
        return createErrorStore(err)
    }
    return n
}

func (ps *Store) Group(f func(*Project) string) *datatypes.OrderedMap[string, store.Store[*Project]] {
    n, err := store.Group(ps, f)
    if err != nil {
        m := datatypes.NewOrderedMap[string, store.Store[*Project]]()
        m.Set("", createErrorStore(err))
        return m
    }
    return n
}

func (ps *Store) Sort(f func(*Project, *Project) bool) store.Store[*Project] {
    n, err := store.Sort(ps, f)
    if err != nil {
        return createErrorStore(err)
    }
    return n
}

func loadProjects(db *sql.DB, url string) ([]*Project, error) {
	req, err := createRequest(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to crete request : %s", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get response : %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response code : %d", res.StatusCode)
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body : %s", err)
	}
	return unmarshalResponse(db, bytes)
}

func createRequest(url string) (*http.Request, error) {
	if req, err := http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", API_VERSION)
		return req, nil
	}
}

func unmarshalResponse(db *sql.DB, b []byte) ([]*Project, error) {
	data := make([]*schema, 0)
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal Json : %s", err)
	}
	projects := make([]*Project, 0, len(data))
	for _, d := range data {
        p, err := loadProject(db, d)
        if err != nil {
            return nil, err
        }
		projects = append(projects, p)
	}
	return projects, nil
}

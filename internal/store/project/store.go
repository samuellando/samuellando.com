package project

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"samuellando.com/internal/store"
)

const URL = "https://api.github.com/users/samuellando/repos?per_page=100"
const API_VERSION = "2022-11-28"

type Store struct {
	run func() ([]*Project, error)
}

func CreateStore() Store {
	return Store{run: func() ([]*Project, error) {
		return loadProjects()
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
	return &Store{run: func() ([]*Project, error) {
		all, err := ps.run()
		if err != nil {
			return nil, err
		}
		return store.Filter(all, f), nil
	}}
}

func (ps *Store) Group(f func(*Project) string) map[string]store.Store[*Project] {
	all, err := ps.run()
	if err != nil {
		return make(map[string]store.Store[*Project])
	}
	groups := store.Group(all, f)
	res := make(map[string]store.Store[*Project])
	for k := range groups {
		res[k] = &Store{run: func() ([]*Project, error) {
			all, err := ps.run()
			if err != nil {
				return nil, err
			}
			return store.Group(all, f)[k], nil
		}}
	}
	return res
}

func (ps *Store) Sort(f func(*Project, *Project) bool) store.Store[*Project] {
	return &Store{run: func() ([]*Project, error) {
		all, err := ps.run()
		if err != nil {
			return nil, err
		}
		return store.Sort(all, f), nil
	}}
}

func loadProjects() ([]*Project, error) {
	req, err := createRequest()
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
    return unmarshalResponse(bytes)
}

func createRequest() (*http.Request, error) {
	if req, err := http.NewRequest("GET", URL, nil); err != nil {
		return nil, err
	} else {
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", API_VERSION) 
		return req, nil
	}
}

func unmarshalResponse(b []byte) ([]*Project, error) {
	data := make([]*schema, 0)
    if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal Json : %s", err)
	}
	projects := make([]*Project, 0, len(data))
	for _, d := range data {
        projects = append(projects, createProject(d))
    }
    return projects, nil
}

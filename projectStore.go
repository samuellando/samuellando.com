package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Project interface {
	Id() int
	Title() string
	Description() string
	Created() string
	Pushed() string
	Url() string
}

type ProjectStore interface {
	Projects() ([]Project, error)
	Close()
}

type basicPs struct{}

type project struct {
	id          int
	name        string
	description string
	created     string
	pushed      string
	url         string
}

func (p *project) Id() int {
	return p.id
}

func (p *project) Title() string {
	return p.name
}

func (p *project) Description() string {
	return p.description
}

func (p *project) Created() string {
	return p.created
}

func (p *project) Pushed() string {
	return p.pushed
}

func (p *project) Url() string {
	return p.url
}

func initializeProjectStore() ProjectStore {
	return &basicPs{}
}

func (ps *basicPs) Close() {
}

func (ps *basicPs) Projects() ([]Project, error) {
	url := "https://api.github.com/users/samuellando/repos?per_page=100"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to crete request : %s", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get response : %s", err)
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response code : %d", res.StatusCode)
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body : %s", err)
	}
	data := make([]interface{}, 0)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal Json : %s", err)
	}
	projects := make([]Project, 0, len(data))
	for _, d := range data {
		var id int
		var name string
		var description string
		var created string
		var pushed string
		var url string
		projectData := d.(map[string]interface{})
		id = int(projectData["id"].(float64))
		name = projectData["name"].(string)
		if projectData["description"] != nil {
			description = projectData["description"].(string)
		}
		if projectData["created_at"] != nil {
			created = projectData["created_at"].(string)
		}
		if projectData["pushed_at"] != nil {
			pushed = projectData["pushed_at"].(string)
		}
		if projectData["url"] != nil {
			url = projectData["url"].(string)
		}
		projects = append(projects, &project{id: id, name: name, description: description, created: created, pushed: pushed, url: url})
	}
	return projects, nil
}

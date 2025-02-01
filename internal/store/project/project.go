package project

import (
    "time"
)

type Project struct {
	id          int
	name        string
	description string
	created     time.Time
	pushed      time.Time
	url         string
}

func (p *Project) Id() int {
	return p.id
}

func (p *Project) Title() string {
	return p.name
}

func (p *Project) Description() string {
	return p.description
}

func (p *Project) Created() time.Time {
	return p.created
}

func (p *Project) Pushed() time.Time {
	return p.pushed
}

func (p *Project) Url() string {
	return p.url
}

func createProject(s *schema) *Project {
    return &Project{
        id: s.ID,
        name: s.Name,
        description: s.Description,
        created: s.CreatedAt,
        pushed: s.PushedAt,
        url: s.HTMLURL,
    }
}

package page

import (
	"io/ioutil"
	"os"
	"regexp"
)

type PageError string

func (e PageError) Error() string {
	return string(e)
}

const INVALID_TITLE PageError = "Invalid title, must match /^[a-zA-Z0-9/]+$/"
const PAGE_EXISTS PageError = "Page already exists"
const PAGE_NOT_FOUND PageError = "Page was not found"

type Page interface {
	Load() error
	Save() error
	Add() error
	Remove() error
	Title() string
	Body() []byte
}

func New(directory, title string, body ...[]byte) *txtPage {
        var b []byte
	if len(body) == 0 {
		b = make([]byte, 0)
	} else {
		b = body[0]
	}
	p := &txtPage{directory: directory, title: title, body: b}
	return p
}

type txtPage struct {
	directory string
	title     string
	body      []byte
}

func (p *txtPage) filePath() string {
	return p.directory + "/" + p.title + ".txt"
}

func (p *txtPage) Title() string {
	return p.title
}

func (p *txtPage) Body() []byte {
	return p.body
}

var validPath = regexp.MustCompile("^[a-zA-z0-9/]+.txt$")

func (p *txtPage) Load() error {
	if !validPath.MatchString(p.filePath()) {
		return INVALID_TITLE
	}
	body, err := ioutil.ReadFile(p.filePath())
	if err != nil {
		return PAGE_NOT_FOUND
	}
	p.body = body
	return nil
}

func (p *txtPage) Save() error {
	if !validPath.MatchString(p.filePath()) {
		return INVALID_TITLE
	}
	return ioutil.WriteFile(p.filePath(), p.body, 0600)
}

func (p *txtPage) Add() error {
	if !validPath.MatchString(p.filePath()) {
		return INVALID_TITLE
	}
	_, err := os.Stat(p.filePath())
	if os.IsExist(err) {
		return PAGE_EXISTS
	}
	return p.Save()
}

func (p *txtPage) Remove() error {
	if !validPath.MatchString(p.filePath()) {
		return INVALID_TITLE
	}
	_, err := os.Stat(p.filePath())
	if os.IsNotExist(err) {
		return PAGE_NOT_FOUND
	}
	return os.Remove(p.filePath())
}

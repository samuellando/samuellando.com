package page

import (
	"../user"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
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
	AddUser(user.User)
	WhiteListed(user.User) bool
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

func List(directory string, user *user.User) []string {
	if !regexp.MustCompile("^[a-zA-Z0-9/]+").MatchString(directory) {
		return nil
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil
	}
	fileList := make([]string, 0)
	for _, file := range files {
		fileList = append(fileList, strings.ReplaceAll(file.Name(), ".txt", ""))
	}
        pageList := make([]string, 0)
        var p *txtPage
        for _, page := range fileList {
          p = New(directory, page)
          p.Load()
          if p.WhiteListed(user) {
            pageList = append(pageList, page)
          }
        }
	return pageList
}

type txtPage struct {
	directory string
	title     string
	body      []byte
	users     []string
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

func (p *txtPage) AddUser(user *user.User) {
        if user != nil {
          p.users = append(p.users, (*user).UserName())
        }
}

func (p *txtPage) WhiteListed(user *user.User) bool {
	if len(p.users) == 0 {
		return true
	}
        if user == nil {
          return false
        }
	isWhiteListed := false
	for i := 0; i < len(p.users); i++ {
		if p.users[i] == (*user).UserName() {
			isWhiteListed = true
			break
		}
	}
	return isWhiteListed
}

var validPath = regexp.MustCompile("^[a-zA-z0-9/]+.txt$")

func (p *txtPage) verifyFile() error {
        if !validPath.MatchString(p.filePath()) {
		return INVALID_TITLE
	}
        _, err := os.Stat(p.filePath())
	if os.IsNotExist(err) {
		return PAGE_NOT_FOUND
	} else {
          return PAGE_EXISTS
        }
}

func (p *txtPage) Load() error {
        err := p.verifyFile()
        if err == INVALID_TITLE || err == PAGE_NOT_FOUND {
          return err
        }
        data, _ := ioutil.ReadFile(p.filePath())
	body := []byte(strings.Split(string(data), "\ufb4f")[0])
	users := strings.Split(string(data), "\ufb4f")[1:]
	p.body = body
	p.users = users
	return nil
}

func (p *txtPage) Save() error {
        err := p.verifyFile()
        if err == INVALID_TITLE {
          return err
        }
	data := p.body
	for i := 0; i < len(p.users); i++ {
		data = append(data, append([]byte("\ufb4f"), []byte(p.users[i])...)...)
	}
	return ioutil.WriteFile(p.filePath(), data, 0600)
}

func (p *txtPage) Add() error {
        err := p.verifyFile()
        if err == INVALID_TITLE || err == PAGE_EXISTS {
          return err
        }
	return p.Save()
}

func (p *txtPage) Remove() error {
        err := p.verifyFile()
        if err == INVALID_TITLE || err == PAGE_NOT_FOUND {
          return err
        }
	return os.Remove(p.filePath())
}

package page

import (
    "io/ioutil"
    "regexp"
    )

type PageError string


func (e PageError) Error() string {
  return string(e)
}

const INVALID_TITLE PageError = "Invalid title, must match /^[a-zA-Z0-9]+$/"
const PAGE_EXISTS PageError = "Page already exists"
const PAGE_NOT_FOUND PageError = "Page was not found"

type Page interface {
  Load() error
  Save() error
  Add() error
  GetTitle() string
  GetBody() []byte
}

type TxtPage struct {
  Directory string
  Title string
  Body []byte
}

func (p *TxtPage) GetTitle() string {
  return p.Title
}

func (p *TxtPage) GetBody() []byte {
  return p.Body
}

var validPath = regexp.MustCompile("^[a-zA-z0-9/]+$")

func (p *TxtPage) Load() error {
  if !validPath.MatchString(p.Title) {
    return INVALID_TITLE
  }
  body, err := ioutil.ReadFile(p.Directory+"/"+p.Title+".txt")
  if err != nil {
    return PAGE_NOT_FOUND
  }
  p.Body = body
  return nil
}

func (p *TxtPage) Save() error {
  if !validPath.MatchString(p.Title) {
    return INVALID_TITLE
  }
  return ioutil.WriteFile(p.Directory+"/"+p.Title+".txt", p.Body, 0600)
}

func (p *TxtPage) Add() error {
  if !validPath.MatchString(p.Title) {
    return INVALID_TITLE
  }
  err := p.Load()
  if err == nil {
    return PAGE_EXISTS
  }
  return p.Save()
}

package main

import (
  "encoding/json"
  "io/ioutil"
)

var lists = make(map[string][]string)

func load() {
  if len(lists) > 0 {
    return
  }
  data, err := ioutil.ReadFile("allowed.json")
  if err != nil {
    panic(err)
  }
  if len(data) > 0 {
    err = json.Unmarshal(data, &lists)
    if err != nil {
      panic(err)
    }
  }
}

func save() {
  data, err := json.Marshal(lists)
  if err != nil {
    panic(err)
  }
  err = ioutil.WriteFile("allowed.json", data, 0600)
  if err != nil {
    panic(err)
  }
}

func List(uid string) []string {
  load()
  return lists[uid]
}

func Allowed(uid, path string) bool {
  load()
  allowed := false
  list := List(uid)
  for _, item := range list {
    if path == item {
      allowed = true
    }
  }
  return allowed
}

func Allow(uid, path string) {
  load()
  lists[uid] = append(lists[uid], path)
  save()
}

func DisAllow(uid, path string) {
  load()
  for i := 0; i < len(lists[uid]); i++ {
    if path == lists[uid][i] {
      lists[uid] = append(lists[uid][:i], lists[uid][i+1:]...)
    }
  }
  save()
}

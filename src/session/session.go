package session

import (
  "../user"
  "time"
  "net/http"
  "errors"
  "container/heap"
  "math/rand"
)

const charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type Session interface {
  Activep(http.ResponseWriter) (*user.User, error)
}

const MAX_TIME = time.Hour

var itemTable = make(map[string]*item)

var pq = make(priorityQueue, 0)

type  session struct {
  lastAction time.Time
  sessionId string
  user *user.User
}

func (s *session) Activep(r *http.Request) (*user.User, error) {
  cookie, err := r.Cookie("userSession")
  if err != nil {
    return nil, errors.New("No valid session id found")
  } else {
    i := itemTable[cookie.Value]
    if i == nil {
      return nil, errors.New("User no assigned user")
    } else {
      return &(*(*i.session).user), nil
    }
  }
}

func sessionAction(f func(http.ResponseWriter, *http.Request, *user.User), w http.ResponseWriter, r *http.Request, u *user.User) {
  clearStaleSessions()
  f(w, r, u)
}

func clearStaleSessions() {
  i := heap.Pop(&pq).(*item)
  for ; i.priority().Sub(time.Now().Add(MAX_TIME)) < 0; i = heap.Pop(&pq).(*item) {
    delete(itemTable, i.session.sessionId)
  }
  heap.Push(&pq, i)
}

func create(w http.ResponseWriter, r *http.Request, u *user.User) {
  _, err := r.Cookie("userSession")
  if err == nil {
    destroy(w, r, nil)
  }
  userSession := randomString(100)
  expiration := time.Now().Add(MAX_TIME)
  cookie := &http.Cookie{Name: "userSession", Value: userSession, Expires: expiration}
  http.SetCookie(w, cookie)
  session := &session{lastAction: time.Now(), user: u, sessionId: userSession}
  item := &item{session: session}
  heap.Push(&pq, item)
  itemTable[userSession] = item
}

func destroy(w http.ResponseWriter, r *http.Request, u *user.User) {
  if u == nil {
    cookie, err := r.Cookie("userSession")
    if err == nil {
      i := itemTable[cookie.Value]
      heap.Remove(&pq, i.index)
      delete(itemTable, cookie.Value)
    }
  }
}

func Create(w http.ResponseWriter, r *http.Request, u *user.User) {
  sessionAction(create, w, r, u)
}

func Destroy(w http.ResponseWriter, r *http.Request) {
  sessionAction(destroy, w, r, nil)
}

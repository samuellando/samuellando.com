package session

import (
	"../user"
	"container/heap"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

const charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

const MAX_TIME = time.Hour

var itemTable = make(map[string]*item)

var pq = make(priorityQueue, 0)

type session struct {
	lastAction time.Time
	sessionId  string
	user       *user.User
}

func Active(r *http.Request) (*user.User, error) {
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

func sessionAction(f func(http.ResponseWriter, *http.Request, user.User), w http.ResponseWriter, r *http.Request, u user.User) {
	clearStaleSessions()
	f(w, r, u)
}

func clearStaleSessions() {
	if pq.Len() > 0 {
		i := heap.Pop(&pq).(*item)
		for ; i.priority().Add(MAX_TIME).Sub(time.Now()) < 0 && pq.Len() > 0; i = heap.Pop(&pq).(*item) {
			delete(itemTable, i.session.sessionId)
		}
		heap.Push(&pq, i)
	}
}

func create(w http.ResponseWriter, r *http.Request, u user.User) {
	_, err := r.Cookie("userSession")
	if err == nil {
		destroy(w, r, nil)
	}
	if u == nil {
		return
	}
	userSession := randomString(100)
	expiration := time.Now().Add(MAX_TIME)
	cookie := &http.Cookie{Name: "userSession", Value: userSession, Expires: expiration}
	http.SetCookie(w, cookie)
	session := &session{lastAction: time.Now(), user: &u, sessionId: userSession}
	item := &item{session: session}
	heap.Push(&pq, item)
	itemTable[userSession] = item
}

func destroy(w http.ResponseWriter, r *http.Request, u user.User) {
	if u == nil {
		cookie, err := r.Cookie("userSession")
		if err == nil {
			i := itemTable[cookie.Value]
			if i != nil {
				heap.Remove(&pq, i.index)
				delete(itemTable, cookie.Value)
			}
			cookie.Expires = time.Now()
			http.SetCookie(w, cookie)
		}
	}
}

func Create(w http.ResponseWriter, r *http.Request, u user.User) {
	sessionAction(create, w, r, u)
}

func Destroy(w http.ResponseWriter, r *http.Request) {
	sessionAction(destroy, w, r, nil)
}

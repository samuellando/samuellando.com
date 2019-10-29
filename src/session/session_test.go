package session

import (
	"../user"
	"container/heap"
	"net/http/httptest"
	"testing"
)

func TestCreateDestroyActive(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	_, err := Active(r)
	if err == nil {
		t.Error("Active returned no error on no valid session id found")
	}
	u := user.New("test.db", "testUser")
	Create(w, r, u)
	if len(itemTable) != 1 || pq.Len() != 1 {
		t.Errorf("Items not added to the datastructures")
	}
	cookie := w.Result().Cookies()[0]
	if (*u).UserName() != (*(itemTable[cookie.Value].session.user)).UserName() {
		t.Error("User was not properly wrapped.")
	}
	if heap.Pop(&pq) != itemTable[cookie.Value] {
		t.Errorf("The values in the priority queue and the map dont match")
	}
	heap.Push(&pq, itemTable[cookie.Value])
	r.AddCookie(cookie)
	u2, err2 := Active(r)
	if err2 != nil {
		t.Error("Active should not return an error with an active user")
	}
	if (*u2).UserName() != (*u).UserName() {
		t.Error("Active did not return the correct user")
	}
	Destroy(w, r)
	if len(itemTable) != 0 || pq.Len() != 0 {
		t.Errorf("The sessions where not completely removed from memory.")
	}
}

package session

import (
  "testing"
  "time"
  "container/heap"
)

func TestLenPushPop(t *testing.T) {
  s1 := &session{lastAction: time.Now()}
  s2 := &session{lastAction: time.Now().Add(time.Minute)}
  i1 := &item{session: s1}
  i2 := &item{session: s2}
  pq := make(priorityQueue, 0)
  heap.Push(&pq, i2)
  heap.Push(&pq, i1)
  if pq.Len() != 2 {
    t.Errorf("The pq should have a length of 2, but got %d", pq.Len())
  }
  i3 := heap.Pop(&pq)
  i4 := heap.Pop(&pq)
  if i3 != i1 {
    t.Errorf("The pq did not pop the lowest priority element")
  }
  if i4 != i2 {
    t.Errorf("The pq did not pop the right element.")
  }
  if pq.Len() != 0 {
    t.Errorf("The pq should have a length of 0, got %d", pq.Len())
  }
}

func TestRemove(t *testing.T) {
  s1 := &session{lastAction: time.Now()}
  s2 := &session{lastAction: time.Now().Add(time.Hour)}
  s3 := &session{lastAction: time.Now().Add(time.Minute)}
  i1 := &item{session: s1}
  i2 := &item{session: s2}
  i3 := &item{session: s3}
  pq := make(priorityQueue, 0)
  heap.Push(&pq, i2)
  heap.Push(&pq, i1)
  heap.Push(&pq, i3)
  heap.Remove(&pq, i3.index)
  if pq.Len() != 2 {
    t.Errorf("The element was not removed")
  }
  i4 := heap.Pop(&pq)
  i5 := heap.Pop(&pq)
   if i4 != i1 || i5 != i2 {
    t.Errorf("The pq removed the wrong element")
  }
}

func TestFix(t *testing.T) {
  s1 := &session{lastAction: time.Now()}
  s2 := &session{lastAction: time.Now().Add(time.Minute)}
  i1 := &item{session: s1}
  i2 := &item{session: s2}
  pq := make(priorityQueue, 0)
  heap.Push(&pq, i2)
  heap.Push(&pq, i1)
  i1.session.lastAction = time.Now().Add(time.Hour)
  heap.Fix(&pq, i1.index)
  i3 := heap.Pop(&pq)
  i4 := heap.Pop(&pq)
  if i3 != i2 || i4 != i1 {
    t.Errorf("The order in the pq was not fixed")
  }
}

package server

import "sync"

type Item struct {
	ID   int         `json:"id"`
	Data interface{} `json:"data"`
}

type Store struct {
	mu    sync.Mutex
	items []Item
	next  int
}

var st = &Store{next: 1}

func (s *Store) add(v interface{}) Item {
	s.mu.Lock()
	it := Item{ID: s.next, Data: v}
	s.next++
	s.items = append(s.items, it)
	s.mu.Unlock()
	return it
}

func (s *Store) all() []Item {
	s.mu.Lock()
	cp := make([]Item, len(s.items))
	copy(cp, s.items)
	s.mu.Unlock()
	return cp
}

func (s *Store) get(id int) (Item, bool) {
	s.mu.Lock()
	for _, it := range s.items {
		if it.ID == id {
			s.mu.Unlock()
			return it, true
		}
	}
	s.mu.Unlock()
	return Item{}, false
}

package memory

import (
	"sync"

	"github.com/pliffdax/sparrow-api/internal/domain"
)

type UserStore struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]domain.User
}

func NewUserStore() *UserStore {
	return &UserStore{items: make(map[int64]domain.User)}
}

func (s *UserStore) Create(name string) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	u := domain.User{ID: s.seq, Name: name}
	s.items[u.ID] = u
	return u, nil
}

func (s *UserStore) GetByID(id int64) (domain.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.items[id]
	return u, ok
}

func (s *UserStore) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}

func (s *UserStore) List() []domain.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.User, 0, len(s.items))
	for _, u := range s.items {
		out = append(out, u)
	}
	return out
}

package memory

import (
	"sync"

	"github.com/pliffdax/sparrow-api/internal/domain"
)

type CategoryStore struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]domain.Category
}

func NewCategoryStore() *CategoryStore {
	return &CategoryStore{items: make(map[int64]domain.Category)}
}

func (s *CategoryStore) Create(title string) (domain.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	c := domain.Category{ID: s.seq, Title: title}
	s.items[c.ID] = c
	return c, nil
}

func (s *CategoryStore) GetByID(id int64) (domain.Category, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.items[id]
	return c, ok
}

func (s *CategoryStore) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}

func (s *CategoryStore) List() []domain.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Category, 0, len(s.items))
	for _, c := range s.items {
		out = append(out, c)
	}
	return out
}

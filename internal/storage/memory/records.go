package memory

import (
	"sync"
	"time"

	"github.com/pliffdax/sparrow-api/internal/domain"
)

type RecordStore struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]domain.Record
}

func NewRecordStore() *RecordStore {
	return &RecordStore{items: make(map[int64]domain.Record)}
}

func (s *RecordStore) Create(r domain.Record) (domain.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	r.ID = s.seq
	s.items[r.ID] = r
	return r, nil
}

func (s *RecordStore) GetByID(id int64) (domain.Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.items[id]
	return rec, ok
}

func (s *RecordStore) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}

func (s *RecordStore) Query(userID, categoryID int64) []domain.Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.Record, 0, len(s.items))
	for _, rec := range s.items {
		if userID > 0 && categoryID > 0 {
			if rec.UserID == userID && rec.CategoryID == categoryID {
				out = append(out, rec)
			}
			continue
		}
		if userID > 0 && rec.UserID == userID {
			out = append(out, rec)
			continue
		}
		if categoryID > 0 && rec.CategoryID == categoryID {
			out = append(out, rec)
			continue
		}
	}
	return out
}

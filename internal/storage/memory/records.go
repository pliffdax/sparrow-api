package memory

import (
	"sync"

	"github.com/pliffdax/sparrow-api/internal/domain"
)

type RecordStore struct {
	mu    sync.Mutex
	seq   int64
	items map[int64]domain.Record
}

func NewRecordStore() *RecordStore {
	return &RecordStore{items: make(map[int64]domain.Record)}
}

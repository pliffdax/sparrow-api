package memory

import (
	"sync"

	"github.com/pliffdax/sparrow-api/internal/domain"
)

type CategoryStore struct {
	mu    sync.Mutex
	seq   int64
	items map[int64]domain.Category
}

func NewCategoryStore() *CategoryStore {
	return &CategoryStore{items: make(map[int64]domain.Category)}
}

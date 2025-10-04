package memory

import (
	"sync"

	"github.com/pliffdax/sparrow-api/internal/domain"
)

type UserStore struct {
	mu    sync.Mutex
	seq   int64
	items map[int64]domain.User
}

func NewUserStore() *UserStore {
	return &UserStore{items: make(map[int64]domain.User)}
}

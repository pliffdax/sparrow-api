package storage

import "github.com/pliffdax/sparrow-api/internal/domain"

type UserStore interface {
	Create(name string) (domain.User, error)
	GetByID(id int64) (domain.User, bool)
	Delete(id int64) bool
	List() []domain.User
}

type CategoryStore interface {
	Create(title string) (domain.Category, error)
	GetByID(id int64) (domain.Category, bool)
	Delete(id int64) bool
	List() []domain.Category
}

type RecordStore interface {
	Create(r domain.Record) (domain.Record, error)
	GetByID(id int64) (domain.Record, bool)
	Delete(id int64) bool
	Query(userID, categoryID int64) []domain.Record
}

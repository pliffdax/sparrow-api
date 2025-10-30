package postgres

import (
	"gorm.io/gorm"
	"github.com/pliffdax/sparrow-api/internal/domain"
)

type UserStore struct{ db *gorm.DB }

func NewUserStore(db *gorm.DB) *UserStore { return &UserStore{db: db} }

func (s *UserStore) Create(name string) (domain.User, error) {
	u := domain.User{Name: name}
	if err := s.db.Create(&u).Error; err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (s *UserStore) GetByID(id int64) (domain.User, bool) {
	var u domain.User
	if err := s.db.First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, false
		}
		return domain.User{}, false
	}
	return u, true
}

func (s *UserStore) Delete(id int64) bool {
	if err := s.db.Delete(&domain.User{}, id).Error; err != nil {
		return false
	}
	return true
}

func (s *UserStore) List() []domain.User {
	var out []domain.User
	_ = s.db.Order("id").Find(&out).Error
	return out
}

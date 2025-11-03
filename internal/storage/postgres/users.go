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
	out := make([]domain.User, 0)
	_ = s.db.Order("id").Find(&out).Error
	return out
}

func (s *UserStore) CreateWithPassword(name, passwordHash string) (domain.User, error) {
	u := domain.User{Name: name}
	tx := s.db.Begin()
	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return domain.User{}, err
	}
	if err := tx.Model(&u).Update("password_hash", passwordHash).Error; err != nil {
		tx.Rollback()
		return domain.User{}, err
	}
	if err := tx.Commit().Error; err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (s *UserStore) FindAuth(name string) (int64, string, error) {
	var id int64
	var ph string
	row := s.db.Raw(`SELECT id, password_hash FROM users WHERE name = ? LIMIT 1`, name).Row()
	if err := row.Scan(&id, &ph); err != nil {
		return 0, "", err
	}
	return id, ph, nil
}

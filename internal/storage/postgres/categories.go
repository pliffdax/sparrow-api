package postgres

import (
	"gorm.io/gorm"
	"github.com/pliffdax/sparrow-api/internal/domain"
)

type CategoryStore struct{ db *gorm.DB }

func NewCategoryStore(db *gorm.DB) *CategoryStore { return &CategoryStore{db: db} }

func (s *CategoryStore) Create(title string) (domain.Category, error) {
	c := domain.Category{Title: title}
	if err := s.db.Create(&c).Error; err != nil {
		return domain.Category{}, err
	}
	return c, nil
}

func (s *CategoryStore) GetByID(id int64) (domain.Category, bool) {
	var c domain.Category
	if err := s.db.First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Category{}, false
		}
		return domain.Category{}, false
	}
	return c, true
}

func (s *CategoryStore) Delete(id int64) bool {
	if err := s.db.Delete(&domain.Category{}, id).Error; err != nil {
		return false
	}
	return true
}

func (s *CategoryStore) List() []domain.Category {
	var out []domain.Category
	_ = s.db.Order("id").Find(&out).Error
	return out
}

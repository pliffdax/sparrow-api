package postgres

import (
	"time"

	"gorm.io/gorm"
	"github.com/pliffdax/sparrow-api/internal/domain"
)

type RecordStore struct{ db *gorm.DB }

func NewRecordStore(db *gorm.DB) *RecordStore { return &RecordStore{db: db} }

func (s *RecordStore) Create(r domain.Record) (domain.Record, error) {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	if err := s.db.Create(&r).Error; err != nil {
		return domain.Record{}, err
	}
	return r, nil
}

func (s *RecordStore) GetByID(id int64) (domain.Record, bool) {
	var rec domain.Record
	if err := s.db.First(&rec, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Record{}, false
		}
		return domain.Record{}, false
	}
	return rec, true
}

func (s *RecordStore) Delete(id int64) bool {
	if err := s.db.Delete(&domain.Record{}, id).Error; err != nil {
		return false
	}
	return true
}

func (s *RecordStore) Query(userID, categoryID int64) []domain.Record {
	q := s.db.Model(&domain.Record{})
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	var out []domain.Record
	_ = q.Order("created_at DESC").Find(&out).Error
	return out
}

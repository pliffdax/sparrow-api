package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pliffdax/sparrow-api/internal/domain"
	"github.com/pliffdax/sparrow-api/internal/storage"
	"github.com/pliffdax/sparrow-api/internal/util"
)

type createRecordReq struct {
	UserID     int64   `json:"user_id"`
	CategoryID int64   `json:"category_id"`
	Amount     float64 `json:"amount"`
	CreatedAt  *string `json:"created_at,omitempty"`
}

func CreateRecord(rs storage.RecordStore, us storage.UserStore, cs storage.CategoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRecordReq
		if err := util.DecodeJSON(r, &req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.UserID <= 0 || req.CategoryID <= 0 {
			http.Error(w, "user_id and category_id must be > 0", http.StatusBadRequest)
			return
		}
		if _, ok := us.GetByID(req.UserID); !ok {
			http.Error(w, "unknown user_id", http.StatusBadRequest)
			return
		}
		if _, ok := cs.GetByID(req.CategoryID); !ok {
			http.Error(w, "unknown category_id", http.StatusBadRequest)
			return
		}

		rec := domain.Record{
			UserID:     req.UserID,
			CategoryID: req.CategoryID,
			Amount:     req.Amount,
		}
		if req.CreatedAt != nil && *req.CreatedAt != "" {
			if t, err := time.Parse(time.RFC3339, *req.CreatedAt); err == nil {
				rec.CreatedAt = t.UTC()
			} else {
				http.Error(w, "created_at must be RFC3339", http.StatusBadRequest)
				return
			}
		}
		rec, _ = rs.Create(rec)
		util.WriteJSON(w, http.StatusCreated, rec)
	}
}

func GetRecord(rs storage.RecordStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		rec, ok := rs.GetByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		util.WriteJSON(w, http.StatusOK, rec)
	}
}

func DeleteRecord(rs storage.RecordStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if !rs.Delete(id) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func QueryRecords(rs storage.RecordStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			uidStr = r.URL.Query().Get("user_id")
			cidStr = r.URL.Query().Get("category_id")
			uid    int64
			cid    int64
			err    error
		)

		if uidStr != "" {
			uid, err = strconv.ParseInt(uidStr, 10, 64)
			if err != nil || uid <= 0 {
				http.Error(w, "user_id must be a positive integer", http.StatusBadRequest)
				return
			}
		}
		if cidStr != "" {
			cid, err = strconv.ParseInt(cidStr, 10, 64)
			if err != nil || cid <= 0 {
				http.Error(w, "category_id must be a positive integer", http.StatusBadRequest)
				return
			}
		}
		if uid == 0 && cid == 0 {
			http.Error(w, "at least one of user_id or category_id is required", http.StatusBadRequest)
			return
		}
		util.WriteJSON(w, http.StatusOK, rs.Query(uid, cid))
	}
}

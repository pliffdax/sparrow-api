package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pliffdax/sparrow-api/internal/security"
	"github.com/pliffdax/sparrow-api/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Users storage.AuthUserStore
}

type registerReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type loginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in registerReq
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" || in.Password == "" {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "hash error", http.StatusInternalServerError)
			return
		}
		u, err := h.Users.CreateWithPassword(in.Name, string(hash))
		if err != nil {
			http.Error(w, "cannot create user: "+err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(u)
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in loginReq
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" || in.Password == "" {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		id, ph, err := h.Users.FindAuth(in.Name)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(ph), []byte(in.Password)) != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		tok, err := security.Sign(fmt.Sprintf("%d", id))
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": tok,
			"token_type":   "Bearer",
		})
	}
}

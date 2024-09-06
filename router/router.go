package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Oyatillohgayratov/fitness-tracking-app/errors"
	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/hash"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage"
)

type UserHandler struct {
	logger  *slog.Logger
	storage storage.Queries
}

func NewMux(logger *slog.Logger, storage *storage.Queries) *http.ServeMux {
	mux := http.NewServeMux()
	u := UserHandler{
		logger:  logger,
		storage: *storage,
	}

	mux.HandleFunc("POST /api/users/register", u.Register)

	return mux
}

func (u UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.logger.Error("failed to decode user registration",
			slog.Any("error", err))
		http.Error(w, errors.ErrDecodeUserRegister.Error(), http.StatusBadRequest)
		return
	}

	password, err := hash.GenerateFromPassword(user.Username)
	if err != nil {
		u.logger.Error("failed to hash password", "error", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	userModel := storage.CreateUserParams{
		Username:     user.Username,
		PasswordHash: password,
		Email:        user.Email,
	}

	resuser, err := u.storage.CreateUser(r.Context(), userModel)
	if err != nil {
		u.logger.Error("failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	res := UserRegisterResponse{
		ID:       int(resuser.ID),
		Username: resuser.Username,
		Email:    resuser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}

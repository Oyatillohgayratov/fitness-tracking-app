package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Oyatillohgayratov/fitness-tracking-app/errors"
	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/email"
	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/hash"
	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/jwt"
	"github.com/Oyatillohgayratov/fitness-tracking-app/models"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage"
)

type UserHandler struct {
	Logger  *slog.Logger
	Storage storage.Queries
}

func NewHandler(logger *slog.Logger, storage *storage.Queries) UserHandler {
	return UserHandler{
		Logger:  logger,
		Storage: *storage,
	}
}

func (u UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.Logger.Error("failed to decode user registration",
			slog.Any("error", err))
		http.Error(w, errors.ErrDecodeUserRegister.Error(), http.StatusBadRequest)
		return
	}

	password, err := hash.GenerateFromPassword(user.Password)
	if err != nil {
		u.Logger.Error("failed to hash password", "error", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	userModel := storage.CreateUserParams{
		Username:     user.Username,
		PasswordHash: password,
		Email:        user.Email,
	}

	resuser, err := u.Storage.CreateUser(r.Context(), userModel)
	if err != nil {
		u.Logger.Error("failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	res := models.UserRegisterResponse{
		ID:       int(resuser.ID),
		Username: resuser.Username,
		Email:    resuser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}

func (u UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	user, err := u.Storage.GetUser(r.Context(), int32(id))
	if err != nil {
		u.Logger.Error("failed to get user", "error", err)
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateUserReq models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := u.Storage.UpdateUser(r.Context(), storage.UpdateUserParams{
		ID:       int32(updateUserReq.ID),
		Username: updateUserReq.Username,
		Email:    updateUserReq.Email,
	})
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "success update user"}`))
	if err != nil {
		u.Logger.Error("failed to write response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (u UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	err = u.Storage.DeleteUser(r.Context(), int32(id))
	if err != nil {
		u.Logger.Error("failed to delete user", "error", err)
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "success dleted user"}`))
	if err != nil {
		u.Logger.Error("failed to write response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (u UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := u.Storage.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "email not found", http.StatusNotFound)
		return
	}

	resetToken, err := jwt.GenerateJWT(int32(user.ID))
	if err != nil {
		http.Error(w, "failed to generate JWT", http.StatusInternalServerError)
		return
	}

	err = u.Storage.SavePasswordResetToken(r.Context(), storage.SavePasswordResetTokenParams{
		UserID: sql.NullInt32{Int32: user.ID, Valid: true},
		Token:  resetToken,
	})
	if err != nil {
		http.Error(w, "failed to save reset token", http.StatusInternalServerError)
		return
	}

	err = email.SendResetEmail(req.Email, resetToken)
	if err != nil {
		u.Logger.Error("failed to send reset email", "error", err)
		http.Error(w, "failed to send reset email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "success"}`))
	if err != nil {
		u.Logger.Error("failed to write response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func (u UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordResetSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resetToken, err := u.Storage.GetPasswordResetToken(r.Context(), req.Token)
	if err != nil {
		http.Error(w, "invalid or expired token", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hash.GenerateFromPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	err = u.Storage.UpdatePassword(r.Context(), storage.UpdatePasswordParams{
		ID:           resetToken.UserID.Int32,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		http.Error(w, "failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "success"}`))
	if err != nil {
		u.Logger.Error("failed to write response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

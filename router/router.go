package router

import (
	"log/slog"
	"net/http"

	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/handlers"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage"
)

func NewMux(logger *slog.Logger, storage *storage.Queries) *http.ServeMux {
	mux := http.NewServeMux()

	u := handlers.NewHandler(logger, storage)

	mux.HandleFunc("POST /api/users/register", u.Register)
	mux.HandleFunc("GET /api/users/get", u.GetUser)
	mux.HandleFunc("PUT /api/users/update", u.UpdateUser)
	mux.HandleFunc("DELETE /api/users/delete", u.DeleteUser)
	mux.HandleFunc("POST /api/users/request_password_reset", u.RequestPasswordReset)
	mux.HandleFunc("PUT /api/users/reset_password", u.ResetPassword)

	mux.HandleFunc("POST /api/workouts", u.CreateWorkout)
	mux.HandleFunc("GET /api/workouts", u.GetWorkoutsByUserID)
	mux.HandleFunc("GET /api/workout", u.GetWorkoutByUserID)

	return mux
}

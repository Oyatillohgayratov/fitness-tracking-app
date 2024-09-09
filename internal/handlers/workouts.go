package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Oyatillohgayratov/fitness-tracking-app/errors"
	"github.com/Oyatillohgayratov/fitness-tracking-app/models"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage"
)

func (u UserHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	workout := models.WorkoutCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		u.Logger.Error("failed to decode user registration",
			slog.Any("error", err))
		http.Error(w, errors.ErrDecodeUserRegister.Error(), http.StatusBadRequest)
		return
	}

	workoutRes, err := u.Storage.CreateWorkout(r.Context(), storage.CreateWorkoutParams{
		UserID:      workout.UserID,
		Name:        workout.Name,
		Description: workout.Description,
	})
	if err != nil {
		u.Logger.Error("failed to create workout", "error", err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	res := models.WorkoutCreateResponse{
		ID:          workoutRes.ID,
		UserID:      workoutRes.UserID,
		Name:        workoutRes.Name,
		Description: workoutRes.Description,
		Date:        workoutRes.Date,
		CreateAt:    workoutRes.CreateAt,
		UpdateAt:    workoutRes.UpdateAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}

func (u UserHandler) GetWorkoutsByUserID(w http.ResponseWriter, r *http.Request) {
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

	workouts, err := u.Storage.GetWorkoutsByUserID(r.Context(), int32(id))
	if err != nil {
		u.Logger.Error("failed to get workouts", "error", err)
		http.Error(w, "failed to get workouts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

func (u UserHandler) GetWorkoutByUserID(w http.ResponseWriter, r *http.Request) {
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
	userId := r.FormValue("user_id")
	if userId == "" {
		http.Error(w, "missing user_id parameter", http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "invalid user_id parameter", http.StatusBadRequest)
		return
	}

	workout, err := u.Storage.GetWorkoutByUserID(r.Context(), storage.GetWorkoutByUserIDParams{ID: int32(id), UserID: int32(userIdInt)})
	if err != nil {
		u.Logger.Error("failed to get workout", "error", err)
		http.Error(w, "failed to get workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)

}

func (u UserHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
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
	useridStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	userid, err := strconv.Atoi(useridStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	var updateWorkout models.WorkoutUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateWorkout); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = u.Storage.UpdateWorkout(r.Context(), storage.UpdateWorkoutParams{
		ID:          int32(id),
		UserID:      int32(userid),
		Name:        updateWorkout.Name,
		Description: updateWorkout.Description,
		Date:        updateWorkout.Date,
	})
	if err != nil {
		u.Logger.Error("failed to update workout", "error", err)
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}
}

func (u UserHandler) DeleteWorkout(w http.ResponseWriter, r http.Request) {
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
	useridStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	userid, err := strconv.Atoi(useridStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	var updateWorkout models.WorkoutUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateWorkout); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = u.Storage.DeleteWorkout(r.Context(), storage.DeleteWorkoutParams{
		ID:     int32(id),
		UserID: int32(userid),
	})
	if err != nil {
		u.Logger.Error("failed to delete workout", "error", err)
		http.Error(w, "failed to delete workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

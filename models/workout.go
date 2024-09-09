package models

import (
	"database/sql"
	"time"
)

type WorkoutCreateRequest struct {
	UserID      int32          `json:"user_id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

type WorkoutCreateResponse struct {
	ID          int32          `json:"id"`
	UserID      int32          `json:"user_id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Date        time.Time      `json:"date"`
	CreateAt    time.Time      `json:"created_at"`
	UpdateAt    time.Time      `json:"updated_at"`
}

type WorkoutUpdateRequest struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Date        time.Time      `json:"date"`
}

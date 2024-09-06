package main

import (
	"context"
	"log/slog"
	"os"

	configloader "github.com/Oyatillohgayratov/config-loader"
	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/config"
	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/server"
	"github.com/Oyatillohgayratov/fitness-tracking-app/router"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage/postgres"
	_ "github.com/lib/pq"
)

var queries *storage.Queries
var logger *slog.Logger

func main() {
	cfg := config.Config{}
	ctx := context.Background()

	err := configloader.LoadYAMLConfig("config.yaml", &cfg)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	connstring := cfg.LoadConfig()

	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := postgres.New(connstring)
	if err != nil {
		logger.Error("Failed to open database")
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping(ctx)
	if err != nil {
		logger.Error("Failed to ping database")
		os.Exit(1)
	}
	queries = storage.New(db)

	mux := router.NewMux(logger, queries)

	srv := server.New(cfg.GetHostPrort(), mux, *logger)
	if err := srv.Run(); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

// func ListUsers(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	users, err := queries.ListUser(ctx)
// 	if err != nil {
// 		logger.Error("Failed to list users", "error", err)
// 		http.Error(w, "Failed to list users", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(users)
// }

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Path[len("/users/"):]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	ctx := context.Background()
// 	user, err := queries.GetUser(ctx, int32(id))
// 	if err != nil {
// 		logger.Error("Failed to get user", "error", err)
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(user)
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Path[len("/users/"):]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	var user storage.UpdateUserParams
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}
// 	user.ID = int32(id)

// 	ctx := context.Background()
// 	err = queries.UpdateUser(ctx, user)
// 	if err != nil {
// 		logger.Error("Failed to update user", "error", err)
// 		http.Error(w, "Failed to update user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func deleteUser(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Path[len("/users/"):]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	ctx := context.Background()
// 	err = queries.DeleteUser(ctx, int32(id))
// 	if err != nil {
// 		logger.Error("Failed to delete user", "error", err)
// 		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		logger.Info("Request received", "method", r.Method, "url", r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }

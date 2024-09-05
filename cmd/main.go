package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	configloader "github.com/Oyatillohgayratov/config-loader"
	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/config"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage"
	_ "github.com/lib/pq"
)

var queries *storage.Queries
var logger *slog.Logger

func main() {
	cfg := config.Config{}
	// cfg.Postgres.Database = "fitness"
	// cfg.Postgres.Host = "localhost"
	// cfg.Postgres.Password = "azamat"
	// cfg.Postgres.Port = "5432"
	// cfg.Postgres.Username = "postgres"

	err := configloader.LoadYAMLConfig("config.yaml", &cfg)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	connstring := cfg.LoadConfig()

	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := sql.Open("postgres", connstring)
	if err != nil {
		logger.Error("Failed to open database")
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error("Failed to ping database")
		os.Exit(1)
	}

	// ctx := context.Background()
	queries = storage.New(db)

	// m := map[string]any{
	// 	"age": 10,
	// 	"bio": "string",
	// }

	// b, err := json.Marshal(m)
	// if err != nil {
	// 	logger.Error("Failed to marshal json", "error", err)
	// 	os.Exit(1)
	// }

	// err = queries.CreateUser(ctx, storage.CreateUserParams{
	// 	Username:     sql.NullString{String: "test", Valid: true},
	// 	Email:        sql.NullString{String: "test@gmail.com", Valid: true},
	// 	PasswordHash: sql.NullString{String: "password123", Valid: true},
	// 	Profile:      pqtype.NullRawMessage{RawMessage: b, Valid: true},
	// })
	// if err != nil {
	// 	logger.Error("Failed to create user", "error", err)
	// 	os.Exit(1)
	// }

	// err = queries.UpdateUser(ctx, storage.UpdateUserParams{
	// 	ID:       1,
	// 	Username: sql.NullString{String: "new_username", Valid: true},
	// 	Email:    sql.NullString{String: "new_email@gmail.com", Valid: true},
	// 	Profile:  pqtype.NullRawMessage{RawMessage: []byte(`{"new_key": "new_value"}`), Valid: true},
	// })
	// if err!= nil {
	//     logger.Error("Failed to update user", "error", err)
	//     os.Exit(1)
	// }

	// err = queries.DeleteUser(ctx, int32(1))
	// if err != nil {
	// 	logger.Error("Failed to delete user", "error", err)
	//     os.Exit(1)
	// }

	// users, err := queries.ListUser(ctx)
	// if err != nil {
	// 	logger.Error("Failed to list users", "error", err)
	// 	os.Exit(1)
	// }

	// for _, user := range users {
	// 	s := user.Profile.RawMessage
	// 	fmt.Printf("User: %+v\n", string(s))
	// }

	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)

	http.Handle("/", LoggingMiddleware(http.DefaultServeMux))

	port := cfg.Server.Port 
	logger.Info("Server is running on port " + port)
	http.ListenAndServe(":"+port, nil)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ListUsers(w, r)
	case "POST":
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUser(w, r)
	case "PUT":
		updateUser(w, r)
	case "DELETE":
		deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users, err := queries.ListUser(ctx)
	if err != nil {
		logger.Error("Failed to list users", "error", err)
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user storage.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := queries.CreateUser(ctx, user)
	if err != nil {
		logger.Error("Failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	user, err := queries.GetUser(ctx, int32(id))
	if err != nil {
		logger.Error("Failed to get user", "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user storage.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	user.ID = int32(id)

	ctx := context.Background()
	err = queries.UpdateUser(ctx, user)
	if err != nil {
		logger.Error("Failed to update user", "error", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = queries.DeleteUser(ctx, int32(id))
	if err != nil {
		logger.Error("Failed to delete user", "error", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request received", "method", r.Method, "url", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

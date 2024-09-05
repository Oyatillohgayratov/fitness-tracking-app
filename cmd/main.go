package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Oyatillohgayratov/fitness-tracking-app/internal/config"
	"github.com/Oyatillohgayratov/fitness-tracking-app/storage"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Config{}
	cfg.Postgres.Database = "fitness"
	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Password = "azamat"
	cfg.Postgres.Port = "5432"
	cfg.Postgres.Username = "postgres"

	connstring := cfg.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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

	ctx := context.Background()
	queries := storage.New(db)

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

	users, err := queries.ListUser(ctx)
	if err != nil {
		logger.Error("Failed to list users", "error", err)
		os.Exit(1)
	}

	for _, user := range users {
		s := user.Profile.RawMessage
		fmt.Printf("User: %+v\n", string(s))
	}

}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

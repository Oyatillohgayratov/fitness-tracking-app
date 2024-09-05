package config

import "fmt"

type Config struct {
	Postgres struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
		SSLMode  string
	}
}

func (c Config) LoadConfig() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Postgres.Username,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database,
		c.Postgres.SSLMode,
	)
}

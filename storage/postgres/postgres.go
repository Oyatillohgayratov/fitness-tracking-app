package postgres

import (
	"context"
	"database/sql"
)

type Postgres struct {
	*sql.DB
}

func New(connectionString string) (*Postgres, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.DB.PingContext(ctx)
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}

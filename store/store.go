package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	config         *Config
	db             *pgxpool.Pool
	UserRepository *UserRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Start() error {
	db, err := pgxpool.New(context.Background(), s.config.DbConnString)
	println(s.config.DbConnString)
	if err != nil {
		return err
	}

	err = db.Ping(context.Background())

	if err != nil {
		return err
	}
	println("BD Connected")
	s.db = db
	return nil
}

func (s *Store) Stop() {
	s.db.Close()
}

func (s *Store) User() *UserRepository {
	if s.UserRepository == nil {
		s.UserRepository = &UserRepository{
			store: s,
		}
	}

	return s.UserRepository
}

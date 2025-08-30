package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	config             *Config
	db                 *pgxpool.Pool
	userRepository     *UserRepository
	tagRepository      *TagRepository
	categoryRepository *CategoryRepository
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

func (s *Store) UserRepository() *UserRepository {
	if s.userRepository == nil {
		s.userRepository = NewUserRepository(s.db)
	}

	return s.userRepository
}

func (s *Store) TagRepository() *TagRepository {
	if s.tagRepository == nil {
		s.tagRepository = NewTagRepository(s.db)
	}

	return s.tagRepository
}

func (s *Store) CategoryRepository() *CategoryRepository {
	if s.categoryRepository == nil {
		s.categoryRepository = NewCategoryRepository(s.db)
	}

	return s.categoryRepository
}

package store

import (
	"arabic/internal/repository"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	config             *Config
	db                 *pgxpool.Pool
	userRepository     *repository.UserRepository
	tagRepository      *repository.TagRepository
	categoryRepository *repository.CategoryRepository
	catalogRepository  *repository.CatalogRepository
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

func (s *Store) UserRepository() *repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = repository.NewUserRepository(s.db)
	}

	return s.userRepository
}

func (s *Store) TagRepository() *repository.TagRepository {
	if s.tagRepository == nil {
		s.tagRepository = repository.NewTagRepository(s.db)
	}

	return s.tagRepository
}

func (s *Store) CategoryRepository() *repository.CategoryRepository {
	if s.categoryRepository == nil {
		s.categoryRepository = repository.NewCategoryRepository(s.db)
	}

	return s.categoryRepository
}

func (s *Store) CatalogRepository() *repository.CatalogRepository {
	if s.catalogRepository == nil {
		s.catalogRepository = repository.NewCatalogRepository(s.db)
	}
	return s.catalogRepository
}

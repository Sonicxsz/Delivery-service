package repository

import (
	"arabic/internal/model"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: pool,
	}
}

type IUserRepository interface {
	Create(cxt context.Context, u *model.User) (*model.User, error)
	FindByEmail(cxt context.Context, email string) (*model.User, error)
}

var (
	insertUser        = "INSERT INTO public.users (email, username, password) VALUES ($1, $2, $3) RETURNING id"
	searchUserByEmail = "SELECT id, username, password, email FROM public.users WHERE email = $1"
)

func (ur *UserRepository) Create(cxt context.Context, u *model.User) (*model.User, error) {
	err := ur.db.QueryRow(cxt, insertUser, u.Email, u.Username, u.Password).Scan(&u.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) FindByEmail(cxt context.Context, email string) (*model.User, error) {
	user := model.User{}
	err := ur.db.QueryRow(cxt, searchUserByEmail, email).Scan(&user.Id, &user.Username, &user.Password, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

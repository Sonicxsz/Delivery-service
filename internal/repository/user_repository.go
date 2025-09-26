package repository

import (
	"arabic/internal/dto"
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
	Create(cxt context.Context, u *dto.UserCreateRequest) error
	FindByEmail(cxt context.Context, email string) (*model.UserFullInfo, error)
	Update(ctx context.Context, query string, values []any) (bool, error)
}

var (
	insertUser        = "INSERT INTO public.users (email, username, password) VALUES ($1, $2, $3) RETURNING id"
	searchUserByEmail = "SELECT id, username, password, email, role_code, first_name, second_name, phone_number, apartment, house, street, city, region   FROM public.users WHERE email = $1"
)

func (ur *UserRepository) Create(cxt context.Context, u *dto.UserCreateRequest) error {
	user := model.User{}
	err := ur.db.QueryRow(cxt, insertUser, u.Email, u.Username, u.Password).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) FindByEmail(cxt context.Context, email string) (*model.UserFullInfo, error) {
	u := model.UserFullInfo{}
	err := ur.db.QueryRow(cxt, searchUserByEmail, email).Scan(&u.Id, &u.Username, &u.Password, &u.Email, &u.RoleCode, &u.FirstName, &u.SecondName, &u.PhoneNumber, &u.Apartment, &u.House, &u.Street, &u.City, &u.Region)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) Update(ctx context.Context, query string, values []any) (bool, error) {

	tag, err := ur.db.Exec(ctx, query, values...)

	if err != nil {
		return false, err
	}

	return tag.RowsAffected() != 0, nil
}

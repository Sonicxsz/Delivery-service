package store

import (
	"arabic/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type CategoryRepository struct {
	store *Store
}

var (
	findAllCategory = "select * from test.category order by id"
	createCategory  = "insert into test.category (name, code) values ($1, $2) returning id, name, code"
	deleteCategory  = "delete from test.category where id = $1"
)

func (t *CategoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	query, err := t.store.db.Query(ctx, findAllCategory)
	if err != nil {
		return nil, err
	}

	var categories []*model.Category
	for query.Next() {
		category := &model.Category{}

		err = query.Scan(&category.Id, &category.Name, &category.Code)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (t *CategoryRepository) Delete(ctx context.Context, id int64) (*pgconn.CommandTag, error) {
	result, err := t.store.db.Exec(ctx, deleteCategory, id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *CategoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	created := &model.Category{}
	err := t.store.db.QueryRow(ctx, createCategory, category.Name, category.Code).Scan(&created.Id, &created.Name, &created.Code)
	if err != nil {
		return nil, err
	}

	return created, nil
}

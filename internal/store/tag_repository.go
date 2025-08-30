package store

import (
	"arabic/internal/model"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5/pgconn"
)

type TagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db: db}
}

type ITagRepository interface {
	FindAll(ctx context.Context) ([]*model.Tag, error)
	Delete(ctx context.Context, id int64) (*pgconn.CommandTag, error)
	Create(ctx context.Context, tag *model.Tag) (*model.Tag, error)
}

var (
	findAllTags = "select * from tags order by id"
	createTag   = "insert into tags (name) values ($1) RETURNING id, name"
	deleteTag   = "delete from tags where id = $1"
)

func (t *TagRepository) FindAll(ctx context.Context) ([]*model.Tag, error) {
	query, err := t.db.Query(ctx, findAllTags)

	if err != nil {
		return nil, err
	}

	var tags []*model.Tag
	for query.Next() {
		tag := &model.Tag{}

		err = query.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (t *TagRepository) Delete(ctx context.Context, id int64) (*pgconn.CommandTag, error) {
	result, err := t.db.Exec(ctx, deleteTag, id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TagRepository) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	created := &model.Tag{}
	err := t.db.QueryRow(ctx, createTag, tag.Name).Scan(&created.Id, &created.Name)
	if err != nil {
		return nil, err
	}

	return created, nil
}

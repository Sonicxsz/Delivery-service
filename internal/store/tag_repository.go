package store

import (
	"arabic/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type TagRepository struct {
	store *Store
}

var (
	findAllTags = "select * from test.tags order by id"
	createTag   = "insert into test.tags (name) values ($1) RETURNING id, name"
	deleteTag   = "delete from test.tags where id = $1"
)

func (t *TagRepository) FindAll(ctx context.Context) ([]*model.Tag, error) {
	query, err := t.store.db.Query(ctx, findAllTags)

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
	result, err := t.store.db.Exec(ctx, deleteTag, id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TagRepository) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	created := &model.Tag{}
	err := t.store.db.QueryRow(ctx, createTag, tag.Name).Scan(&created.Id, &created.Name)
	if err != nil {
		return nil, err
	}

	return created, nil
}

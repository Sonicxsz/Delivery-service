package store

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"context"
	"fmt"
)

type TagRepository struct {
	store *Store
}

var (
	findAll = "select * from test.tags order by id"
	create  = "insert into test.tags (name) values ($1) RETURNING id"
)

func (t *TagRepository) FindAll(ctx context.Context) ([]model.Tag, error) {
	query, err := t.store.db.Query(ctx, findAll)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var tags []model.Tag
	for query.Next() {
		var tag model.Tag
		err = query.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil

}

//
//func (t *TagRepository) Delete(ctx context.Context, id int64) (*model.Tag, error) {
//
//}

func (t *TagRepository) Create(ctx context.Context, req *dto.TagRequest) (*dto.TagRequest, error) {
	err := t.store.db.QueryRow(ctx, create, req.Name).Scan(&req.Id)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return req, nil
}

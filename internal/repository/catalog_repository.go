package repository

import (
	"arabic/internal/model"
	"arabic/pkg/logger"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct {
	db *pgxpool.Pool
}

type ICatalogRepository interface {
	FindAll(ctx context.Context) ([]*model.Catalog, error)
	Delete(ctx context.Context, id uint) (bool, error)
	Create(ctx context.Context, category *model.Catalog) (*model.Catalog, error)
	Update(ctx context.Context, query string, values []any) (bool, error)
	FindById(ctx context.Context, id uint) (*model.Catalog, bool, error)
}

func NewCatalogRepository(db *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{
		db: db,
	}
}

func (c *CatalogRepository) Update(ctx context.Context, queryParts string, values []any) (bool, error) {
	query := "update public.catalogs set " + " " + queryParts + " WHERE id = $1"
	tag, err := c.db.Exec(ctx, query, values...)

	if err != nil {
		return false, err
	}

	return tag.RowsAffected() != 0, nil
}

func (c *CatalogRepository) FindById(ctx context.Context, id uint) (*model.Catalog, bool, error) {
	query := "SELECT id, name, price, discount_percent, amount, category_id, description, sku FROM public.catalogs WHERE id = $1"

	item := &model.Catalog{}

	err := c.db.QueryRow(ctx, query, id).Scan(
		&item.Id,
		&item.Name,
		&item.Price,
		&item.DiscountPercent,
		&item.Amount,
		&item.CategoryId,
		&item.Description,
		&item.Sku,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return item, true, nil
}

func (c *CatalogRepository) Delete(ctx context.Context, id uint) (bool, error) {
	query := "delete from public.catalogs where id = $1"
	tag, err := c.db.Exec(ctx, query, id)

	if err != nil {
		return false, err
	}

	return tag.RowsAffected() != 0, nil
}

func (c *CatalogRepository) FindAll(ctx context.Context) ([]*model.Catalog, error) {
	query := "SELECT id, name, price, discount_percent, amount, category_id FROM public.catalogs ORDER BY id"
	rows, err := c.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	var catalogItems []*model.Catalog
	for rows.Next() {
		item := &model.Catalog{}
		err = rows.Scan(&item.Id, &item.Name, &item.Price, &item.DiscountPercent, &item.Amount, &item.CategoryId)
		if err != nil {
			logger.Log.Error("Catalog repository -> FindAll -> error: " + err.Error())
			continue
		}
		catalogItems = append(catalogItems, item)
	}

	return catalogItems, nil

}

func (c *CatalogRepository) Create(ctx context.Context, ci *model.Catalog) (*model.Catalog, error) {
	query := "insert into public.catalogs (name, description, price, amount, discount_percent, sku, category_id) values ($1, $2, $3, $4, $5, $6, $7) returning id"
	err := c.db.QueryRow(ctx, query,
		ci.Name,
		ci.Description,
		ci.Price,
		ci.Amount,
		ci.DiscountPercent,
		ci.Sku,
		ci.CategoryId).Scan(&ci.Id)

	if err != nil {
		return ci, err
	}

	return ci, nil

}

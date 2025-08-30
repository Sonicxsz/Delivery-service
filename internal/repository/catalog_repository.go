package repository

import (
	"arabic/internal/model"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct {
	db *pgxpool.Pool
}

type ICatalogRepository interface {
	FindAll(ctx context.Context) ([]*model.Catalog, error)
	//FindById(ctx context.Context, ids []int64) ([]*model.Catalog, error)
	Delete(ctx context.Context, id int64) (bool, error)
	Create(ctx context.Context, category *model.Catalog) (*model.Catalog, error)
}

var (
	findAllCatalogItems = "select id, name, price, discount_percent, amount, category_id from catalog order by id"
	createCatalogItem   = "insert into catalog (name, description, price, amount, discount_percent, sku, category_id) values ($1, $2, $3, $4, $5, $6, $7) returning id"
	deleteCatalogItem   = "delete from catalog where id = $1"
)

func NewCatalogRepository(db *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{
		db: db,
	}
}

func (c *CatalogRepository) Delete(ctx context.Context, id int64) (bool, error) {
	tag, err := c.db.Exec(ctx, deleteCatalogItem, id)

	if err != nil {
		return false, err
	}

	return tag.RowsAffected() != 0, nil
}

func (c *CatalogRepository) FindAll(ctx context.Context) ([]*model.Catalog, error) {
	rows, err := c.db.Query(ctx, findAllCatalogItems)

	if err != nil {
		return nil, err
	}

	var catalogItems []*model.Catalog
	for rows.Next() {
		item := &model.Catalog{}
		err1 := rows.Scan(&item.Id, &item.Name, &item.Price, &item.DiscountPercent, &item.Amount, &item.CategoryId)
		if err1 != nil {
			println(err1.Error())
			continue
		}
		catalogItems = append(catalogItems, item)
	}

	return catalogItems, nil

}

func (c *CatalogRepository) Create(ctx context.Context, ci *model.Catalog) (*model.Catalog, error) {
	err := c.db.QueryRow(ctx, createCatalogItem,
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

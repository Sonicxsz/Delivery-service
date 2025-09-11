package model

import "arabic/internal/dto"

type Catalog struct {
	Id              uint    `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	DiscountPercent float32 `json:"discountPercent"`
	Sku             string  `json:"sku"`
	CategoryId      uint    `json:"categoryId"`
	ImageUrl        string  `json:"imageUrl"`
}

func (c *Catalog) ToResponse(imagePrefix string) *dto.CatalogResponse {
	return &dto.CatalogResponse{
		Id:              c.Id,
		Description:     c.Description,
		Sku:             c.Sku,
		Name:            c.Name,
		Price:           c.Price,
		Amount:          c.Amount,
		CategoryId:      c.CategoryId,
		DiscountPercent: c.DiscountPercent,
		ImageUrl:        imagePrefix + c.ImageUrl,
	}
}

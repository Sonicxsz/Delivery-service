package dto

import "arabic/pkg/validator"

type CatalogBase struct {
	Name            string  `json:"name"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	DiscountPercent float32 `json:"discount_percent"`
	CategoryId      int64   `json:"category_id"`
}

type CatalogResponse struct {
	Id int64 `json:"id"`
	CatalogBase
}

type CatalogCreateRequest struct {
	CatalogBase
	Description string `json:"description"`
	Sku         string `json:"sku"`
}

type GetCatalogByIdResponse struct {
	Id int64 `json:"id"`
	CatalogBase
	Description string `json:"description"`
	Sku         string `json:"sku"`
}

type CatalogUpdateRequest struct {
	Id              int64    `json:"id"`
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	Price           *float32 `json:"price"`
	Amount          *int     `json:"amount"`
	DiscountPercent *float32 `json:"discount_percent"`
	Sku             *string  `json:"sku"`
	CategoryId      *int64   `json:"category_id"`
}

type CatalogGetManyRequest struct {
	ID []any `json:"id"`
}

func (c *CatalogGetManyRequest) IsValid() (bool, []string) {
	v := validator.New()
	v.CheckNumber(len(c.ID), "Id").IsMin(1)

	return v.HasErrors(), v.GetErrors()
}

func (c *CatalogCreateRequest) IsValid() (bool, []string) {
	v := validator.Validator{}
	v.CheckString(c.Name, "Name").IsMin(3).IsMax(50)
	v.CheckString(c.Description, "Description").IsMin(20).IsMax(1500)
	v.CheckNumber(c.Price, "Price").IsMin(1).IsMax(100000)
	v.CheckNumber(c.DiscountPercent, "Discount").IsMin(0).IsMax(100)
	v.CheckString(c.Sku, "Sku").IsMin(10).IsMax(64)
	v.CheckNumber(c.CategoryId, "CategoryId").IsMin(1)

	return !v.HasErrors(), v.GetErrors()
}

func (c *CatalogUpdateRequest) IsValid() (bool, []string) {
	v := validator.Validator{}

	if c.Name != nil {
		v.CheckString(*c.Name, "Name").IsMin(3).IsMax(50)
	}
	if c.Description != nil {
		v.CheckString(*c.Description, "Description").IsMin(20).IsMax(1500)
	}
	if c.Price != nil {
		v.CheckNumber(*c.Price, "Price").IsMin(1).IsMax(100000)
	}
	if c.DiscountPercent != nil {
		v.CheckNumber(*c.DiscountPercent, "Discount").IsMax(100)
	}
	if c.Sku != nil {
		v.CheckString(*c.Sku, "Sku").IsMin(10).IsMax(64)
	}
	if c.CategoryId != nil {
		v.CheckNumber(*c.CategoryId, "CategoryId").IsMin(1)
	}

	return !v.HasErrors(), v.GetErrors()
}

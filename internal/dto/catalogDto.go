package dto

import "arabic/pkg/validator"

type CatalogResponse struct {
	Id              uint    `json:"id"`
	Name            string  `json:"name"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	DiscountPercent float32 `json:"discount_percent"`
	CategoryId      uint    `json:"category_id"`
	Description     string  `json:"description"`
	Sku             string  `json:"sku"`
	ImageUrl        string  `json:"imageUrl"`
}

type CatalogCreateRequest struct {
	Name            string  `json:"name"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	DiscountPercent float32 `json:"discount_percent"`
	CategoryId      uint    `json:"category_id"`
	Description     string  `json:"description"`
	Sku             string  `json:"sku"`
}

type CatalogUpdateRequest struct {
	Id              uint     `json:"id"`
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	Price           *float32 `json:"price"`
	Amount          *int     `json:"amount"`
	DiscountPercent *float32 `json:"discount_percent"`
	Sku             *string  `json:"sku"`
	CategoryId      *uint    `json:"category_id"`
}

type AddImageRequest struct {
	Image string `json:"image"`
	Id    uint   `json:"id"`
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

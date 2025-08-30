package dto

import "arabic/pkg/validator"

type CatalogCreateRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	DiscountPercent float32 `json:"discountPercent"`
	Sku             string  `json:"sku"`
	CategoryId      int64   `json:"categoryId"`
}

type CatalogResponse struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	DiscountPercent float32 `json:"discountPercent"`
	CategoryId      int64   `json:"categoryId"`
}

type CatalogUpdateRequest struct {
	Id              int64    `json:"id"`
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	Price           *float32 `json:"price"`
	Amount          *int     `json:"amount"`
	DiscountPercent *float32 `json:"discountPercent"`
	Sku             *string  `json:"sku"`
	CategoryId      *int64   `json:"categoryId"`
}

type CatalogUploadImageRequest struct {
	ImageBase64 string `json:"image_base64"`
}

func (c *CatalogCreateRequest) IsValid() (bool, []string) {
	v := validator.Validator{}
	v.CheckString(c.Name, "Name").IsMin(3).IsMax(50)
	v.CheckString(c.Description, "Description").IsMin(20).IsMax(1500)
	v.CheckNumber(c.Price, "Price").IsMin(1).IsMax(100000)
	v.CheckNumber(c.DiscountPercent, "Discount").IsMax(100)
	v.CheckString(c.Sku, "Sku").IsMin(10).IsMax(64)
	v.CheckNumber(c.CategoryId, "CategoryId").IsMin(1)

	return !v.HasErrors(), v.GetErrors()
}

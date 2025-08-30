package model

type Catalog struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	DiscountPercent float32 `json:"discountPercent"`
	Sku             string  `json:"sku"`
	CategoryId      int64   `json:"categoryId"`
}

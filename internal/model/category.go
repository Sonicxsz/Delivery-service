package model

type Category struct {
	Id         int64  `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
	UsageCount int    `json:"usage_count"`
}

package model

type Tag struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
	UsageCount int    `json:"usage_count"`
	Color      string `json:"color"`
}

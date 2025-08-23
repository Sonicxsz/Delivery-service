package model

type PartOfSpeech struct {
	Id   int64  `json:"id"`
	Code int    `json:"code"`
	Name string `json:"name"`
}

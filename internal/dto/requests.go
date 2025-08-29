package dto

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TagRequest struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type TagResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type CategoryRequest struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoryResponse struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

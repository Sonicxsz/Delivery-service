package dto

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserGetResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type TagRequest struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TagResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	IsActive bool   `json:"is_active"`
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

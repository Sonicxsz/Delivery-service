package dto

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TagRequest struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

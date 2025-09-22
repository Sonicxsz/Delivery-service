package model

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Id       int64  `json:"id"`
	RoleCode string `json:"role_code"`
}

type UserFullInfo struct {
	User
	UserAddress
	FirstName   string `json:"first_name"`
	SecondName  string `json:"second_name"`
	PhoneNumber string `json:"phone_number"`
}

type UserAddress struct {
	Apartment  string `json:"apartment"`
	House      string `json:"house"`
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Region     string `json:"region"`
}

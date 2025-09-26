package dto

import "arabic/pkg/validator"

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserGetResponse struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	Id         int64  `json:"id"`
	SecondName string `json:"second_name"`
}
type UserCreateRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Id         int64
	FirstName  *string `json:"first_name"`
	SecondName *string `json:"second_name"`
}

func (u *UserUpdateRequest) IsValid() (bool, []string) {
	v := validator.New()

	if u.FirstName != nil {
		v.CheckString(*u.FirstName, "FirstName").IsMin(2).IsMax(10)
	}
	if u.SecondName != nil {
		v.CheckString(*u.SecondName, "SecondName").IsMin(4).IsMax(20)
	}

	return !v.HasErrors(), v.GetErrors()
}

type UserAddressUpdateRequest struct {
	Id        int64
	Apartment string `json:"apartment"`
	House     string `json:"house"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Region    string `json:"region"`
}

func (u *UserAddressUpdateRequest) IsValid() (bool, []string) {
	v := validator.New()

	v.CheckString(u.Apartment, "Apartment").IsMin(0).IsMax(10)
	v.CheckString(u.House, "House").IsMin(1).IsMax(5)
	v.CheckString(u.Street, "Street").IsMin(2).IsMax(173)
	v.CheckString(u.City, "City").IsMin(3).IsMax(25)
	v.CheckString(u.Region, "Region").IsMin(4).IsMax(25)

	return !v.HasErrors(), v.GetErrors()
}

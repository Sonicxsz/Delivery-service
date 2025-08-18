package validator

import (
	"regexp"
)

type Validator struct {
	errors []string
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9.!_@#$%^&*].{8,}`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,16}$`)
)

func New() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

func (v *Validator) CheckEmail(email string) {
	if !emailRegex.MatchString(email) {
		v.errors = append(v.errors, "Incorrect email address")
	}
}

func (v *Validator) CheckPassword(password string) {

	for _, re := range []*regexp.Regexp{passwordRegex} {
		if !re.MatchString(password) {
			v.errors = append(v.errors, "Incorrect password")
			return
		}
	}
}

func (v *Validator) CheckUsername(username string) {

	if !usernameRegex.MatchString(username) {
		v.errors = append(v.errors, "Incorrect username")
		return
	}
}

func (v *Validator) HasErrors() (bool, []string) {
	hasErrors := len(v.errors) > 0
	return hasErrors, v.errors
}

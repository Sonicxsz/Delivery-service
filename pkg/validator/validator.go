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
		v.errors = append(v.errors, "Invalid email. Must contain: letters (a-z, A-Z), digits (0-9), or symbols ._%+- before @, followed by a valid domain with a TLD (e.g., .com, .org).")
	}
}

func (v *Validator) CheckPassword(password string) {
	if !passwordRegex.MatchString(password) {
		v.errors = append(v.errors, "Invalid password. Must be at least 8 characters long and can include letters, digits, and symbols .!_@#$%^&*.")
		return
	}
}

func (v *Validator) CheckUsername(username string) {
	if !usernameRegex.MatchString(username) {
		v.errors = append(v.errors, "Invalid username. Must be 3–16 characters long and contain only letters, digits, hyphens (-), or underscores (_).")
		return
	}
}

func (v *Validator) HasErrors() (bool, []string) {
	hasErrors := len(v.errors) > 0
	return hasErrors, v.errors
}

package validator

import (
	"fmt"
	"regexp"
)

type Validator struct {
	errors []string
}

// Проверяет есть ли ошибки
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *Validator) GetErrors() []string {
	return v.errors
}

type StringValidator struct {
	value string
	errs  *Validator
}

type NumberValidator struct {
	value any
	errs  *Validator
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

// Создаем структуру с методами для проверки строк
func (v *Validator) CheckString(value string) *StringValidator {
	return &StringValidator{
		value: value,
		errs:  v,
	}
}

// // Создаем структуру с методами для проверки чисел
func (v *Validator) CheckNumber(value any) *NumberValidator {
	return &NumberValidator{
		value: value,
		errs:  v,
	}
}

// Проверяет что длина строки не больше указанного
func (v *StringValidator) IsMax(max int) *StringValidator {
	if len(v.value) > max {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("Max aviable length is %d, Provided: %s", max, v.value))
	}
	return v
}

// Проверяет что длина строки не меньше указанного
func (v *StringValidator) IsMin(min int) *StringValidator {
	if len(v.value) < min {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("Min required length is %d, Provided: %s", min, v.value))
	}
	return v
}

// Проверяет что строка соответствует шаблону электронной почта Email
func (v *StringValidator) IsEmail() *StringValidator {
	if !emailRegex.MatchString(v.value) {
		v.errs.errors = append(v.errs.errors, "Invalid email. Must contain: letters (a-z, A-Z), digits (0-9), or symbols ._%+- before @, followed by a valid domain with a TLD (e.g., .com, .org).")
	}
	return v
}

// Проверяет что строка соответствует требованиям надежного пароля
func (v *StringValidator) IsPassword() *StringValidator {
	if !passwordRegex.MatchString(v.value) {
		v.errs.errors = append(v.errs.errors, "Invalid password. Must be at least 8 characters long and can include letters, digits, and symbols .!_@#$%^&*.")

	}
	return v
}

// Проверяет что строка соответствует требованиям имени пользователя
func (v *StringValidator) IsValidUsername() *StringValidator {
	if !usernameRegex.MatchString(v.value) {
		v.errs.errors = append(v.errs.errors, "Invalid username. Must be 3–16 characters long and contain only letters, digits, hyphens (-), or underscores (_).")

	}
	return v
}

// Проверяет что число не меньше указанного
func (v *NumberValidator) IsMin(min int) *NumberValidator {
	hasErr, err := false, fmt.Sprintf("Min required: %d, Provided: %d", min, v.value)

	switch val := v.value.(type) {
	case int:
		if val < min {
			hasErr = true
		}
	case int8:
		if int(val) < min {
			hasErr = true
		}
	case int16:
		if int(val) < min {
			hasErr = true
		}
	case int32:
		if int(val) < min {
			hasErr = true
		}
	case int64:
		if val < int64(min) {
			hasErr = true
		}
	default:
		hasErr = true
		err = fmt.Sprintf("Unsupported type: %T", v.value)
	}

	if hasErr {
		v.errs.errors = append(v.errs.errors, err)
	}

	return v
}

// Проверяет что число не больше указанного
func (v *NumberValidator) IsMax(max int) *NumberValidator {
	hasErr, err := false, fmt.Sprintf("Max aviable: %d, Provided: %d", max, v.value)

	switch val := v.value.(type) {
	case int:
		if val > max {
			hasErr = true
		}
	case int8:
		if int(val) < max {
			hasErr = true
		}
	case int16:
		if int(val) < max {
			hasErr = true
		}
	case int32:
		if int(val) < max {
			hasErr = true
		}
	case int64:
		if val < int64(max) {
			hasErr = true
		}
	default:
		hasErr = true
		err = fmt.Sprintf("Unsupported type: %T", v.value)
	}

	if hasErr {
		v.errs.errors = append(v.errs.errors, err)
	}
	return v
}

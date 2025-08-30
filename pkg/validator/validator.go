package validator

import (
	"fmt"
	"regexp"
	"unicode/utf8"
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
	name  string
}

type NumberValidator struct {
	value any
	errs  *Validator
	name  string
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
func (v *Validator) CheckString(value string, name string) *StringValidator {
	return &StringValidator{
		value: value,
		errs:  v,
		name:  name,
	}
}

// // Создаем структуру с методами для проверки чисел
func (v *Validator) CheckNumber(value any, name string) *NumberValidator {
	return &NumberValidator{
		value: value,
		errs:  v,
		name:  name,
	}
}

// Проверяет что длина строки не больше указанного
func (v *StringValidator) IsMax(max int) *StringValidator {
	length := utf8.RuneCountInString(v.value)
	if length > max {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("[%s] - Max aviable length is %d, Provided: %d", v.name, max, length))
	}
	return v
}

// Проверяет что длина строки не меньше указанного
func (v *StringValidator) IsMin(min int) *StringValidator {
	length := utf8.RuneCountInString(v.value)
	if length < min {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("[%s] - Min required length is %d, Provided: %d", v.name, min, length))
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
		v.errs.errors = append(v.errs.errors, "[%s] - Invalid username. Must be 3–16 characters long and contain only letters, digits, hyphens (-), or underscores (_).")

	}
	return v
}

// Проверяет что число не меньше указанного
func (v *NumberValidator) IsMin(min float64) *NumberValidator {
	value, ok := v.toInt64()

	if !ok {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("[%s] - Unsupported type: %T", v.name, v.value))
		return v
	}
	if value < min {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("[%s] - Min required: %g, Provided: %g", v.name, min, value))

	}
	return v
}

// Проверяет что число не больше указанного
func (v *NumberValidator) IsMax(max float64) *NumberValidator {
	value, ok := v.toInt64()

	if !ok {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("[%s] - Unsupported type: %T", v.name, v.value))
		return v
	}

	if value > max {
		v.errs.errors = append(v.errs.errors, fmt.Sprintf("[%s] - Max aviable: %g, Provided: %g", v.name, max, value))
	}
	return v
}

// Преобразует любой числовой тип в int64 для единого сравнения
func (v *NumberValidator) toInt64() (float64, bool) {
	switch val := v.value.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		// Проверяем переполнение
		if val > 9223372036854775807 { // math.MaxInt64
			return 0, false
		}
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

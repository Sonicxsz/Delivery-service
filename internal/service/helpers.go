package service

import "strings"

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}

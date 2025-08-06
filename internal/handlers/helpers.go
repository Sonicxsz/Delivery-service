package handlers

import "net/http"

type Message[T any] struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Data    T      `json:"data,omitempty"`
}

func NewSuccessMessage[T any](msg string, code int, data T) *Message[T] {
	return &Message[T]{
		Message: msg,
		Code:    code,
		Data:    data,
		Error:   false,
	}
}

func NewErrorMessage(msg string, code int) *Message[interface{}] {
	return &Message[interface{}]{
		Message: msg,
		Code:    code,
		Error:   true,
	}
}

func initHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
}

package server

import (
	"net/http"
)

type CustomContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Session *session
	*Server
}
type response struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Error   string   `json:"error,omitempty"`
	Data    any      `json:"data,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type CustomMW func(c *CustomContext) (bool, int, *response)
type CustomHandlerFunc func(c *CustomContext) error

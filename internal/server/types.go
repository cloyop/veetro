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
type ResponseStatus struct {
	Code    int      `json:"code,omitempty"`
	Message string   `json:"message,omitempty"`
	Error   string   `json:"error,omitempty"`
	Errs    []string `json:"errs,omitempty"`
}

type CustomMW func(c *CustomContext) (bool, *ResponseStatus)
type CustomHandlerFunc func(c *CustomContext) error

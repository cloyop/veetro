package server

import (
	"net/http"
	"strings"
)

func corsMid(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		n.ServeHTTP(w, r)
	})
}
func EmployeeOnly(c *CustomContext) (bool, *ResponseStatus) {
	if c.Session.Employee {
		return true, nil
	}
	res := &ResponseStatus{}
	res.Code = http.StatusUnauthorized
	res.Error = "Only for Employees"
	return false, res
}
func CustomerOnly(c *CustomContext) (bool, *ResponseStatus) {
	if c.Session.Customer {
		return true, nil
	}
	res := &ResponseStatus{}
	res.Code = http.StatusUnauthorized
	res.Error = "Only for customers"
	return false, res
}
func Auth(c *CustomContext) (bool, *ResponseStatus) {
	sessionStr, is := getSessionIdFromHeader(c.Request)
	if is {
		session, found := c.Sessions.GetSession(sessionStr)
		if found && session.Valid {
			c.Session = session
			return true, nil
		}
	}
	return false, &ResponseStatus{Code: http.StatusForbidden, Error: "Forbidden"}
}
func getSessionIdFromHeader(r *http.Request) (string, bool) {
	sessionSplit := strings.Split(r.Header.Get("Authorization"), " ")
	if len(sessionSplit) < 2 {
		return "", false
	}
	str := sessionSplit[1]
	if str == "" {
		return "", false
	}
	return str, true
}

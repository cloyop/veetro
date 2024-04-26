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
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		n.ServeHTTP(w, r)
	})
}
func EmployeeOnly(c *CustomContext) (bool, int, *response) {
	if c.Session.Employee {
		return true, 0, nil
	}
	return false, http.StatusUnauthorized, &response{Error: "Only for Employees"}
}
func CustomerOnly(c *CustomContext) (bool, int, *response) {
	if c.Session.Customer {
		return true, 0, nil
	}
	return false, http.StatusUnauthorized, &response{Error: "Only for customers"}
}
func Auth(c *CustomContext) (bool, int, *response) {
	sessionStr, is := getSessionIdFromHeader(c.Request)
	if is {
		session, found := c.Sessions.GetSession(sessionStr)
		if found && session.Valid {
			c.Session = session
			return true, 0, nil
		}
	}
	return false, http.StatusForbidden, &response{Error: "Forbidden"}
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

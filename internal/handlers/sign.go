package handlers

import (
	"encoding/json"
	"net/http"

	s "github.com/cloyop/veetro/internal/server"
	"github.com/cloyop/veetro/internal/storage"
)

func login(c *s.CustomContext) error {
	sreq := &storage.User{}
	if err := json.NewDecoder(c.Request.Body).Decode(&sreq); err != nil {
		return s.ResponseBadJSON(c.Writer, "", nil)
	}
	u, exist, err := c.Storage.GetUserByEmail(sreq.Email)
	if err != nil {
		return err
	}
	if !exist || u.CheckPassword(sreq.Password) != nil {
		return s.ResponseBadJSON(c.Writer, "invalid credentials", nil)
	}
	ses := c.Sessions.NewSession(u)
	u.Password = ""
	return s.ResponseJSON(c.Writer, http.StatusOK, &map[string]interface{}{"user": u, "token": ses.SessionId, "code": 200})
}
func logout(c *s.CustomContext) error {
	if !c.Sessions.RemoveSession(c.Session.SessionId) {
		return s.ResponseBadJSON(c.Writer, "cant found session", nil)
	}
	return s.ResponseJSON(c.Writer, http.StatusOK, &s.ResponseStatus{Code: 200})
}
func signup(c *s.CustomContext) error {
	sreq := &storage.User{}
	if err := json.NewDecoder(c.Request.Body).Decode(&sreq); err != nil {
		return s.ResponseBadJSON(c.Writer, "", nil)
	}
	if errs := sreq.Validate(); len(errs) != 0 {
		return s.ResponseBadJSON(c.Writer, "", &errs)
	}
	exist, err := c.Storage.UserExist(sreq.Email)
	if err != nil {
		return err
	}
	if exist {
		return s.ResponseBadJSON(c.Writer, "user exists", nil)
	}
	if sreq.HashPassword() != nil {
		return err
	}
	u := storage.NewUser(sreq.Name, sreq.Email, sreq.Password, sreq.Location, sreq.Employee, sreq.Customer)
	if err := c.Storage.CreateUser(u); err != nil {
		return err
	}
	ses := c.Sessions.NewSession(u)
	u.Password = ""
	return s.ResponseJSON(c.Writer, http.StatusOK, &map[string]interface{}{"user": u, "token": ses.SessionId, "code": 200})
}

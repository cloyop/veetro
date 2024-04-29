package handlers

import (
	"encoding/json"

	s "github.com/cloyop/veetro/internal/server"
	"github.com/cloyop/veetro/internal/storage"
)

func LoadHandlers(srv *s.Server) {

	srv.Handle("/user/up", signup, "POST")
	srv.Handle("/user/in", login, "POST")
	srv.Handle("/user/out", logout, "POST", s.Auth)
	// Employee
	srv.Handle("/user/application/{app_id}", employeeApplication, "GET PUT DELETE", s.Auth, s.EmployeeOnly)
	srv.Handle("/user/application", employeeApplications, "GET POST DELETE", s.Auth, s.EmployeeOnly)
	// Customer
	srv.Handle("/user/offer/{offer_id}/application/{app_id}", customerOfferApplication, "GET DELETE", s.Auth, s.CustomerOnly) //✔
	srv.Handle("/user/offer/{offer_id}/application", customerOfferApplications, "GET", s.Auth, s.CustomerOnly)
	srv.Handle("/user/offer/{offer_id}", customerOffer, "GET PUT DELETE", s.Auth, s.CustomerOnly)
	srv.Handle("/user/offer", customerOffers, "GET POST DELETE", s.Auth, s.CustomerOnly)
	srv.Handle("/user", user, "GET PUT DELETE", s.Auth)
	// Offers
	srv.Handle("/offers/{offer_id}", offerHandler, "GET")
	srv.Handle("/offers", offersHandler, "GET")
}
func user(c *s.CustomContext) error {
	switch c.Request.Method {
	case "GET":
		u := c.Session.User
		u.Password = ""
		if c.Session.Employee {
			aps, err := c.Storage.GetAllApplications("", u.Id)
			if err != nil {
				return err
			}
			return s.ResponseJSON(c.Writer, "success", &map[string]interface{}{"user": u, "applications": aps})
		}
		ofs, err := c.Storage.GetAllOffersWithApplications(u.Id)
		if err != nil {
			return err
		}
		return s.ResponseJSON(c.Writer, "success", &map[string]interface{}{"user": u, "offers": ofs})
	case "PUT":
		u := &storage.User{}
		if err := json.NewDecoder(c.Request.Body).Decode(u); err != nil {
			return s.ResponseBadJSON(c.Writer, "", nil)
		}
		edits, errs := u.ParseUpdateFields(&c.Session.User)
		if len(*errs) > 0 {
			return s.ResponseBadJSON(c.Writer, "", errs)
		}
		if len(*edits) == 0 {
			return s.ResponseBadJSON(c.Writer, "missing fields", nil)
		}
		is, err := c.Storage.UpdateUser(c.Session.Id, edits)
		if err != nil {
			return err
		}
		if !is {
			return s.ResponseBadJSON(c.Writer, "nothing new to update", nil)
		}
		u, _, _ = c.Storage.GetUserByEmail(c.Session.Email)
		c.Session.User = *u
		u.Password = ""
		return s.ResponseJSON(c.Writer, "user "+u.Id+" updated", u)
	case "DELETE":
		var success bool
		if c.Session.Customer {
			sc, err := c.Storage.DeleteCustomer(c.Session.Id)
			if err != nil {
				return err
			}
			success = sc
		} else {
			sc, err := c.Storage.DeleteEmployee(c.Session.Id)
			if err != nil {
				return err
			}
			success = sc
		}
		if !success {
			return s.ResponseBadJSON(c.Writer, "user not delete", nil)
		}
		c.Sessions.RemoveSession(c.Session.SessionId)
		return s.ResponseJSON(c.Writer, "user "+c.Session.Id+" deleted", nil)
	}
	return nil
}

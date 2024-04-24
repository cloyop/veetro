package handlers

import (
	"encoding/json"
	"fmt"

	s "github.com/cloyop/veetro/internal/server"
	"github.com/cloyop/veetro/internal/storage"
)

func employeeApplications(c *s.CustomContext) error {
	switch c.Request.Method {
	case "GET":
		offs, err := c.Storage.GetAllApplications("", c.Session.Id)
		if err != nil {
			return err
		}
		return s.ResponseJSON(c.Writer, 200, offs)
	case "POST":
		appl := &storage.Application{}
		if err := json.NewDecoder(c.Request.Body).Decode(appl); err != nil {
			return s.ResponseBadJSON(c.Writer, "", nil)
		}
		is, err := c.Storage.OfferExist(appl.OfferId)
		if err != nil {
			return err
		}
		if !is {
			return s.ResponseBadJSON(c.Writer, "Offer not exist", nil)
		}
		is, err = c.Storage.ApplicationExist("", appl.OfferId, c.Session.Id)
		if err != nil {
			return err
		}
		if is {
			return s.ResponseBadJSON(c.Writer, "Already applied to this offer", nil)
		}
		if appl.Resume == "" {
			if !c.Session.ResumeUploaded || c.Session.Resume == "" {
				return s.ResponseBadJSON(c.Writer, "User has not resume to send", nil)
			}
			appl.Resume = c.Session.Resume
		}
		if appl.EmployeeLocation == "" {
			appl.EmployeeLocation = c.Session.Location
		}
		appl = storage.NewApplication(appl.OfferId, c.Session.Id, appl.Resume, appl.CoverLetter, appl.EmployeeLocation, &appl.Anwsers)
		if err := c.Storage.CreateApplication(appl); err != nil {
			return err
		}
		return s.ResponseJSON(c.Writer, 200, appl)
	case "DELETE":
		ids := &[]string{}
		if err := json.NewDecoder(c.Request.Body).Decode(ids); err != nil {
			return s.ResponseBadJSON(c.Writer, "", nil)
		}
		if len(*ids) < 1 {
			return s.ResponseBadJSON(c.Writer, "Missing applications ids", nil)
		}
		dc, err := c.Storage.DeleteApplicationsByIds(ids, c.Session.Id)
		if err != nil {
			return err
		}
		if dc == 0 {
			return s.ResponseBadJSON(c.Writer, "nothing to delete", nil)
		}
		return s.ResponseJSON(c.Writer, 200, &s.ResponseStatus{Code: 200, Message: fmt.Sprintf("%d applications deleted", dc)})
	}
	return nil
}
func employeeApplication(c *s.CustomContext) error {
	app_id := c.Request.PathValue("app_id")
	switch c.Request.Method {
	case "GET":
		appl, is, err := c.Storage.GetApplication(app_id, "", c.Session.Id)
		if err != nil {
			return err
		}
		if !is {
			return s.ResponseBadJSON(c.Writer, "application not exists", nil)
		}
		return s.ResponseJSON(c.Writer, 200, appl)
	case "PUT":
		a := &storage.Application{}
		if err := json.NewDecoder(c.Request.Body).Decode(a); err != nil {
			return s.ResponseBadJSON(c.Writer, "", nil)
		}
		updts, errs := a.ParseUpdateFields()
		if len(*errs) > 0 {
			return s.ResponseBadJSON(c.Writer, "", errs)
		}
		if len(*updts) == 0 {
			return s.ResponseBadJSON(c.Writer, "missing fields", nil)
		}
		found, upt, err := c.Storage.UpdateApplication(app_id, c.Session.Id, updts)
		if err != nil {
			return err
		}
		if !found {
			return s.ResponseBadJSON(c.Writer, "application not found", nil)
		}
		if !upt {
			return s.ResponseBadJSON(c.Writer, "nothing new to update", nil)
		}
		a, _, err = c.Storage.GetApplication(app_id, "", c.Session.Id)
		if err != nil {
			return err
		}
		return s.ResponseJSON(c.Writer, 200, a)
	case "DELETE":
		is, err := c.Storage.DeleteApplication(app_id, "", c.Session.Id)
		if err != nil {
			return err
		}
		if !is {
			return s.ResponseBadJSON(c.Writer, "nothing deleted", nil)
		}
		c.Writer.WriteHeader(200)
		return nil
	}
	return nil
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	s "github.com/cloyop/veetro/internal/server"
	"github.com/cloyop/veetro/internal/storage"
)

func customerOffers(c *s.CustomContext) error {
	switch c.Request.Method {
	case "GET":
		offers, err := c.Storage.GetAllOffersWithApplications(c.Session.Id)
		if err != nil {
			return err
		}
		return s.ResponseJSON(c.Writer, "success", offers)
	case "POST":
		offerR := &storage.Offer{}
		if err := json.NewDecoder(c.Request.Body).Decode(&offerR); err != nil {
			return s.ResponseBadJSON(c.Writer, "", nil)
		}
		if errs := offerR.Validate(); len(errs) != 0 {
			return s.ResponseBadJSON(c.Writer, "", &errs)
		}
		offerR = storage.NewOffer(c.Session.Id, offerR.Title, offerR.Role, offerR.Description, offerR.Location, offerR.Open, &offerR.Questions)
		if err := c.Storage.CreateOffer(offerR); err != nil {
			return err
		}
		if offerR.Open {
			c.State.Change(true)
		}
		return s.ResponseJSON(c.Writer, "Offer created", offerR)
	case "DELETE":
		offersIds := &[]string{}
		if err := json.NewDecoder(c.Request.Body).Decode(offersIds); err != nil {
			return s.ResponseBadJSON(c.Writer, "", nil)
		}
		if len(*offersIds) < 1 {
			return s.ResponseBadJSON(c.Writer, "missing offers ids", nil)
		}
		quantity, err := c.Storage.DeleteOffersByIds(offersIds, c.Session.Id)
		if err != nil {
			return nil
		}
		if quantity == 0 {
			return s.ResponseBadJSON(c.Writer, "nothing to delete", nil)
		}
		c.State.Change(true)
		return s.ResponseBadJSON(c.Writer, fmt.Sprintf("%d offers deleted", quantity), nil)
	}
	return nil
}
func customerOffer(c *s.CustomContext) error {
	offerId := c.Request.PathValue("offer_id")
	switch c.Request.Method {
	case "GET":
		offer, is, err := c.Storage.GetOfferWithApplications(offerId, c.Session.Id)
		if err != nil {
			return err
		}
		if !is {
			return s.ResponseBadJSON(c.Writer, "offer not Found", nil)
		}
		return s.ResponseJSON(c.Writer, "success", offer)
	case "PUT":
		o := &storage.Offer{}
		if err := json.NewDecoder(c.Request.Body).Decode(o); err != nil {
			return s.ResponseBadJSON(c.Writer, "", nil)
		}
		filters, errs := o.ParseUpdateFields()
		if len(*errs) > 0 {
			return s.ResponseBadJSON(c.Writer, "", errs)
		}
		var ChangeOpen bool
		if isOpen := c.Request.URL.Query().Get("open"); isOpen == "true" || isOpen == "false" {
			if isOpen == "true" {
				filters["open"] = true
			} else {
				filters["open"] = false
			}
			ChangeOpen = true
		}
		if len(filters) == 0 {
			return s.ResponseBadJSON(c.Writer, "missing fields", nil)
		}
		found, upt, err := c.Storage.UpdateOffer(offerId, c.Session.Id, &filters)
		if err != nil {
			return err
		}
		if !found {
			return s.ResponseBadJSON(c.Writer, "offer not found", nil)
		}
		if !upt {
			return s.ResponseBadJSON(c.Writer, "nothing new to update", nil)
		}
		o, _, err = c.Storage.GetOfferWithApplications(offerId, c.Session.Id)
		if err != nil {
			return err
		}
		if ChangeOpen {
			c.State.Change(true)
		}
		return s.ResponseJSON(c.Writer, "offer "+o.Id+" updated", o)
	case "DELETE":
		sucess, err := c.Storage.DeleteOffer(offerId, c.Session.Id)
		if err != nil {
			return err
		}
		if !sucess {
			return s.ResponseBadJSON(c.Writer, "offer not found", nil)
		}
		c.State.Change(true)
		return s.ResponseJSON(c.Writer, "Offer Deleted", nil)
	}
	return nil
}
func customerOfferApplications(c *s.CustomContext) error {
	offerId := c.Request.PathValue("offer_id")
	appl, err := c.Storage.GetAllApplications(offerId, "")
	if err != nil {
		return err
	}
	return s.ResponseJSON(c.Writer, "success", appl)
}
func customerOfferApplication(c *s.CustomContext) error {
	offerId := c.Request.PathValue("offer_id")
	ApplyId := c.Request.PathValue("app_id")
	switch c.Request.Method {
	case "GET":
		appl, is, err := c.Storage.GetApplication(ApplyId, offerId, "")
		if err != nil {
			return err
		}
		if !is {
			return s.ResponseBadJSON(c.Writer, "application not found", nil)
		}
		return s.ResponseJSON(c.Writer, "success", appl)
	case "DELETE":
		success, err := c.Storage.DeleteApplication(ApplyId, offerId, "")
		if err != nil {
			return err
		}
		if !success {
			return s.ResponseBadJSON(c.Writer, "couldnt delete application", nil)
		}
		c.Writer.WriteHeader(http.StatusOK)
		return nil
	}
	return nil
}

package handlers

import (
	"net/http"

	"github.com/cloyop/veetro/internal/server"
)

func offersHandler(c *server.CustomContext) error {
	q := c.Request.URL.Query()
	fsize := len(q)
	if fsize == 0 && !c.State.HasChanged() {
		return server.ResponseJSON(c.Writer, http.StatusOK, &map[string]interface{}{"offers": c.State.CurrentOffers(), "totalOffers": c.State.OpenOffers()})
	}
	offers, err := c.Storage.GetAllOffers(q.Get("keyword"), q.Get("role"), q.Get("location"))
	if err != nil {
		return err
	}
	if fsize == 0 && c.State.HasChanged() {
		c.State.UpdateOpenOffers(offers)
	}
	return server.ResponseJSON(c.Writer, http.StatusOK, &map[string]interface{}{"offers": offers, "totalOffers": c.State.OpenOffers()})
}
func offerHandler(c *server.CustomContext) error {
	Offerid := c.Request.PathValue("offer_id")
	o, exist, err := c.Storage.GetOffer(Offerid)
	if err != nil {
		return err
	}
	if !exist {
		return server.ResponseBadJSON(c.Writer, "offer Dont Exists", nil)
	}
	return server.ResponseJSON(c.Writer, http.StatusOK, o)
}

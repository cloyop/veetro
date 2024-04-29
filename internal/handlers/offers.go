package handlers

import (
	"strconv"

	"github.com/cloyop/veetro/internal/server"
)

func offersHandler(c *server.CustomContext) error {
	q := c.Request.URL.Query()
	page := 1
	if pageStr := q.Get("page"); pageStr != "" {
		pint, err := strconv.Atoi(pageStr)
		if err != nil || pint < 1 {
			return server.ResponseBadJSON(c.Writer, "Invalid Pagination Param", nil)
		}
		page = pint
	}
	offers, n, err := c.Storage.GetAllOffers(q.Get("keyword"), q.Get("role"), q.Get("location"), page)
	if err != nil {
		return err
	}
	return server.ResponseJSON(c.Writer, "success", &map[string]interface{}{"offers": offers, "totalOffers": n})
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
	return server.ResponseJSON(c.Writer, "success", o)
}

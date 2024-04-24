package storage

import (
	"fmt"
	"time"
)

// MarkDown string format desc
type Offer struct {
	Id           string        `json:"id"`
	Title        string        `json:"title"`
	Role         string        `json:"role"`
	Description  string        `json:"description"`
	OwnerId      string        `json:"ownerId" bson:"owner_id"`
	Open         bool          `json:"open"`
	Location     string        `json:"location"`
	Created      int64         `json:"created"`
	Questions    []string      `json:"questions"`
	Applications []Application `json:"applications" bson:"omitempty"`
}

func NewOffer(ownerId, title, role, description, location string, open bool, questions *[]string) *Offer {
	now := time.Now().Unix()
	if questions == nil {
		questions = &[]string{}
	}
	return &Offer{
		Id:          fmt.Sprintf("%s%d%s", ownerId[:8], now, ownerId[len(ownerId)-4:]),
		Title:       title,
		Description: description,
		Role:        role,
		Open:        open,
		OwnerId:     ownerId,
		Location:    location,
		Created:     now,
		Questions:   *questions,
	}
}

func (or *Offer) Validate() (errs []string) {
	if !or.validTitle() {
		errs = append(errs, "Invalid Title")
	}
	if !or.validRole() {
		errs = append(errs, "Invalid Role")
	}
	if !or.validDescription() {
		errs = append(errs, "Invalid Description")
	}
	if or.Location == "" {
		or.Location = "remote"
	}
	if !validLocation(or.Location) {
		errs = append(errs, "Invalid Location")
	}
	return
}
func (o *Offer) ParseUpdateFields() (map[string]any, *[]string) {
	filters := map[string]any{}
	errs := []string{}
	if o.Title != "" {
		if o.validTitle() {
			filters["title"] = o.Title
		} else {
			errs = append(errs, fmt.Sprintf(`Invalid Title %v`, o.Title))
		}
	}
	if o.Role != "" {
		if o.validRole() {
			filters["role"] = o.Role
		} else {
			errs = append(errs, fmt.Sprintf(`Invalid Role %v`, o.Role))
		}
	}
	if o.Description != "" {
		if o.validDescription() {
			filters["description"] = o.Description
		} else {
			errs = append(errs, fmt.Sprintf(`Invalid Description %v`, o.Description))
		}
	}
	if o.Location != "" {
		if validLocation(o.Location) {
			filters["location"] = o.Location
		} else {
			errs = append(errs, fmt.Sprintf(`Invalid location %v`, o.Description))
		}
	}
	return filters, &errs
}
func (o *Offer) validTitle() bool {
	return len(o.Title) > 5 && len(o.Title) < 50
}
func (o *Offer) validDescription() bool {
	return len(o.Description) > 100
}
func (o *Offer) validRole() bool {
	return len(o.Role) > 5 && len(o.Role) < 100
}
func validLocation(loc string) bool {
	return len(loc) > 5 && len(loc) < 30
}

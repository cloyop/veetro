package server

import (
	"sync"

	"github.com/cloyop/veetro/internal/storage"
)

type state struct {
	mu          sync.Mutex
	offersCount int
	changed     bool
	offers      *[]storage.Offer
}

func (s *state) UpdateOpenOffers(o *[]storage.Offer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offers = o
	s.offersCount = len(*o)
	s.Change(false)
}
func (s *state) HasChanged() bool {
	return s.changed
}
func (s *state) Change(b bool) {
	s.changed = b
}
func (s *state) CurrentOffers() *[]storage.Offer {
	return s.offers
}
func (s *state) OpenOffers() int {
	return s.offersCount
}

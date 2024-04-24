package server

import (
	"sync"
	"time"

	"github.com/cloyop/veetro/internal/storage"
	"github.com/google/uuid"
)

type session struct {
	SessionId string
	Valid     bool
	CreatedAt int64
	storage.User
}
type sessions struct {
	mu       sync.Locker
	sessions []session
}

func (ss *sessions) NewSession(u *storage.User) *session {
	id := uuid.New().String()
	s := &session{
		SessionId: id,
		User:      *u,
		Valid:     true,
		CreatedAt: time.Now().Unix(),
	}
	ss.mu.Lock()
	ss.sessions = append(ss.sessions, *s)
	ss.mu.Unlock()
	return s
}
func (ss *sessions) GetSession(id string) (*session, bool) {
	for _, s := range ss.sessions {
		if s.SessionId == id {
			return &s, true
		}
	}
	return nil, false
}
func (ss *sessions) RemoveSession(id string) bool {
	for i, s := range ss.sessions {
		if s.SessionId == id {
			ss.mu.Lock()
			ns := []session{}
			ns = append(ns, ss.sessions[:i]...)
			ns = append(ns, ss.sessions[i+1:]...)
			ss.sessions = ns
			ss.mu.Unlock()
			return true
		}
	}
	return false
}

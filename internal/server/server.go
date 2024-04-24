package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cloyop/veetro/internal/storage"
)

type Server struct {
	State      *state
	listenAddr string
	Sessions   *sessions
	router     *http.ServeMux
	Storage    storage.StorageService
}

func New(ln string, s storage.StorageService) *Server {
	return &Server{
		router:     http.NewServeMux(),
		Sessions:   &sessions{},
		listenAddr: ln,
		Storage:    s,
		State:      &state{changed: true},
	}
}
func (s *Server) Run() error {
	s.Storage.Init()
	defer s.Storage.Close()
	fmt.Printf("Running on %v\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, corsMid(s.router))
}
func (s *Server) Handle(path string, handler CustomHandlerFunc, methods string, mws ...CustomMW) {
	s.router.HandleFunc(path, s.h(handler, methods, mws...))
}
func (s *Server) h(f CustomHandlerFunc, methods string, mws ...CustomMW) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(methods, r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ctx := &CustomContext{Writer: w, Request: r, Server: s}
		for _, mw := range mws {
			if success, res := mw(ctx); !success {
				ResponseJSON(w, res.Code, res)
				return
			}
		}

		if err := f(ctx); err != nil {
			log.Printf("Error: %v -> %v - %v\n", err, r.URL, r.Method)
			ResponseErrInternalSrv(w)
		}
	}
}

func ResponseJSON(w http.ResponseWriter, status int, v any) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	return err
}
func ResponseBadJSON(w http.ResponseWriter, e string, errs *[]string) error {
	r := &ResponseStatus{Error: "Bad Request", Code: 400}
	if errs != nil && len(*errs) > 0 {
		r.Errs = *errs
	}
	if e != "" {
		r.Error = e
	}
	res, err := json.Marshal(r)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(400)
	_, err = w.Write(res)
	return err
}
func ResponseErrInternalSrv(w http.ResponseWriter) {
	w.WriteHeader(500)
	if err := json.NewEncoder(w).Encode(&ResponseStatus{Code: 500, Error: "Internal Error"}); err != nil {
		fmt.Println(err)
	}
}

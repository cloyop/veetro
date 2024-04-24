package main

import (
	"flag"

	"github.com/cloyop/veetro/internal/handlers"
	s "github.com/cloyop/veetro/internal/server"
	"github.com/cloyop/veetro/pkg/mongo"
)

func main() {
	port := flag.String("port", ":8080", "http port")
	flag.Parse()
	srv := s.New(*port, mongo.New())
	handlers.LoadHandlers(srv)
	srv.Run()

}

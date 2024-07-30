package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/sean9999/good-graph/api"
	society "github.com/sean9999/good-graph/graph"
)

func main() {

	//	graph
	bus := society.NewBus()
	db := society.NewJsonStore("testdata")
	g, _ := society.NewGraph(db, bus, rand.Reader)
	// society, err := g.Db.Load()
	// if err != nil {
	// 	panic(err)
	// }

	// router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//	REST api
	r.Mount("/api", api.Routes(g))

	//	websockets
	r.Mount("/ws", g.Broker)

	// HTML
	r.Mount("/", http.FileServer(http.Dir("./www/dist")))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8282"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Info().Str("addr", addr).Msg("Starting server")
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 7 * time.Second,
		IdleTimeout:  43 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Error().Msgf("%v", err)
	}

}

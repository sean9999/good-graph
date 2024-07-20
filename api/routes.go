package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sean9999/good-graph/db"
	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/good-graph/ws"
)

func Routes(db db.Database, graph graph.Society, msgs chan ws.Msg) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middlewareJson)

	router.Get("/vertices", GetVertices(db, graph))
	router.Get("/vertex/{nick}", GetVertexByNick(db, graph))
	router.Get("/neighbours/{nick}", GetNeighbours(db, graph))
	router.Get("/edge/{from}/{to}", GetEdge(db, graph))

	router.Post("/vertices", AddVertex(db, graph, msgs))
	//apiRouter.Get("/edges", ListEdges)
	//apiRouter.Post("/befriend/{nick1}/{nick2}", api.CreateEdge)
	router.Post("/befriend", Befriend(db, graph, msgs))
	return router
}

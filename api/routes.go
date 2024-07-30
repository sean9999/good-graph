package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sean9999/good-graph/graph"
)

func Routes(g graph.Graph) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middlewareJson)

	router.Get("/vertices", GetVertices(g))
	router.Get("/vertex/{nick}", GetVertexByNick(g))
	// router.Get("/neighbours/{nick}", GetNeighbours(g))
	// router.Get("/edge/{from}/{to}", GetEdge(g))

	// router.Post("/vertices", AddVertex(g))
	// //apiRouter.Get("/edges", ListEdges)
	// //apiRouter.Post("/befriend/{nick1}/{nick2}", api.CreateEdge)
	// router.Post("/befriend", Befriend(g))
	return router
}

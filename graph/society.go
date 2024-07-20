package graph

import (
	"fmt"

	"github.com/hmdsefi/gograph"
)

// a Society is a Graph of Citizens
type Society struct {
	gograph.Graph[Peer]
}

func (soc Society) AddPeer(p Peer) *gograph.Vertex[Peer] {
	return soc.AddVertexByLabel(p)
}

func (soc Society) RemovePeer(p Peer) {
	vert := soc.GetVertexByID(p)
	soc.RemoveVertices(vert)
}

func (soc Society) RelationshipExists(rel Relationship) bool {
	v1 := gograph.NewVertex(rel.From)
	v2 := gograph.NewVertex(rel.To)
	return soc.ContainsEdge(v1, v2)
}

// func (soc Society) GetRelationships() []Citizen {
// 	return soc.GetAllEdges(from *gograph.Vertex[Peer], to *gograph.Vertex[Peer])
// }

func (soc Society) GetNeighbours(p Peer) []Citizen {
	return soc.Graph.GetVertexByID(p).Neighbors()
}

func (soc Society) Follow(p1 Peer, p2 Peer) error {
	v1 := soc.GetVertexByID(p1)
	v2 := soc.GetVertexByID(p2)
	if v1 == nil || v2 == nil {
		return ErrNoPeer
	}
	_, err := soc.AddEdge(v1, v2)
	return err
}

func (soc Society) Befriend(p1 Peer, p2 Peer) error {
	err1 := soc.Follow(p1, p2)
	err2 := soc.Follow(p2, p1)
	if err1 != nil && err2 != nil {
		return fmt.Errorf("%w - %w", err1, err2)
	}
	return nil
}

func (soc Society) GetVertexByNick(nick string) Citizen {
	for _, vert := range soc.GetAllVertices() {
		if vert.Label().Nickname() == nick {
			return vert
		}
	}
	return nil
}

func NewSociety() Society {
	g := gograph.New[Peer](gograph.Directed())
	return Society{g}
}

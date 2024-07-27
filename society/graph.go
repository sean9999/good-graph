package society

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/hmdsefi/gograph"
)

var ErrNoPeer = errors.New("no such peer")

// a Graph is a [Society] represented as a graph data-structure
// with the ability to persist itself to a [Database]
// and to communicate with the [OutsideWorld].
type Graph struct {
	engine     gograph.Graph[Peer]
	Db         Database
	Broker     OutsideWorld
	Randomness io.Reader
	Society    Society
}

var NoGraph Graph

func NewGraph(db Database, broker OutsideWorld, randy io.Reader) (Graph, error) {

	engine := gograph.New[Peer](gograph.Directed())
	graph := Graph{
		engine,
		db,
		broker,
		randy,
		nil,
	}
	go graph.advertise("info/graphInitialized", nil, 0)

	soc, err := graph.Db.Load()
	if err != nil {
		return NoGraph, err
	}
	graph.Society = soc
	go graph.advertise("info/graphLoaded", nil, 0)

	go func() {
		for ev := range broker.Listen() {
			switch ev.Title {
			case "please/addNode":
				graph.GeneratePeer()
				// graph.AddPeer(p)
				// j := p.Json()

				// msg := Message{
				// 	Title: "nodeAdded",
				// 	Body:  j,
				// 	ID:    p.ToInt(),
				// }
				// bus.Outbox() <- msg
			default:
				fmt.Printf("unknown message title: inbox: %v\n", ev)
			}
		}
	}()

	return graph, nil
}

func (g Graph) GeneratePeer() Peer {
	p := NewPeer(g.Randomness)
	g.AddPeer(p)
	return p
}

func (g Graph) AddPeer(p Peer) error {
	err := g.Db.AddPeer(p)
	if err != nil {
		return err
	}
	err = g.Society.AddPeer(p)
	if err != nil {
		return err
	}
	vert := gograph.NewVertex(p)
	g.engine.AddVertex(vert)
	go g.advertise("note/peerAdded", p, p.ToInt())
	return nil
}

func (g Graph) Peer(hash string) (Peer, error) {
	return g.Society.Peer(hash)
}

func (g Graph) Peers() (map[string]Peer, error) {
	return g.Peers()
}

func (g Graph) advertise(typ string, msg json.Unmarshaler, n int) {
	j, _ := json.Marshal(msg)
	m := Message{
		Title: fmt.Sprintf("society/%s", typ),
		Body:  json.RawMessage(j),
		ID:    n,
	}
	g.Broker.Outbox() <- m
}

func (g Graph) RemovePeer(p Peer) {
	vert := g.engine.GetVertexByID(p)
	if vert != nil {
		g.engine.RemoveVertices(vert)
		g.Db.RemovePeer(p)
		go g.advertise("note/peerRemoved", p, 1)
	}
}

func (pol Graph) RelationshipExists(rel Relationship) bool {
	v1 := gograph.NewVertex(rel.From)
	v2 := gograph.NewVertex(rel.To)
	return pol.engine.ContainsEdge(v1, v2)
}

func (g Graph) GetNeighbours(p Peer) []Peer {
	ps := []Peer{}
	for _, vertex := range g.engine.GetVertexByID(p).Neighbors() {
		ps = append(ps, vertex.Label())
	}
	return ps
}

func (g Graph) Follow(p1 Peer, p2 Peer) error {
	v1 := g.engine.GetVertexByID(p1)
	v2 := g.engine.GetVertexByID(p2)
	if v1 == nil || v2 == nil {
		return ErrNoPeer
	}
	rel := NewRelationship(p1, p2)
	_, err := g.engine.AddEdge(v1, v2)
	if err == nil {
		g.Db.AddRelationship(rel)
		go g.advertise("addFollow", &rel, 2)
	}
	return err
}

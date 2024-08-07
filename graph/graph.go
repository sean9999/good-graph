package graph

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hmdsefi/gograph"
)

var ErrNoPeer = errors.New("no such peer")

// a Graph is a set of [Peer]s, connected via [Relationship]s
// having the ability to persist itself to a [Database]
// and to send and receive [Message]s via a [Broker].
type Graph struct {
	engine     gograph.Graph[Peer]
	Db         Database
	Broker     Broker
	Randomness io.Reader
}

var NoGraph Graph

func loadGraph(g Graph) {

	time.Sleep(5 * time.Second)

	//	load data from database
	snap, _ := g.Db.Load()
	for _, peer := range snap.Peers {
		g.AddPeer(peer)
	}
	for _, rel := range snap.Relationships {
		g.AddRelationship(rel)
	}
	go g.advertise(Message{Subject: "graph/loaded"})
}

func NewGraph(db Database, broker Broker, randy io.Reader) (Graph, error) {

	//	initialize
	engine := gograph.New[Peer](gograph.Directed())
	graph := Graph{
		engine,
		db,
		broker,
		randy,
	}
	go graph.advertise(Message{
		Subject: "graph/initialized",
	})

	go loadGraph(graph)

	//	listen for messages
	go func() {
		for ev := range broker.Listen() {
			switch ev.Subject {
			case "please/addNode":
				fmt.Println(ev)
				graph.GeneratePeer()
			case "please/addRelationship":
				graph.AddRelationship(*ev.Relationship)
			case "please/removePeer":
				graph.RemovePeer(*ev.Peer)
			case "please/removeRelationship":
				graph.RemoveRelationship(*ev.Relationship)
			default:
				fmt.Printf("unknown subject %q: %v\n", ev.Subject, ev)
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
	vert := gograph.NewVertex(p)
	if g.engine.ContainsVertex(vert) {
		return errors.New("vertex exists")
	}
	g.engine.AddVertex(vert)
	err := g.Db.AddPeer(p)
	if err != nil {
		g.engine.RemoveVertices(vert)
		return err
	}
	msg := Message{
		Subject: "mutation/peerAdded",
		Peer:    &p,
	}

	fmt.Println("adding pper", msg)

	go g.advertise(msg)
	return nil
}

func (g Graph) Peer(hash string) (Peer, error) {
	for _, vert := range g.engine.GetAllVertices() {
		if hash == vert.Label().Hash() {
			return vert.Label(), nil
		}
	}
	return NoPeer, ErrNoPeer
}

func (g Graph) Peers() map[string]Peer {
	allVerts := g.engine.GetAllVertices()
	m := make(map[string]Peer, len(allVerts))
	for _, vert := range allVerts {
		m[vert.Label().Hash()] = vert.Label()
	}
	return m
}

func (g Graph) advertise(msg Message) {
	if msg.MessageID == 0 {
		msg.MessageID = GlobalID.Add(1)
	}
	if msg.ThreadID == 0 {
		msg.ThreadID = msg.MessageID
	}
	g.Broker.Outbox() <- msg
}

func (g Graph) RemovePeer(p Peer) {
	vert := g.engine.GetVertexByID(p)
	if vert != nil {
		g.engine.RemoveVertices(vert)
		g.Db.RemovePeer(p)
		go g.advertise(Message{
			Subject: "mutation/peerRemoved",
			Peer:    &p,
		})
	}
}

func (g Graph) RemoveRelationship(rel Relationship) {
	v1 := gograph.NewVertex(rel.From)
	v2 := gograph.NewVertex(rel.To)
	e := gograph.NewEdge(v1, v2)
	g.engine.RemoveEdges(e)
}

func (pol Graph) RelationshipExists(rel Relationship) bool {
	v1 := gograph.NewVertex(rel.From)
	v2 := gograph.NewVertex(rel.To)
	return pol.engine.ContainsEdge(v1, v2)
}

func (g Graph) AddRelationship(rel Relationship) error {
	v1 := gograph.NewVertex(rel.From)
	v2 := gograph.NewVertex(rel.To)
	e, err := g.engine.AddEdge(v1, v2)
	if err != nil {
		return err
	}
	err = g.Db.AddRelationship(rel)
	if err != nil {
		g.engine.RemoveEdges(e)
		return err
	}
	go g.advertise(Message{
		Subject:      "mutation/relationshipAdded",
		Relationship: &rel,
	})
	return nil
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
		go g.advertise(Message{
			Subject:      "mutation/addFollow",
			Relationship: &rel,
		})
	}
	return err
}

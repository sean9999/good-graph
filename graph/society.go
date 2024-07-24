package graph

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/hmdsefi/gograph"
	"github.com/sean9999/good-graph/transport"
)

// a Society is a Graph of Citizens
type Society struct {
	gograph.Graph[Peer]
	Inbox  chan transport.Msg
	Outbox chan transport.Msg
}

func (soc Society) AddPeer(p Peer) *gograph.Vertex[Peer] {
	go soc.advertise("peerAdded", p, 1)
	return soc.AddVertexByLabel(p)
}

func (soc Society) advertise(typ, msg any, n int) {
	j, _ := json.Marshal(msg)
	m := transport.Msg{
		MsgType: fmt.Sprintf("society/%s", typ),
		Msg:     json.RawMessage(j),
		N:       n,
	}
	soc.Outbox <- m
}

func (soc Society) RemovePeer(p Peer) {
	vert := soc.GetVertexByID(p)
	if vert != nil {
		soc.RemoveVertices(vert)
		go soc.advertise("removePeer", p.String(), 1)
	}
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

	rel := NewRelationship(p1, p2)
	j, _ := rel.MarshalJson()

	_, err := soc.AddEdge(v1, v2)
	if err == nil {
		go soc.advertise("addFollow", j, 2)
	}
	return err
}

func (soc Society) Befriend(p1 Peer, p2 Peer) error {
	err1 := soc.Follow(p1, p2)
	err2 := soc.Follow(p2, p1)
	if err1 != nil && err2 != nil {
		return fmt.Errorf("%w - %w", err1, err2)
	}
	go soc.advertise("befriend", NewRelationship(p1, p2).String(), 4)
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

func NewSociety(inbox, outbox chan transport.Msg) Society {
	g := gograph.New[Peer](gograph.Directed())
	soc := Society{
		g,
		inbox,
		outbox,
	}
	go soc.advertise("graphCreated", "graphCreated", 0)

	go func() {
		for ev := range inbox {
			fmt.Printf("inbox: %v\n", ev)
			p, err := NewPeer(rand.Reader)
			if err != nil {
				fmt.Println("creating new peer", err)
			}
			soc.AddPeer(p)
			j, err := p.MarshalJSON()
			if err != nil {
				fmt.Println("marshaling peer")
			}
			msg := transport.NewMsg("addThisNode", string(j), p.ToInt())
			outbox <- msg
		}
	}()

	return soc
}

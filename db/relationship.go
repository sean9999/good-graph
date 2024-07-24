package db

import (
	"encoding/json"
	"fmt"

	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/harebrain"
)

type Relationship interface {
	From() Peer
	To() Peer
	harebrain.EncodeHasher
	ToGraph() graph.Relationship
}

var _ Relationship = (*relationship)(nil)

type relationship struct {
	FromPeer Peer `json:"from"`
	ToPeer   Peer `json:"to"`
}

func (rel *relationship) ToGraph() graph.Relationship {
	g := graph.Relationship{
		From: graph.Peer(rel.FromPeer),
		To:   graph.Peer(rel.ToPeer),
	}
	return g
}

func (rel *relationship) From() Peer {
	return rel.FromPeer
}
func (rel *relationship) To() Peer {
	return rel.ToPeer
}

func (rel *relationship) Hash() string {
	return fmt.Sprintf("%s-to-%s.json", rel.From().Nickname(), rel.ToPeer.Nickname())
}
func (rel *relationship) MarshalBinary() ([]byte, error) {
	return rel.MarshalJSON()
}
func (rel *relationship) UnmarshalBinary(p []byte) error {
	return rel.UnmarshalJSON(p)
}

func NewRelationship(from [64]byte, to [64]byte) *relationship {
	p1 := Peer(from)
	p2 := Peer(to)
	return &relationship{p1, p2}
}

// for efficient storage
type thinRel struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (rel *relationship) MarshalJSON() ([]byte, error) {
	tr := thinRel{
		From: rel.FromPeer.String(),
		To:   rel.ToPeer.String(),
	}
	return json.Marshal(tr)
}

func (rel *relationship) UnmarshalJSON(data []byte) error {
	var tr thinRel
	err := json.Unmarshal(data, &tr)
	if err != nil {
		return err
	}
	sourcePeer, err := graph.PeerFromHex([]byte(tr.From))
	if err != nil {
		return err
	}
	destPeer, err := graph.PeerFromHex([]byte(tr.To))
	if err != nil {
		return err
	}
	rel.FromPeer = Peer(sourcePeer)
	rel.ToPeer = Peer(destPeer)
	return nil
}

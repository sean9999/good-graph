package db

import (
	"encoding/json"
	"fmt"

	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/harebrain"
)

type relationship struct {
	FromPeer Peer `json:"from"`
	ToPeer   Peer `json:"to"`
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
func (rel *relationship) Clone() harebrain.EncodeHasher {
	r2 := new(relationship)
	p, _ := rel.MarshalBinary()
	r2.UnmarshalBinary(p)
	return r2
}

func RelationshipFrom(from [64]byte, to [64]byte) *relationship {
	p1 := PeerFrom(from)
	p2 := PeerFrom(to)
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
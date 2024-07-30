package graph

import (
	"encoding/json"
	"fmt"

	"github.com/sean9999/go-oracle"
)

// The Nature of a Relationship is one of Spouse, Friend, or Follows
type Nature uint8

const (
	NoNature Nature = iota
	Spouse
	Friend
	Follows
)

var NoRelationship Relationship

// a Relationship is a directed relationship between two [Peer]s.
//
//	It represents an edge in the graph.
type Relationship struct {
	From   Peer   `json:"from"`
	To     Peer   `json:"to"`
	Nature Nature `json:"nature"`
}

func (rel Relationship) Exists() bool {
	return !rel.Equal(NoRelationship)
}

// Relationship implements Hasher
var _ Hasher = (Relationship)(NoRelationship)

func (r Relationship) Hash() string {
	return fmt.Sprintf("%s:%s:%d", r.From.Hash(), r.To.Hash(), r.Nature)
}
func (r Relationship) Equal(h Hasher) bool {
	return r.Hash() == h.Hash()
}
func (r Relationship) Json() []byte {
	j, _ := json.Marshal(r)
	return j
}
func (r *Relationship) UnmarshalJSON(b []byte) error {
	type ft struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	fff := new(ft)
	err := json.Unmarshal(b, fff)
	if err != nil {
		return err
	}
	f, _ := oracle.PeerFromHex([]byte(fff.From))
	t, _ := oracle.PeerFromHex([]byte(fff.To))
	r.From = Peer(f)
	r.To = Peer(t)
	r.Nature = Follows
	return nil
}

// func (r Relationship) UnmarshalBinary(b []byte) error {
// 	return r.UnmarshalJSON(b)
// }

func (r Relationship) MarshalJSON() ([]byte, error) {
	type fromto struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	h1 := r.From.String()
	h2 := r.To.String()
	x := fromto{
		From: h1,
		To:   h2,
	}
	return json.Marshal(x)
}

func NewRelationship(p1, p2 Peer) Relationship {
	return Relationship{
		From:   p1,
		To:     p2,
		Nature: Follows,
	}
}

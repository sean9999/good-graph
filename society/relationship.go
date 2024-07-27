package society

import (
	"encoding/json"
	"fmt"
)

// a Relationship is a directed relationship between two Citizens.
// The Nature of a Relationship is one of Spouse, Friend, or Follows
type Nature uint8

const (
	NoNature Nature = iota
	Spouse
	Friend
	Follows
)

var NoRelationship Relationship

type Relationship struct {
	From   Peer   `json:"from"`
	To     Peer   `json:"to"`
	Nature Nature `json:"nature"`
}

func (r Relationship) Hash() Hash {
	str := fmt.Sprintf("%s:%s:%d", r.From.Hash(), r.To.Hash(), r.Nature)
	return Hash(str)
}
func (r Relationship) Equal(h Hasher) bool {
	return r.Hash() == h.Hash()
}
func (r Relationship) Json() []byte {
	j, _ := json.Marshal(r)
	return j
}
func (r *Relationship) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, r)
}
func (r Relationship) UnmarshalBinary(b []byte) error {
	return r.UnmarshalJSON(b)
}

func NewRelationship(p1, p2 Peer) Relationship {
	return Relationship{
		From:   p1,
		To:     p2,
		Nature: Follows,
	}
}

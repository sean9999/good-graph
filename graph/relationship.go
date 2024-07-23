package graph

import (
	"encoding/json"
)

type Relationship struct {
	From Peer `json:"from"`
	To   Peer `json:"to"`
}

func (rel Relationship) MarshalJson() ([]byte, error) {
	ps := [2]string{rel.From.String(), rel.To.String()}
	return json.Marshal(ps)
}

func (rel Relationship) String() string {
	j, _ := rel.MarshalJson()
	return string(j)
}

func NewRelationship(p1, p2 Peer) Relationship {
	return Relationship{
		From: p1,
		To:   p2,
	}
}

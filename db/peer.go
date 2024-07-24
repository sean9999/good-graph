package db

import (
	"github.com/sean9999/good-graph/graph"
)

var NoPeer Peer

type Peer graph.Peer

func (p Peer) String() string {
	return p.Graph().String()
}

func (p Peer) Graph() graph.Peer {
	return graph.Peer(p)
}

func (p Peer) Nickname() string {
	return p.Graph().Nickname()
}
func (p Peer) MarshalBinary() ([]byte, error) {
	return p.Graph().MarshalJSON()
}
func (p Peer) Hash() string {
	return p.Nickname() + ".json"
}

func (p Peer) UnmarshalBinary(b []byte) error {
	gp := new(graph.Peer)
	err := gp.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	copy(p[:], gp[:])
	return nil
}

func (p Peer) Bytes() []byte {
	return p[:]
}

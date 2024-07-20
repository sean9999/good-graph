package db

import (
	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/harebrain"
)

type Peer graph.Peer

func (p Peer) String() string {
	return graph.PeerFromBytes(p[:]).String()
}
func (p Peer) Nickname() string {
	return graph.PeerFromBytes(p[:]).Nickname()
}
func (p Peer) Clone() harebrain.EncodeHasher {
	var p2 Peer
	copy(p2[:], p[:])
	return p2
}
func (p Peer) MarshalBinary() ([]byte, error) {
	return graph.PeerFromBytes(p[:]).MarshalJSON()
}
func (p Peer) Hash() string {
	return graph.PeerFromBytes(p[:]).Nickname() + ".json"
}
func (p Peer) UnmarshalBinary(b []byte) error {
	copy(p[:], b)
	return nil
}

func PeerFrom(b [64]byte) Peer {
	var p Peer
	copy(p[:], b[:])
	return p
}

package db

import (
	"encoding/json"

	"github.com/sean9999/good-graph/graph"
)

var NoPeer Peer

type Peer graph.Peer

func (p Peer) String() string {
	return p.Graph().String()
}

func (p Peer) Graph() graph.Peer {
	return graph.PeerFrom(p)
}

func (p Peer) Nickname() string {
	return p.Graph().Nickname()
}
func (p Peer) MarshalBinary() ([]byte, error) {
	return graph.PeerFromBytes(p[:]).MarshalJSON()
}
func (p Peer) Hash() string {
	return graph.PeerFromBytes(p[:]).Nickname() + ".json"
}
func (p Peer) UnmarshalBinary(b []byte) error {
	type smalPeer struct {
		Pubkey string `json:"pubkey"`
	}
	m := smalPeer{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	if m.Pubkey == "" {
		return graph.ErrNoPeer
	}
	p2, err := graph.PeerFromHex([]byte(m.Pubkey))
	if err != nil {
		return err
	}
	copy(p[:], p2[:])
	return nil
}

func (p Peer) Bytes() []byte {
	return p[:]
}

func PeerFrom(b [64]byte) Peer {
	var p Peer
	copy(p[:], b[:])
	return p
}

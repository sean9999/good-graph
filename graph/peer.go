package graph

import (
	"encoding/json"
	"io"
	"math/big"
	"slices"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity"
)

// a Hasher can hash itself to a unique string
// and can distinguish itself from other Hashers
type Hasher interface {
	Hash() string
	Equal(Hasher) bool
}

// a Peer is a Hasher made up of public key material.
// It represents a node in the graph.
type Peer [64]byte

var NoPeer Peer

// a zero-value Peer is said to not exist
func (p Peer) Exists() bool {
	return slices.Equal(p[:], NoPeer[:])
}

func (p Peer) Hash() string {
	return oracle.Peer(p).Nickname()
}
func (p Peer) Equal(h Hasher) bool {
	return p.Hash() == h.Hash()
}
func (p Peer) MarshalJSON() ([]byte, error) {
	return oracle.Peer(p).MarshalJSON()
}

func (p Peer) ToHex() []byte {
	hex, _ := oracle.Peer(p).MarshalHex()
	return hex
}

func (p Peer) String() string {
	return string(p.ToHex())
}

func (p *Peer) UnmarshalJSON(b []byte) error {

	type js struct {
		Address  string `json:"address"`
		Nickname string `json:"nick"`
		Pubkey   string `json:"pubkey"`
	}

	pj := new(js)

	json.Unmarshal(b, pj)

	ph, err := oracle.PeerFromHex([]byte(pj.Pubkey))
	if err != nil {
		return err
	}
	copy(p[:], ph[:])
	return nil
}
func (p Peer) UnmarshalBinary(b []byte) error {
	return p.UnmarshalJSON(b)
}
func (p Peer) MarshalBinary() ([]byte, error) {
	return p.MarshalJSON()
}
func (p Peer) Json() []byte {
	j, _ := p.MarshalJSON()
	return j
}
func (p Peer) ToInt() int {
	fourBytes := make([]byte, 4)
	copy(fourBytes, p[:])
	i := int(big.NewInt(0).SetBytes(fourBytes).Uint64())
	return i
}

func NewPeer(randy io.Reader) Peer {
	p, _ := polity.NewPeer(randy)
	return Peer(p)
}

func PeerFromBytes(b []byte) Peer {
	var p Peer
	copy(p[:], b)
	return p
}

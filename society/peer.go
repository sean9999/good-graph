package society

import (
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
func (p Peer) UnmarshalJSON(b []byte) error {
	o := oracle.Peer(p)
	err := o.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	copy(p[:], o[:])
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

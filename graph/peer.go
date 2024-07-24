package graph

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity"
)

type Peer [64]byte

var NoPeer Peer

func (p Peer) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.AsMap())
}

func (p *Peer) UnmarshalJSON(data []byte) error {
	var m map[string]string
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	pubkey, ok := m["pubkey"]
	if !ok {
		return errors.New("no pubkey")
	}
	ph, err := PeerFromHex([]byte(pubkey))
	copy(p[:], ph[:])
	return err
}

func (p Peer) Nickname() string {
	return polity.Peer(p).Nickname()
}
func (p Peer) ToHex() []byte {
	dst := make([]byte, hex.EncodedLen(64))
	hex.Encode(dst, p[:])
	return dst
}

func (p Peer) ToInt() int {
	fourBytes := make([]byte, 4)
	copy(fourBytes, p[:])
	i := int(big.NewInt(0).SetBytes(fourBytes).Uint64())
	return i
}

func (p Peer) String() string {
	return string(p.ToHex())
}
func (p Peer) AsMap() map[string]string {
	nick := p.Nickname()
	m := map[string]string{
		"pubkey":  p.String(),
		"nick":    nick,
		"address": fmt.Sprintf("mem://%s", nick),
	}
	return m
}
func PeerFromHex(hex []byte) (Peer, error) {
	var p oracle.Peer
	err := p.UnmarshalHex(hex)
	if err != nil {
		return NoPeer, err
	}
	return Peer(p), nil
}
func PeerFromBytes(bs []byte) Peer {
	return Peer(bs)
}
func NewPeer(randy io.Reader) (Peer, error) {
	pol, err := polity.NewPeer(randy)
	return Peer(pol), err
}

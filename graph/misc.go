package graph

import (
	"net"
	"time"

	"github.com/sean9999/polity"
)

type inMemAddress struct {
	nickname string
}

func (a *inMemAddress) Network() string {
	return "inmem"
}
func (a *inMemAddress) String() string {
	return a.nickname
}

type inMemConnection struct {
	Pubkey polity.Peer
}

func (conn *inMemConnection) ReadFrom(p []byte) (int, net.Addr, error) {
	return 0, nil, nil
}
func (conn *inMemConnection) WriteTo(p []byte, addr net.Addr) (int, error) {
	return 0, nil
}
func (conn *inMemConnection) Close() error {
	return nil
}
func (conn *inMemConnection) LocalAddr() net.Addr {
	return nil
}
func (conn *inMemConnection) SetDeadline(t time.Time) error {
	return nil
}
func (conn *inMemConnection) SetReadDeadline(t time.Time) error {
	return nil
}
func (conn *inMemConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
func (conn *inMemConnection) Address() net.Addr {
	return &inMemAddress{
		nickname: conn.Pubkey.Nickname(),
	}
}
func (conn *inMemConnection) Join() error {
	return nil
}
func (conn *inMemConnection) Leave() error {
	return nil
}
func (conn *inMemConnection) AddressFromPubkey(p []byte) net.Addr {
	peer, _ := polity.PeerFromBytes(p)
	return peer.Address(conn)
}
func NewInMemConn(p []byte) *inMemConnection {
	peer, _ := polity.PeerFromBytes(p)
	return &inMemConnection{peer}
}

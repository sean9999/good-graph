package db

import (
	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/harebrain"
)

// type Relationship = graph.Relationship
type Society = graph.Society

type Relationship interface {
	From() Peer
	To() Peer
	harebrain.EncodeHasher
}

type Database interface {
	Open() error
	Close() error
	Load() (Society, error)
	Save(Society) error
	Peers() map[string]Peer
	AddPeer(Peer) error
	RemovePeer(Peer) error
	Relationships() []Relationship
	AddRelationship(Relationship) error
	RemoveRelationship(Relationship) error
}

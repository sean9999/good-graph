package db

import (
	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/good-graph/transport"
)

type Society = graph.Society

type Database interface {
	Open() error
	Close() error
	Load(chan transport.Msg) (Society, error)
	Save(Society) error
	Peers() (map[string]Peer, error)
	AddPeer(Peer) error
	RemovePeer(Peer) error
	Relationships() ([]Relationship, error)
	AddRelationship(Relationship) error
	RemoveRelationship(Relationship) error
}

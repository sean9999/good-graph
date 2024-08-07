package graph

import (
	"errors"
)

var ErrNoRelationship = errors.New("no such relationship")

// Database represents the storage layer for a [Graph]
type Database interface {
	Open() error
	Close() error
	Load() (Snapshot, error)
	Save(Snapshot) error
	Peers() (map[string]Peer, error)
	AddPeer(Peer) error
	RemovePeer(Peer) error
	Relationships() ([]Relationship, error)
	AddRelationship(Relationship) error
	RemoveRelationship(Relationship) error
	GetPeer(string) (Peer, error)
	GetRelationship(string) (Relationship, error)
}

// a [Snapshot] is a snapshot of the entire [Database]
type Snapshot struct {
	Peers         []Peer
	Relationships []Relationship
}

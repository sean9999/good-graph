package db

import (
	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/harebrain"
)

// jsonstore implements Database, using harebrain as back-end
type jsonstore struct {
	Brain *harebrain.Database
	Database
}

func New(rootPath string) *jsonstore {
	brain := harebrain.NewDatabase()
	brain.Open(rootPath)
	store := jsonstore{
		Brain: brain,
	}
	return &store
}

// load an entire graph from database
func (store *jsonstore) Load() (graph.Society, error) {
	society := graph.NewSociety()
	for _, peer := range store.Peers() {
		society.AddPeer(graph.PeerFrom(peer))
	}
	for _, rel := range store.Relationships() {
		society.Befriend(graph.PeerFrom(rel.From()), graph.PeerFrom(rel.To()))
	}
	return society, nil
}

func (store *jsonstore) AddRelationship(rel Relationship) error {
	return store.Brain.Table("relationships").Insert(rel)
}

func (store *jsonstore) AddPeer(p Peer) error {
	return store.Brain.Table("peers").Insert(p)
}

func (store *jsonstore) Peers() map[string]Peer {
	var p Peer
	m := store.Brain.Table("peers").LoadAll(&p)
	m2 := map[string]Peer{}
	for k, v := range m {
		m2[k] = v.(Peer)
	}
	return m2
}

func (store *jsonstore) Relationships() []Relationship {
	var rs []Relationship
	var r relationship
	m := store.Brain.Table("peers").LoadAll(&r)

	for _, v := range m {
		rs = append(rs, v.(*relationship))
	}
	return rs
}

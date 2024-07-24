package db

import (
	"errors"

	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/good-graph/transport"
	"github.com/sean9999/harebrain"
)

var _ Database = (*jsonstore)(nil)

// jsonstore implements Database, using harebrain as back-end
type jsonstore struct {
	Brain *harebrain.Database
	Database
	Inbox  chan transport.Msg
	Outbox chan transport.Msg
}

func New(rootPath string, inbox, outbox chan transport.Msg) *jsonstore {
	brain := harebrain.NewDatabase()
	brain.Open(rootPath)
	store := jsonstore{
		Brain:  brain,
		Inbox:  inbox,
		Outbox: outbox,
	}
	return &store
}

// load an entire graph from database
func (store *jsonstore) Load() (graph.Society, error) {
	society := graph.NewSociety(store.Inbox, store.Outbox)
	peers, err := store.Peers()
	if err != nil {
		return society, err
	}
	rels, err := store.Relationships()
	if err != nil {
		return society, err
	}
	for _, peer := range peers {
		society.AddPeer(graph.Peer(peer))
	}
	for _, rel := range rels {
		society.Befriend(rel.ToGraph().From, rel.ToGraph().To)
	}
	return society, nil
}

// persist an entire graph to database
func (store *jsonstore) Persist(g graph.Society) error {
	return errors.New("not implemented")
}

func (store *jsonstore) AddRelationship(rel Relationship) error {
	return store.Brain.Table("relationships").Insert(rel)
}

func (store *jsonstore) RelationshipExists(rel Relationship) bool {
	_, err := store.Brain.Table("relationships").Get(rel.Hash())
	return (err == nil)
}

func (store *jsonstore) RemoveRelationship(rel Relationship) error {
	return store.Brain.Table("relationships").Delete(rel.Hash())
}

func (store *jsonstore) AddPeer(p Peer) error {
	return store.Brain.Table("peers").Insert(p)
}

func (store *jsonstore) GetPeer(hash string) (Peer, error) {
	b, err := store.Brain.Table("peers").Get(hash)
	if err != nil {
		return NoPeer, err
	}
	p := Peer{}
	err = p.UnmarshalBinary(b)
	if err != nil {
		return NoPeer, err
	}
	return p, nil
}

func (store *jsonstore) RemovePeer(p Peer) error {
	return store.Brain.Table("peers").Delete(p.Hash())
}

func (store *jsonstore) Peers() (map[string]Peer, error) {
	m1, err := store.Brain.Table("peers").GetAll()
	if err != nil {
		return nil, err
	}
	m2 := map[string]Peer{}
	for k, v := range m1 {
		p := new(Peer)
		err := p.UnmarshalBinary(v)
		if err != nil {
			return nil, err
		}
		m2[k] = *p
	}
	return m2, nil
}

func (store *jsonstore) Relationships() ([]Relationship, error) {
	var rs []Relationship
	m1, err := store.Brain.Table("relationships").GetAll()
	if err != nil {
		return nil, err
	}
	for _, v := range m1 {
		rel := new(relationship)
		err := rel.UnmarshalBinary(v)
		if err != nil {
			return nil, err
		}
		rs = append(rs, rel)
	}
	return rs, nil
}

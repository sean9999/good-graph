package graph

import "github.com/sean9999/harebrain"

// *jsonstore implements Database
var _ Database = (*jsonstore)(nil)

// jsonstore is a [Database] that deals in ".json" files on disk
type jsonstore struct {
	Brain *harebrain.Database
}

func NewJsonStore(rootPath string) *jsonstore {
	brain := harebrain.NewDatabase()
	brain.Open(rootPath)
	store := jsonstore{
		Brain: brain,
	}
	return &store
}

func (store *jsonstore) GetPeer(hash string) (Peer, error) {
	b, err := store.Brain.Table("peers").Get(hash)
	if err != nil {
		return NoPeer, ErrNoPeer
	}
	return PeerFromBytes(b), nil
}

func (store *jsonstore) RemovePeer(p Peer) error {
	return store.Brain.Table("peers").Delete(p.Hash())
}
func (store *jsonstore) RemoveRelationship(r Relationship) error {
	return store.Brain.Table("relationships").Delete(r.Hash())
}

func (store *jsonstore) Close() error {
	//	a jsonstore doesn't need to open or close connections
	return nil
}
func (store *jsonstore) Open() error {
	//	a jsonstore doesn't need to open or close connections
	return nil
}

func (store *jsonstore) GetRelationship(hash string) (Relationship, error) {
	b, err := store.Brain.Table("relationships").Get(hash)
	if err != nil {
		return NoRelationship, ErrNoRelationship
	}
	rel := new(Relationship)
	err = rel.UnmarshalJSON(b)
	if err != nil {
		return NoRelationship, err
	}
	return *rel, nil
}

func (store *jsonstore) Peers() (map[string]Peer, error) {
	m1, err := store.Brain.Table("peers").GetAll()
	if err != nil {
		return nil, err
	}
	m2 := map[string]Peer{}
	for k, v := range m1 {
		p := new(Peer)
		err := p.UnmarshalJSON(v)
		if err != nil {
			return nil, err
		}
		m2[k] = *p
	}
	return m2, nil
}

func (store *jsonstore) AddRelationship(rel Relationship) error {
	return store.Brain.Table("relationships").Delete(rel.Hash())
}

func (store *jsonstore) Relationships() ([]Relationship, error) {
	var rs []Relationship
	m1, err := store.Brain.Table("relationships").GetAll()
	if err != nil {
		return nil, err
	}
	for _, v := range m1 {
		rel := new(Relationship)
		err := rel.UnmarshalJSON(v)
		if err != nil {
			return nil, err
		}
		rs = append(rs, *rel)
	}
	return rs, nil
}

func (store *jsonstore) AddPeer(p Peer) error {
	return store.Brain.Table("peers").Insert(p)
}

func (store *jsonstore) Save(snap Snapshot) error {
	for _, peer := range snap.Peers {
		store.AddPeer(peer)
	}
	for _, rel := range snap.Relationships {
		store.AddRelationship(rel)
	}
	return nil
}

func (store *jsonstore) Load() (Snapshot, error) {
	snap := Snapshot{
		Peers:         []Peer{},
		Relationships: []Relationship{},
	}
	peers, _ := store.Peers()
	for _, p := range peers {
		snap.Peers = append(snap.Peers, p)
	}
	rels, _ := store.Relationships()
	for _, rel := range rels {
		snap.Relationships = append(snap.Relationships, rel)
	}
	return snap, nil
}

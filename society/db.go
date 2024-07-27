package society

import "github.com/sean9999/harebrain"

type Database interface {
	Open() error
	Close() error
	Load() (Society, error)
	Save(Society) error
	Peers() (map[string]Peer, error)
	AddPeer(Peer) error
	RemovePeer(Peer) error
	Relationships() ([]Relationship, error)
	AddRelationship(Relationship) error
	RemoveRelationship(Relationship) error
	GetPeer(string) (Peer, error)
	GetRelationship(string) (Relationship, error)
}

type jsonstore struct {
	Brain *harebrain.Database
}

type Dump struct {
	Peers         []Peer
	Relationships []Relationship
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
	return nil
}
func (store *jsonstore) Open() error {
	return nil
}

func (store *jsonstore) GetRelationship(hash string) (Relationship, error) {
	b, err := store.Brain.Table("relationships").Get(hash)
	if err != nil {
		return NoRelationship, ErrNoRelationship
	}
	rel := new(Relationship)
	err = rel.UnmarshalBinary(b)
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
		err := p.UnmarshalBinary(v)
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
		err := rel.UnmarshalBinary(v)
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

func (store *jsonstore) Save(soc Society) error {
	for _, p := range soc.Peers() {
		err := store.AddPeer(p)
		if err != nil {
			return err
		}
	}
	for _, r := range soc.Relationships() {
		err := store.AddRelationship(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *jsonstore) Load() (Society, error) {

	soc := &society{}

	peers, err := store.Peers()
	if err != nil {
		return nil, err
	}
	soc.peers = peers

	rels, err := store.Relationships()
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
		soc.AddPeer(peer)
	}
	for _, rel := range rels {
		soc.AddRelationship(rel.From, rel.To)
	}
	return soc, nil
}

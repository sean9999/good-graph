package graph

// *memStore implements Database
var _ Database = (*memStore)(nil)

// memStore is an in-memory implementation of the Database interface
type memStore struct {
	peers map[string]Peer
	rels  map[string]Relationship
}

func (m *memStore) Open() error {
	return nil
}

func (m *memStore) Close() error {
	return nil
}

func (m *memStore) Load() (Snapshot, error) {
	return Snapshot{
		Peers:         mapToSlice(m.peers),
		Relationships: mapToSlice(m.rels),
	}, nil
}

func (m *memStore) Save(snap Snapshot) error {
	m.peers = sliceToMap(snap.Peers)
	m.rels = sliceToMap(snap.Relationships)
	return nil
}

func (m *memStore) Peers() (map[string]Peer, error) {
	return m.peers, nil
}

func (m *memStore) AddPeer(p Peer) error {
	m.peers[p.Hash()] = p
	return nil
}

func (m *memStore) RemovePeer(p Peer) error {
	delete(m.peers, p.Hash())
	return nil
}

func (m *memStore) Relationships() ([]Relationship, error) {
	return mapToSlice(m.rels), nil
}

func (m *memStore) AddRelationship(rel Relationship) error {
	m.rels[rel.Hash()] = rel
	return nil
}

func (m *memStore) RemoveRelationship(rel Relationship) error {
	delete(m.rels, rel.Hash())
	return nil
}

func (m *memStore) GetPeer(s string) (Peer, error) {
	p, exists := m.peers[s]
	if !exists {
		return p, ErrNoPeer
	}
	return p, nil
}

func (m *memStore) GetRelationship(s string) (Relationship, error) {
	r, exists := m.rels[s]
	if !exists {
		return r, ErrNoPeer
	}
	return r, nil
}

// constructor
func NewMemStore() *memStore {
	return &memStore{
		peers: map[string]Peer{},
		rels:  map[string]Relationship{},
	}
}

// helper functions
func mapToSlice[T Hasher](m map[string]T) []T {
	arr := make([]T, 0, len(m))
	for _, v := range m {
		arr = append(arr, v)
	}
	return arr
}

func sliceToMap[T Hasher](s []T) map[string]T {
	m := make(map[string]T, len(s))
	for _, item := range s {
		m[item.Hash()] = item
	}
	return m
}

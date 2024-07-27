package society

import "errors"

//	society is a collection of Peers and Relationships, working in concert

var ErrNoCitizen = errors.New("no such citizen")
var ErrNoRelationship = errors.New("no such relationship")
var ErrRelationshipExists = errors.New("relationship exists")

// a Society is a set of Peers, and Relationships between them
type Society interface {
	Peers() []Peer
	Relationships() []Relationship
	Peer(Hash) (Peer, error)
	Relationship(Peer, Peer) (Relationship, error)
	AddRelationship(Peer, Peer) error
	AddPeer(Peer) error
	RemoveRelationship(Relationship) error
	RemovePeer(Peer) error
}

func NewSociety() *society {
	s := society{
		peers:             map[string]Peer{},
		fromRelationships: map[Peer]Peer{},
		toRelationships:   map[Peer]Peer{},
	}
	return &s
}

var _ Society = (*society)(nil)

type society struct {
	peers             map[string]Peer
	fromRelationships map[Peer]Peer
	toRelationships   map[Peer]Peer
}

func (s *society) RemoveRelationship(rel Relationship) error {
	delete(s.fromRelationships, rel.From)
	delete(s.toRelationships, rel.To)
	return nil
}

func (s *society) RemovePeer(p Peer) error {
	hash := p.Hash()
	delete(s.peers, hash)
	delete(s.fromRelationships, p)
	delete(s.toRelationships, p)
	return nil
}

func (s *society) Relationship(p1, p2 Peer) (Relationship, error) {
	for from, to := range s.fromRelationships {
		if from == p1 && to == p2 {
			return NewRelationship(from, to), nil
		}
	}
	return NoRelationship, ErrNoRelationship
}

func (s *society) Peers() []Peer {
	ps := make([]Peer, 0, len(s.peers))
	for _, p := range s.peers {
		ps = append(ps, p)
	}
	return ps
}

func (s *society) Peer(hash string) (Peer, error) {
	p, exists := s.peers[hash]
	if !exists {
		return NoPeer, ErrNoPeer
	}
	return p, nil
}

func (s *society) Citizens() map[Hash]Citizen {
	c := map[Hash]Citizen{}
	for k, p := range s.peers {
		c[k] = Citizen{p, s}
	}
	return c
}

func (s *society) Citizen(hash Hash) (Citizen, error) {
	citizen, exists := s.peers[hash]
	if !exists {
		return NoCitizen, ErrNoCitizen
	}
	return Citizen{citizen, s}, nil
}

func (s *society) Relationships() []Relationship {
	rs := make([]Relationship, 0, len(s.fromRelationships))
	for k, v := range s.fromRelationships {
		r := Relationship{
			From:   k,
			To:     v,
			Nature: Follows,
		}
		rs = append(rs, r)
	}
	return rs
}

func (s *society) AddRelationship(from Peer, to Peer) error {
	if _, exists := s.fromRelationships[from]; !exists {
		return ErrRelationshipExists
	}
	s.fromRelationships[from] = to
	s.toRelationships[to] = from
	return nil
}

func (s *society) AddPeer(p Peer) error {
	s.peers[p.Hash()] = p
	return nil
}

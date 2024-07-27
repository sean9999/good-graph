package society

// a Citizen is a Peer in the context of a Society
type Citizen struct {
	Peer    Peer
	Society *society
}

var NoCitizen Citizen

func (c Citizen) Hash() Hash {
	return c.Peer.Hash()
}

func (c Citizen) Exists() bool {
	return c.Peer.Exists()
}
func (c Citizen) Equal(h Hasher) bool {
	return c.Hash() == h.Hash()
}
func (c Citizen) Followees() []Peer {
	cm := []Peer{}
	for k, v := range c.Society.fromRelationships {
		if k.Hash() == c.Hash() {
			cm = append(cm, v)
		}
	}
	return cm
}

func (c Citizen) Followers() []Peer {
	cm := []Peer{}
	for k, v := range c.Society.toRelationships {
		if k.Hash() == c.Hash() {
			cm = append(cm, v)
		}
	}
	return cm
}

func (c Citizen) Mutuals() []Peer {
	cs := make([]Peer, 0, len(c.Society.fromRelationships))
	m := map[string]Peer{}
	for _, p := range c.Followees() {
		m[p.Hash()] = p
	}
	for _, p := range c.Followers() {
		m[p.Hash()] = p
	}
	for _, v := range m {
		cs = append(cs, v)
	}
	return cs
}

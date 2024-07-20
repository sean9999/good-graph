package graph

// a Gaggle is a group of Peers
type Gaggle []Peer

func (gaggle Gaggle) AsMaps() []map[string]string {
	ms := []map[string]string{}
	for _, p := range gaggle {
		ms = append(ms, p.AsMap())
	}
	return ms
}

package graph

type Relationship struct {
	From Peer `json:"from"`
	To   Peer `json:"to"`
}

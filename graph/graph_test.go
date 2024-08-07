package graph

import (
	"crypto/rand"
	"io"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

type MockBroker struct {
	mock.Mock
	http.Handler
	inbox  chan Message
	outbox chan Message
	logger zerolog.Logger
}

func NewMockBroker() *MockBroker {
	m := MockBroker{
		inbox:  make(chan Message, 1),
		outbox: make(chan Message, 1024),
		logger: zerolog.New(io.Discard).Level(zerolog.FatalLevel),
	}
	return &m
}

// listen for events and drop them on the floor
func (b *MockBroker) Listen() chan Message {
	go func() {
		for range b.inbox {
			// drop
		}
	}()
	return b.inbox
}

func (b *MockBroker) Outbox() chan Message {
	return b.outbox
}

func (b *MockBroker) Logger() zerolog.Logger {
	return b.logger
}

type MockDatabase struct {
	mock.Mock
	Database
	peers map[string]Peer
	rels  map[string]Relationship
}

func (mdb *MockDatabase) Load() (Snapshot, error) {
	ps := make([]Peer, 0, len(mdb.peers))
	rs := make([]Relationship, 0, len(mdb.rels))
	for _, p := range mdb.peers {
		ps = append(ps, p)
	}
	for _, rel := range mdb.rels {
		rs = append(rs, rel)
	}
	return Snapshot{
		Peers:         ps,
		Relationships: rs,
	}, nil
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}

// func NewGraph(db Database, broker Broker, randy io.Reader) (Graph, error) {

func TestNewGraph(t *testing.T) {

	graph, err := NewGraph(NewMockDatabase(), NewMockBroker(), rand.Reader)

	if err != nil {
		t.Fatal(err)
	}

	p := graph.GeneratePeer()

	t.Error(p)

}

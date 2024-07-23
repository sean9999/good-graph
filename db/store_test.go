package db_test

import (
	"crypto/rand"
	"testing"

	"github.com/sean9999/good-graph/db"
	"github.com/sean9999/good-graph/graph"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_AddPeer(t *testing.T) {
	gdb := db.New("../testdata")
	graphPeer, err := graph.NewPeer(rand.Reader)
	assert.Nil(t, err, "creating a new peer should not produce an error")
	p1 := db.PeerFrom(graphPeer)
	err = gdb.AddPeer(p1)
	assert.Nil(t, err, "adding a peer should not produce an error")
	p2, err := gdb.GetPeer(p1.Nickname())
	assert.Nil(t, err, "getting a peer should be fine")
	assert.Equal(t, p1.Nickname(), p2.Nickname(), "peers are not the same")
}

func TestDatabase_AddRelationship(t *testing.T) {
	d := db.New("../testdata")
	// frog := new(db.Peer)
	// wood := new(db.Peer)
	frog, err := d.GetPeer("young-sea")
	assert.Nil(t, err, "getting a peer should be fine")
	wood, err := d.GetPeer("late-frog")

	rel := db.RelationshipFrom(*frog, *wood)
	err = d.AddRelationship(rel)
	assert.Nil(t, err)
}

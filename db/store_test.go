package db_test

import (
	"crypto/rand"
	"testing"

	"github.com/sean9999/good-graph/db"
	"github.com/sean9999/good-graph/graph"
	"github.com/sean9999/good-graph/transport"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_AddPeer(t *testing.T) {
	ch1 := make(chan transport.Msg, 10)
	ch2 := make(chan transport.Msg, 10)
	gdb := db.New("../testdata", ch1, ch2)
	graphPeer, err := graph.NewPeer(rand.Reader)
	assert.Nil(t, err, "creating a new peer should not produce an error")
	p1 := db.Peer(graphPeer)
	err = gdb.AddPeer(p1)
	assert.Nil(t, err, "adding a peer should not produce an error")
	p2, err := gdb.GetPeer(p1.Nickname())
	assert.Nil(t, err, "getting a peer should be fine")
	assert.Equal(t, p1.Nickname(), p2.Nickname(), "peers are not the same")
}

func TestDatabase_AddRelationship(t *testing.T) {
	ch1 := make(chan transport.Msg, 10)
	ch2 := make(chan transport.Msg, 10)
	d := db.New("../testdata", ch1, ch2)
	morn, err := d.GetPeer("silent-morning")
	assert.Nil(t, err, "getting a peer should be fine")
	frog, err := d.GetPeer("late-frog")
	rel := db.NewRelationship(frog, morn)
	err = d.AddRelationship(rel)
	assert.Nil(t, err)
}

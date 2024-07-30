package graph

// func TestDatabase_AddPeer(t *testing.T) {
// 	ch1 := make(chan Message, 10)
// 	ch2 := make(chan Message, 10)
// 	gdb := Message("../testdata", ch1, ch2)
// 	graphPeer := NewPeer(rand.Reader)
// 	assert.Nil(t, err, "creating a new peer should not produce an error")
// 	p1 := Peer(graphPeer)
// 	err = gdb.AddPeer(p1)
// 	assert.Nil(t, err, "adding a peer should not produce an error")
// 	p2, err := gdb.GetPeer(p1.Hash())
// 	assert.Nil(t, err, "getting a peer should be fine")
// 	assert.Equal(t, p1.Hash(), p2.Hash(), "peers are not the same")
// }

// func TestDatabase_AddRelationship(t *testing.T) {
// 	ch1 := make(chan Message, 10)
// 	ch2 := make(chan Message, 10)
// 	d := NewJsonStore("../testdata", ch1, ch2)
// 	morn := d.GetPeer("silent-morning")
// 	assert.Nil(t, err, "getting a peer should be fine")
// 	frog, err := d.GetPeer("late-frog")
// 	rel := NewRelationship(frog, morn)
// 	err = d.AddRelationship(rel)
// 	assert.Nil(t, err)
// }

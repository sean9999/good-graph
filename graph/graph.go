package graph

import (
	"errors"
)

var ErrAlreaderFriends = errors.New("these peers are already friends")
var ErrNoPeer = errors.New("peer doesn't exist")
var ErrUnmarshalRelationship = errors.New("couldn't unmarshal relationship")

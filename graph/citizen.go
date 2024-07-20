package graph

import "github.com/hmdsefi/gograph"

// a Citizen is a Peer in a Society
type Citizen = *gograph.Vertex[Peer]

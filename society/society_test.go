package society_test

// type allZeros struct{}

// func (az allZeros) Read(p []byte) (int, error) {
// 	zeroes := make([]byte, len(p))
// 	i := copy(p, zeroes)
// 	return i, io.EOF
// }

// type SocietySuite struct {
// 	suite.Suite
// 	Graph society.Graph
// 	io.Reader
// }

// func (suite *SocietySuite) SetupTest() {
// 	ch1 := make(chan society.Message, 100)
// 	ch2 := make(chan society.Message, 100)
// 	suite.Graph = society.NewGraph(ch1, ch2)
// }

// func (s *SocietySuite) TestNewSociety() {
// 	assert.Equal(s.T(), false, s.Graph.IsAcyclic(), "graph should not be acyclic")
// 	assert.Equal(s.T(), false, s.Graph.IsDirected(), "graph should be directed")
// }

// func (s *SocietySuite) AddPeer() {
// 	p, err := graph.NewPeer(allZeros{})
// 	assert.Nil(s.T(), err, "err should be nil")
// 	s.Graph.AddPeer(p)
// 	assert.Equal(s.T(), p.Nickname(), s.Graph.GetVertexByID(p).Label().Nickname(), "expected nicknames to match")
// }

package api

import (
	"encoding/json"
	"net/http"

	"github.com/sean9999/good-graph/society"
)

func GetVertices(soc society.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peers, err := soc.Peers()
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("not cool, dude"))
		}
		j, _ := json.Marshal(peers)
		w.Write(j)
	}
}

func GetVertexByNick(g society.Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nick := r.PathValue("nick")
		p, err := g.Peer(nick)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("not cool, dude"))
		}
		w.Write(p.Json())
	}
}

// func GetEdge(g society.Graph) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fromHex := r.PathValue("from")
// 		toHex := r.PathValue("to")

// 		fromPeer, err := graph.PeerFromHex([]byte(fromHex))
// 		if err != nil {
// 			msg := jMap{
// 				"err":  err.Error(),
// 				"msg":  "no such vertex",
// 				"side": "source",
// 			}
// 			w.WriteHeader(404)
// 			w.Write(msg.Json())
// 			return
// 		}
// 		toPeer, err := graph.PeerFromHex([]byte(toHex))
// 		if err != nil {
// 			msg := jMap{
// 				"err":  err.Error(),
// 				"msg":  "no such vertex",
// 				"side": "dest",
// 			}
// 			w.WriteHeader(404)
// 			w.Write(msg.Json())
// 			return
// 		}

// 		rel := graph.Relationship{
// 			From: fromPeer,
// 			To:   toPeer,
// 		}
// 		j, _ := json.Marshal(&rel)

// 		if soc.RelationshipExists(rel) {
// 			w.Write(j)
// 		} else {
// 			msg := jMap{
// 				"msg":    "no such relationship",
// 				"source": fromPeer.Nickname(),
// 				"dest":   toPeer.Nickname(),
// 			}
// 			w.WriteHeader(404)
// 			w.Write(msg.Json())
// 		}

// 	}
// }

// func GetNeighbours(g society.Graph) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		nick := r.PathValue("nick")
// 		vert := soc.GetVertexByNick(nick)
// 		vertices := soc.GetNeighbours(vert.Label())
// 		gaggle := graph.Gaggle{}
// 		for _, v := range vertices {
// 			p := v.Label()
// 			gaggle = append(gaggle, p)
// 		}
// 		j, _ := json.Marshal(gaggle.AsMaps())
// 		w.Write(j)
// 	}
// }

// func AddVertex(g graph.Graph) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		p, err := graph.NewPeer(rand.Reader)
// 		if err != nil {
// 			msg := jMap{
// 				"err": err,
// 			}
// 			w.WriteHeader(400)
// 			w.Write(msg.Json())
// 		}
// 		soc.AddPeer(p)
// 		database.AddPeer(db.Peer(p))
// 		msg := jMap{
// 			"msg": fmt.Sprintf("%s was added to the graph", p.Nickname()),
// 		}

// 		j, _ := p.MarshalJSON()

// 		//	advertise event
// 		ev := transport.Msg{
// 			MsgType: "AddVertex",
// 			Msg:     j,
// 			N:       1,
// 		}
// 		msgs <- ev

// 		w.Write(msg.Json())
// 	}
// }

// func Befriend(g society.Graph) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		twoPeers := [2]string{}
// 		buf := new(bytes.Buffer)
// 		buf.ReadFrom(r.Body)
// 		err := json.Unmarshal(buf.Bytes(), &twoPeers)
// 		msg := jMap{}
// 		p1, err := graph.PeerFromHex([]byte(twoPeers[0]))
// 		if err != nil {
// 			msg["err"] = err.Error()
// 			w.WriteHeader(404)
// 			w.Write(msg.Json())
// 			return
// 		}
// 		p2, err := graph.PeerFromHex([]byte(twoPeers[1]))
// 		if err != nil {
// 			msg["err"] = err.Error()
// 			w.WriteHeader(404)
// 			w.Write(msg.Json())
// 			return
// 		}
// 		var err1, err2 error
// 		err = soc.Befriend(p1, p2)
// 		if err == nil {
// 			rel1 := db.NewRelationship(p1, p2)
// 			rel2 := db.NewRelationship(p2, p1)
// 			j, _ := rel1.MarshalJSON()
// 			events <- transport.Msg{
// 				MsgType: "Befriend",
// 				Msg:     j,
// 				N:       2,
// 			}
// 			err1 = database.AddRelationship(rel1)
// 			err2 = database.AddRelationship(rel2)
// 		}
// 		if err1 != nil && err2 != nil {
// 			msg["err1"] = err1.Error()
// 			msg["err2"] = err2.Error()
// 			w.WriteHeader(507)
// 			w.Write(msg.Json())
// 			return
// 		}
// 		msg["from"] = p1.AsMap()
// 		msg["to"] = p2.AsMap()
// 		msg["msg"] = fmt.Sprintf("%s and %s are now friends", p1.Nickname(), p2.Nickname())
// 		w.Write(msg.Json())
// 	}
// }

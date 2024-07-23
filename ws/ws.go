package ws

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/sean9999/good-graph/transport"
)

// MotherShip contains pointers to all connections, and handles websockets
type MotherShip struct {
	Connections map[*websocket.Conn]bool
	Logger      zerolog.Logger
	Inbox       chan transport.Msg
	Outbox      chan transport.Msg
}

// a constructor, just in case
func NewMotherShip() *MotherShip {
	ms := MotherShip{
		Connections: map[*websocket.Conn]bool{},
		Logger:      zerolog.New(os.Stdout).Level(zerolog.DebugLevel),
		Inbox:       make(chan transport.Msg, 1024),
		Outbox:      make(chan transport.Msg, 2024),
	}
	return &ms
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

// our main http.Handler, mounted to "/ws" probably
func (m *MotherShip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.Logger.Print("upgrade:", err)
		return
	}

	defer func() {
		m.Logger.Info().Msg("websockets going down")
		c.Close()
		r.Body.Close()
	}()

	//	re-use this for each message
	var msg transport.Msg

WebSocketListener:
	for {

		_, msgBytes, err := c.ReadMessage()

		if err != nil {
			m.Logger.Println("error reading from websocket", err)
			break
		}

		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			m.Logger.Err(err)
		} else {
			m.Outbox <- msg
			switch msg.MsgType {
			case "marcoPolo":
				if msg.Msg == "marco" {
					msg.Msg = "polo"
				} else {
					msg.Msg = "marco"
				}
				msg.N++
				m.Logger.Info().Msgf("%v", msg)
				err = c.WriteJSON(msg)
			case "killYourself":
				m.Logger.Info().Msgf("%v", msg)
				err = c.WriteJSON(msg)
				break WebSocketListener
			default:
				m.Logger.Info().Msgf("%q and %q and %d", msg.MsgType, msg.Msg, msg.N)
			}

		}
	}
}

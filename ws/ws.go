package ws

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/sean9999/good-graph/transport"
)

var ErrMarcoPolo = errors.New("marco polo")
var ErrSuicide = errors.New("kill yourself")

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

var addr = flag.String("addr", "localhost:8282", "http service address")

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

	//	send anything in outbox to websocket channel
	go func() {
		for msg := range m.Outbox {
			err := c.WriteJSON(msg)
			if err != nil {
				m.Logger.Err(err)
			}
		}
	}()

	//	re-use this for each message
	var msg transport.Msg

WebSocketListener:
	for {

		_, msgBytes, err := c.ReadMessage()

		if err != nil {
			m.Logger.Err(err).Msg("reading from websocket")
			break
		}

		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			m.Logger.Err(err)
		} else {
			m.Inbox <- msg
			switch msg.MsgType {
			case "marcoPolo":
				if msg.String() == "marco" {
					msg.Set("polo")
				} else {
					msg.Set("marco")
				}
				msg.N++
				err = c.WriteJSON(msg)
				if err != nil {
					m.Logger.Err(fmt.Errorf("%w :%w", err, ErrMarcoPolo))
				} else {
					m.Logger.Info().Msgf("marco polo: %v", msg)
				}
			case "killYourself":
				err = c.WriteJSON(msg)
				if err != nil {
					m.Logger.Err(err)
				} else {
					m.Logger.Info().Msgf("kill yourself: %v", msg)
				}
				break WebSocketListener
			case "society/addNode":
				m.Logger.Info().Msgf("SOCIETY: %s: %v (%d)", msg.MsgType, msg.Msg, msg.N)
			default:
				m.Logger.Info().Msgf("%q and %q and %d", msg.MsgType, msg.Msg, msg.N)
			}

		}
	}
}

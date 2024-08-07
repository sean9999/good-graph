package graph

import (
	"encoding/json"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// a Broker handles [Message]s between two or more parties
type Broker interface {
	http.Handler
	Listen() chan Message
	Outbox() chan Message
	Logger() zerolog.Logger
}

var GlobalID atomic.Int64

// a Message is a uniquely identifiable peice of information with an intended purpose
type Message struct {
	MessageID    int64         `json:"mid"`
	ThreadID     int64         `json:"tid"`
	Subject      string        `json:"subject"`
	Peer         *Peer         `json:"peer"`
	Relationship *Relationship `json:"relationship"`
}

// *bus implements Broker
var _ Broker = (*bus)(nil)

type bus struct {
	logger zerolog.Logger
	inbox  chan Message
	outbox chan Message
}

func (b *bus) Listen() chan Message {
	return b.inbox
}
func (b *bus) Outbox() chan Message {
	return b.outbox
}
func (b *bus) Logger() zerolog.Logger {
	return b.logger
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (m *bus) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.logger.Print("upgrade:", err)
		return
	}

	defer func() {
		m.logger.Info().Msg("websockets going down")
		r.Body.Close()
		c.Close()
	}()

	//	send anything in outbox to websocket channel
	go func() {
		for msg := range m.outbox {
			err := c.WriteJSON(msg)
			if err != nil {
				m.logger.Err(err)
			}
		}
	}()

	//	send anything comming in on the websocket to inbox
	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			m.logger.Println(err)
			return
		} else {
			msg := new(Message)
			err := json.Unmarshal(p, msg)
			if err != nil {
				m.logger.Println("can't marshal to a message ", msg)
			}
			m.inbox <- *msg
		}
	}

}

func NewBus() *bus {
	b := bus{
		logger: zerolog.New(os.Stdout).Level(zerolog.DebugLevel),
		inbox:  make(chan Message, 1),
		outbox: make(chan Message, 1),
	}
	return &b
}

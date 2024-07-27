package society

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// a OutsideWorld handles messages and connections
type OutsideWorld interface {
	http.Handler
	Listen() chan Message
	Outbox() chan Message
	Logger() zerolog.Logger
	Connections() map[*websocket.Conn]bool
}

// *bus implements OutsideWorld
var _ OutsideWorld = (*bus)(nil)

type bus struct {
	connections map[*websocket.Conn]bool
	logger      zerolog.Logger
	inbox       chan Message
	outbox      chan Message
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
func (b *bus) Connections() map[*websocket.Conn]bool {
	return b.connections
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
		c.Close()
		r.Body.Close()
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
}

func NewBus() *bus {
	b := bus{
		connections: map[*websocket.Conn]bool{},
		logger:      zerolog.New(os.Stdout).Level(zerolog.DebugLevel),
		inbox:       make(chan Message, 1024),
		outbox:      make(chan Message, 2024),
	}
	return &b
}

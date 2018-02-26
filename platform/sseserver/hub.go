package sseserver

import (
	"strings"
	"time"
)

// A connection hub keeps track of all the active client connections, and
// handles broadcasting messages out to those connections that match the
// appropriate namespace.
type Hub struct {
	Broadcast   chan SSEMessage      // Inbound messages to propogate out.
	Connections map[*Connection]bool // Registered connections.
	Register    chan *Connection     // Register requests from the connections.
	Unregister  chan *Connection     // Unregister requests from connections.
	SentMsgs    uint64               // Msgs broadcast since startup
	StartupTime time.Time            // Time hub was created
}

func newHub() *Hub {
	return &Hub{
		Broadcast:   make(chan SSEMessage),
		Connections: make(map[*Connection]bool),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
		StartupTime: time.Now(),
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
		case c := <-h.Unregister:

			delete(h.Connections, c)
			close(c.Send)
		case msg := <-h.Broadcast:
			h.SentMsgs++
			formattedMsg := msg.sseFormat()
			for c := range h.Connections {
				if strings.HasPrefix(msg.Namespace, c.Namespace) {
					select {
					case c.Send <- formattedMsg:
					default:
						// log.Println("cant pass to a connection send chan, buffer is full -- kill it with fire")
						delete(h.Connections, c)
						close(c.Send)
						// go c.ws.Close()
						/* TODO: figure out what to do here...
						we are already closing the send channel, in *theory* shouldn't the
						connection clean up? I guess possible it doesnt if its deadlocked or
						something... is it?

						we want to make sure to always close the HTTP connection though,
						so server can never fill up max num of open sockets.
						*/
					}
				}
			}
		}
	}
}

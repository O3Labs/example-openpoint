package sseserver

import "sync"

// Interface to a SSE server.
//
// Exposes a send-only chan `broadcast`, any SSEMessage sent to this channel
// will be broadcast out to any connected clients subscribed to a namespace
// that matches the message.
type Server struct {
	Broadcast chan<- SSEMessage
	Hub       *Hub
}

var sseServer *Server

var once sync.Once

// Creates a new Server and returns a reference to it.
func CurrentServer() *Server {
	if sseServer != nil {
		return sseServer
	}
	once.Do(func() {
		// set up the public interface
		sseServer = &Server{}

		// start up our actual internal connection hub
		// which we keep in the server struct as private
		var h = newHub()
		sseServer.Hub = h
		go h.run()

		// expose just the broadcast chanel to public
		// will be typecast to send-only
		sseServer.Broadcast = h.Broadcast
	})
	return sseServer
}

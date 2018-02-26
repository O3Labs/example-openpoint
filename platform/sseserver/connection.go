package sseserver

import (
	"net/http"
	"time"
)

type Connection struct {
	R         *http.Request          // The HTTP request
	W         http.ResponseWriter    // The HTTP response
	Created   time.Time              // Timestamp for when connection was opened
	Send      chan []byte            // Buffered channel of outbound messages
	Namespace string                 // Conceptual "channel" SSE client is requesting
	Data      map[string]interface{} //useful data field if needed
	MsgsSent  uint64                 // Msgs the connection has sent (all time)
}

type ConnectionStatus struct {
	Path      string `json:"request_path"`
	Namespace string `json:"namespace"`
	Created   int64  `json:"created_at"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
	MsgsSent  uint64 `json:"msgs_sent"`
}

func (c *Connection) Status() ConnectionStatus {
	return ConnectionStatus{
		Path:      c.R.URL.Path,
		Namespace: c.Namespace,
		Created:   c.Created.Unix(),
		ClientIP:  c.R.RemoteAddr,
		UserAgent: c.R.UserAgent(),
		MsgsSent:  c.MsgsSent,
	}
}

func (c *Connection) Writer() {
	cn := c.W.(http.CloseNotifier)
	closer := cn.CloseNotify()

	for {
		select {
		case msg := <-c.Send:
			_, err := c.W.Write(msg)
			if err != nil {
				break
			}
			if f, ok := c.W.(http.Flusher); ok {
				f.Flush()
				c.MsgsSent++
			}
		case <-closer:

			return
		}
	}
}

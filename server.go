package main

import (
	"github.com/rs/zerolog/hlog"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type IncomingCallServer struct {
	Interval time.Duration
}

func NewIncomingCallServer(interval time.Duration) *IncomingCallServer {
	server := &IncomingCallServer{
		Interval: interval,
	}

	return server
}

func (s IncomingCallServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upgrade the websocket")
		return
	}

	// this a no-op on a properly closed connection
	defer c.Close(websocket.StatusInternalError, "deferred close")

	// we will never read, so we call CloseRead() to inform the package to handle
	// pong's (warning: writes to the websocket will now error!)
	ctx := c.CloseRead(r.Context())

	logger.Info().Msg("Received a new connection!")

	client := NewWebsocketConnection(c)
	client.SendEvents(ctx, logger, time.Second*5)
}

package main

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/rs/zerolog"
	"math/rand"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type RandomMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type IncomingCallEvent struct {
	FirstName   string    `json:"first_name" fake:"{firstname}"`
	LastName    string    `json:"last_name" fake:"{lastname}"`
	Timestamp   time.Time `json:"timestamp" fake:"skip"`
	SIP         string    `json:"sip" fake:"skip"`
	City        string    `json:"city" fake:"{city}"`
	State       string    `json:"state" fake:"{stateabr}"`
	PhoneNumber string    `json:"phone_number" fake:"{phone}"`
	Priority    uint8     `json:"priority" fake:"{number:1,100}"`
}

type WebsocketConnection struct {
	*websocket.Conn
	rand *rand.Rand
}

func NewWebsocketConnection(conn *websocket.Conn) *WebsocketConnection {
	c := &WebsocketConnection{}
	c.Conn = conn
	c.rand = rand.New(rand.NewSource(time.Now().Unix()))

	return c
}

func (c *WebsocketConnection) SendEvent(ctx context.Context) error {
	messageCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// do we want to send garbage? (yes, 20% of the time!)
	x := c.rand.Intn(100)
	if x <= 10 {
		return wsjson.Write(messageCtx, c.Conn, "this is intentionally bad data and should be discarded")
	}
	if x > 10 && x <= 20 {
		var message RandomMessage
		message.Timestamp = time.Now()
		message.Message = "Hello world!"

		return wsjson.Write(messageCtx, c.Conn, &message)
	}

	// legit data
	var event IncomingCallEvent
	gofakeit.Struct(&event)
	event.SIP = "https://127.0.0.1:33213/" + gofakeit.UUID()
	event.Timestamp = time.Now()
	event.Priority = uint8(c.rand.Intn(100))

	return wsjson.Write(messageCtx, c.Conn, &event)
}

func (c *WebsocketConnection) SendEvents(ctx context.Context, logger *zerolog.Logger, interval time.Duration) {
	// loop, sending out an item every few seconds
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.
				Info().
				Msg("Closing connection")
			c.Close(websocket.StatusNormalClosure, "context is done")
			return
		case <-ticker.C:
			err := c.SendEvent(ctx)
			if err != nil {
				logger.
					Error().
					Err(err).
					Msg("failed to write event to the websocket")
				return
			}
			logger.Info().Msg("Sent an event")
		}
	}
}

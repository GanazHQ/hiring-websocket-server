package main

import (
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
	"os"
	"time"
)

func main() {
	// set up logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// get port
	port := os.Getenv("PORT")
	if port == "" {
		// default
		port = "7777"
	}

	websocketServer := NewIncomingCallServer(time.Second * 3)

	middleware := alice.New(hlog.NewHandler(logger))

	srv := &http.Server{
		Addr: "0.0.0.0:" + port,

		// default timeouts are no good for websockets, since the HTTP request effectively
		// remains open indefinitely
		ReadHeaderTimeout: time.Second * 30,
		WriteTimeout:      0,
		IdleTimeout:       0,
		Handler:           middleware.Then(websocketServer),
	}

	logger.Info().Str("PORT", port).Msg("Serving...")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

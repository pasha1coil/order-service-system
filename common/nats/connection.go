package nats

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Configuration struct {
	URL  string `env:"NATS_URL"`
	Name string `env:"NATS_CLIENT_NAME"`
}

func Connect(cfg Configuration) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS connection closed: %v", nc.LastError())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("NATS error: %v", err)
		}),
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	if !conn.IsConnected() {
		conn.Close()
		return nil, fmt.Errorf("NATS connection not established")
	}

	return conn, nil
}

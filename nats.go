package nats

import (
	"encoding/json"
	"time"

	"github.com/zoobr/csxlib/logger"

	"github.com/nats-io/nats.go"
)

type Config struct {
	ClientName       string        // optional name label which will be sent to the server on CONNECT to identify the client.
	URL              string        // NATS URL
	ReconnectionTime time.Duration // wait time between reconnect attempts (optional)
	User             string
	Password         string
}

type Conn struct {
	*nats.Conn
}

func Connect(cfg *Config) (*Conn, error) {
	opts := []nats.Option{
		nats.Name(cfg.ClientName),
		nats.RetryOnFailedConnect(true),
		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			logger.Error(err, "url", nc.ConnectedUrl(), "queue", s.Queue, "subject", s.Subject, "stats", nc.Stats())
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error(err, "msg", "disconnected", "url", nc.ConnectedUrl(), "stats", nc.Stats())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Info("connection closed", "url", nc.ConnectedUrl(), "stats", nc.Stats())
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("successfully reconnected", "url", nc.ConnectedUrl(), "stats", nc.Stats())
		}),
	}
	if cfg.ReconnectionTime != 0 {
		opts = append(opts, nats.ReconnectWait(cfg.ReconnectionTime))
	}
	if cfg.User != "" && cfg.Password != "" {
		opts = append(opts, nats.UserInfo(cfg.User, cfg.Password))
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: nc}, nil
}

// Publish publishes the data interface to the given subject. The data
// argument is left untouched and needs to be correctly interpreted on
// the receiver.
func (c *Conn) Publish(subject string, data interface{}) error {
	if d, ok := data.([]byte); ok {
		return c.Conn.Publish(subject, d)
	}

	encodedData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Conn.Publish(subject, encodedData)
}

package queue

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/nats-io/nuid"
	"math"
)

const (
	defaultHost = "localhost"
	defaultPort = 50000
)

type options struct {
	host              string
	port              int
	clientId          string
	authToken         string
	channels          []string
	expirationSeconds int
	delaySeconds      int
	maxReceiveCount   int
	deadLetterQueue   string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.host, o.port, err = cfg.MustParseAddress("address", defaultHost)
	if err != nil {
		return options{}, fmt.Errorf("error parsing address value, %w", err)
	}
	o.authToken = cfg.ParseString("auth_token", "")
	o.clientId = cfg.ParseString("client_id", nuid.Next())
	o.channels = cfg.ParseStringList("channels")

	o.expirationSeconds, err = cfg.ParseIntWithRange("expiration_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing expiration seconds, %w", err)
	}
	o.delaySeconds, err = cfg.ParseIntWithRange("delay_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing delay seconds, %w", err)
	}
	o.maxReceiveCount, err = cfg.ParseIntWithRange("max_receive_count", 1, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error max receive count seconds")
	}
	o.deadLetterQueue = cfg.ParseString("dead_letter_queue", "")
	return o, nil
}

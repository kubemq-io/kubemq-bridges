package queue

import (
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
	"math"
)

const (
	defaultHost = "localhost:50000"
)

type options struct {
	host              string
	port              int
	clientId          string
	authToken         string
	channel           string
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
	o.clientId = cfg.ParseString("client_id", uuid.New().String())
	o.channel = cfg.ParseString("channel", "")
	if o.channel != "" {
		o.channels = append(o.channels, o.channel)
	} else {
		o.channels = cfg.ParseStringList("channels")
		if len(o.channels) == 0 {
			return options{}, fmt.Errorf("error parsing channles, cannot be empty")
		}
	}
	o.expirationSeconds, err = cfg.ParseIntWithRange("expiration_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing expiration seconds, %w", err)
	}
	o.delaySeconds, err = cfg.ParseIntWithRange("delay_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing delay seconds, %w", err)
	}
	o.maxReceiveCount, err = cfg.ParseIntWithRange("max_receive_count", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error max receive count seconds")
	}
	o.deadLetterQueue = cfg.ParseString("dead_letter_queue", "")
	return o, nil
}

package query

import (
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
	"math"
)

const (
	defaultHost           = "localhost:5000"
	defaultTimeoutSeconds = 600
)

type options struct {
	host           string
	port           int
	clientId       string
	channel        string
	authToken      string
	defaultChannel string
	timeoutSeconds int
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
		o.defaultChannel = o.channel
	} else {
		o.defaultChannel = cfg.ParseString("default_channel", "")
		if o.defaultChannel == "" {
			return options{}, fmt.Errorf("error parsing channel, cannot be empty")
		}
	}

	o.timeoutSeconds, err = cfg.ParseIntWithRange("timeout_seconds", defaultTimeoutSeconds, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing timeout seconds value, %w", err)
	}
	return o, nil
}

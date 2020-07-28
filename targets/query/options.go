package query

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/nats-io/nuid"
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
	o.clientId = cfg.ParseString("client_id", nuid.Next())
	o.defaultChannel = cfg.ParseString("default_channel", "")
	o.timeoutSeconds, err = cfg.ParseIntWithRange("timeout_seconds", defaultTimeoutSeconds, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing timeout seconds value, %w", err)
	}
	return o, nil
}

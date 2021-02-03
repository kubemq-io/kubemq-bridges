package queue_stream

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/uuid"
)

const (
	defaultWaitTimeout       = 3600
	defaultVisibilityTimeout = 3600
	defaultSources           = 1
)

type options struct {
	host              string
	port              int
	clientId          string
	authToken         string
	channel           string
	sources           int
	visibilityTimeout int
	waitTimeout       int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.host, o.port, err = cfg.MustParseAddress("address", "")
	if err != nil {
		return options{}, fmt.Errorf("error parsing address value, %w", err)
	}
	o.authToken = cfg.ParseString("auth_token", "")

	o.clientId = cfg.ParseString("client_id", uuid.New().String())

	o.channel, err = cfg.MustParseString("channel")
	if err != nil {
		return options{}, fmt.Errorf("error parsing channel value, %w", err)
	}
	o.sources, err = cfg.ParseIntWithRange("sources", defaultSources, 1, 100)
	if err != nil {
		return options{}, fmt.Errorf("error parsing sources value, %w", err)
	}

	o.visibilityTimeout, err = cfg.ParseIntWithRange("visibility_timeout_seconds", defaultVisibilityTimeout, 1, 24*60*60)
	if err != nil {
		return options{}, fmt.Errorf("error parsing visibility timeout value, %w", err)
	}
	o.waitTimeout, err = cfg.ParseIntWithRange("wait_timeout", defaultWaitTimeout, 1, 24*60*60)
	if err != nil {
		return options{}, fmt.Errorf("error parsing wait timeout value, %w", err)
	}
	return o, nil
}

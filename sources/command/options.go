package command

import (
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
	"time"
)

const (
	defaultAddress       = "0.0.0.0:50000"
	defaultAutoReconnect = true
	defaultSources       = 1
)

type options struct {
	host                     string
	port                     int
	clientId                 string
	authToken                string
	channel                  string
	group                    string
	autoReconnect            bool
	reconnectIntervalSeconds time.Duration
	maxReconnects            int
	sources                  int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.host, o.port, err = cfg.MustParseAddress("address", defaultAddress)
	if err != nil {
		return options{}, fmt.Errorf("error parsing address value, %w", err)
	}
	o.authToken = cfg.ParseString("auth_token", "")

	o.clientId = cfg.ParseString("client_id", uuid.New().String())

	o.channel, err = cfg.MustParseString("channel")
	if err != nil {
		return o, fmt.Errorf("error parsing channel value, %w", err)
	}
	o.sources, err = cfg.ParseIntWithRange("sources", defaultSources, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing sources value, %w", err)
	}
	o.group = cfg.ParseString("group", "")
	o.autoReconnect = cfg.ParseBool("auto_reconnect", defaultAutoReconnect)
	interval, err := cfg.ParseIntWithRange("reconnect_interval_seconds", 1, 1, 1000000)
	if err != nil {
		return o, fmt.Errorf("error parsing reconnect interval seconds value, %w", err)
	}

	o.reconnectIntervalSeconds = time.Duration(interval) * time.Second

	o.maxReconnects = cfg.ParseInt("max_reconnects", 0)

	return o, nil
}

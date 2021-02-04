package queue

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/uuid"
)

const (
	defaultAddress     = "0.0.0.0:50000"
	defaultBatchSize   = 1
	defaultWaitTimeout = 60
	defaultSources     = 1
)

type options struct {
	host        string
	port        int
	clientId    string
	authToken   string
	channel     string
	sources     int
	batchSize   int
	waitTimeout int
	maxRequeue  int
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
		return options{}, fmt.Errorf("error parsing channel value, %w", err)
	}
	o.sources, err = cfg.ParseIntWithRange("sources", defaultSources, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing sources value, %w", err)
	}

	o.batchSize, err = cfg.ParseIntWithRange("batch_size", defaultBatchSize, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing batch size value, %w", err)
	}
	o.waitTimeout, err = cfg.ParseIntWithRange("wait_timeout", defaultWaitTimeout, 1, 24*60*60)
	if err != nil {
		return options{}, fmt.Errorf("error parsing wait timeout value, %w", err)
	}
	o.maxRequeue, err = cfg.ParseIntWithRange("max_requeue", 0, 0, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max requeue value, %w", err)
	}
	return o, nil
}

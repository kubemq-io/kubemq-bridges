package events_store

import (
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
)

const (
	defaultHost = "localhost:50000"
)

type options struct {
	host      string
	port      int
	clientId  string
	authToken string
	channels  []string
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
	o.channels = cfg.ParseStringList("channels")
	return o, nil
}

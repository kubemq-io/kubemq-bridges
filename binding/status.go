package binding

import (
	"github.com/kubemq-io/kubemq-bridges/config"
)

type Status struct {
	Binding      string            `json:"binding"`
	Ready        bool              `json:"ready"`
	SourceType   string            `json:"source_type"`
	SourceConfig []config.Metadata `json:"source_config"`
	TargetType   string            `json:"target_type"`
	TargetConfig []config.Metadata `json:"target_config"`
}

func newStatus(cfg config.BindingConfig) *Status {
	return &Status{
		Binding:      cfg.Name,
		Ready:        false,
		SourceType:   cfg.Sources.Kind,
		SourceConfig: cfg.Sources.Connections,
		TargetType:   cfg.Targets.Kind,
		TargetConfig: cfg.Targets.Connections,
	}
}

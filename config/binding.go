package config

import (
	"fmt"
)

type BindingConfig struct {
	Name       string   `json:"name"`
	Sources    Spec     `json:"sources"`
	Targets    Spec     `json:"targets"`
	Properties Metadata `json:"properties"`
}

func (b BindingConfig) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("binding must have name")
	}
	if err := b.Sources.Validate(); err != nil {
		return fmt.Errorf("binding sources error, %w", err)
	}
	if err := b.Targets.Validate(); err != nil {
		return fmt.Errorf("binding targets error, %w", err)
	}
	return nil
}

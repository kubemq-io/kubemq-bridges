package config

import (
	"fmt"
)

type Metadata map[string]string

type Spec struct {
	Name        string     `json:"name"`
	Kind        string     `json:"kind"`
	Connections []Metadata `json:"connections"`
}

func (s Spec) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if s.Kind == "" {
		return fmt.Errorf("kind cannot be empty")
	}
	if len(s.Connections) == 0 {
		return fmt.Errorf("no connections found")
	}
	return nil
}

package common

import (
	"fmt"
)

type Metadata struct {
	Name          string   `json:"name"`
	Kind          string   `json:"kind"`
	Description   string   `json:"description"`
	Default       string   `json:"default"`
	Options       []string `json:"options"`
	Must          bool     `json:"must"`
	Min           int      `json:"min"`
	Max           int      `json:"max"`
	LoadedOptions string
}

func NewMetadata() *Metadata {
	return &Metadata{}
}

func (m *Metadata) SetName(value string) *Metadata {
	m.Name = value
	return m
}

func (m *Metadata) SetKind(value string) *Metadata {
	m.Kind = value
	return m
}
func (m *Metadata) SetDescription(value string) *Metadata {
	m.Description = value
	return m
}

func (m *Metadata) SetDefault(value string) *Metadata {
	m.Default = value
	return m
}
func (m *Metadata) SetLoadedOptions(value string) *Metadata {
	m.LoadedOptions = value
	return m
}
func (m *Metadata) SetOptions(value []string) *Metadata {
	m.Options = value
	return m
}
func (m *Metadata) SetMust(value bool) *Metadata {
	m.Must = value
	return m
}

func (m *Metadata) SetMin(value int) *Metadata {
	m.Min = value
	return m
}
func (m *Metadata) SetMax(value int) *Metadata {
	m.Max = value
	return m
}
func (m *Metadata) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("metadata name cannot be empty")
	}
	if m.Kind == "" {
		return fmt.Errorf("metadata kind cannot be empty")
	}

	if m.Description == "" {
		return fmt.Errorf("metadata description cannot be empty")
	}

	return nil
}

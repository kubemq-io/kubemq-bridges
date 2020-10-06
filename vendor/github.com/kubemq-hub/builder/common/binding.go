package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Binding struct {
	Name              string            `json:"name"`
	Source            Spec              `json:"source"`
	Target            Spec              `json:"target"`
	Properties        map[string]string `json:"properties"`
	loadedOptions     DefaultOptions
	targetsList       Connectors
	sourcesList       Connectors
	takenBindingNames []string
}

func NewBinding() *Binding {
	return &Binding{
		Name:              "",
		Source:            Spec{},
		Target:            Spec{},
		Properties:        map[string]string{},
		loadedOptions:     nil,
		targetsList:       nil,
		sourcesList:       nil,
		takenBindingNames: nil,
	}
}
func (b *Binding) SetDefaultOptions(value DefaultOptions) *Binding {
	b.loadedOptions = value
	return b
}
func (b *Binding) SetTargetsList(value Connectors) *Binding {
	b.targetsList = value
	return b
}
func (b *Binding) SetSourcesList(value Connectors) *Binding {
	b.sourcesList = value
	return b
}
func (b *Binding) SetTakenBindingNames(value []string) *Binding {
	b.takenBindingNames = value
	return b
}
func (b *Binding) SourceName() string {
	return b.Source.Name
}
func (b *Binding) TargetName() string {
	return b.Target.Name
}
func (b *Binding) askKind(kinds []string) (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Select Connector Kind").
		SetDefault(kinds[0]).
		SetOptions(kinds).
		SetHelp("Select Connector Kind").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}
func (b *Binding) askSource() error {
	var err error
	if b.Source.Name, err = NewName().
		RenderSource(); err != nil {
		return err
	}
	var kinds []string
	sources := make(map[string]*Connector)
	for _, c := range b.sourcesList {
		kinds = append(kinds, c.Kind)
		sources[c.Kind] = c
	}
	if len(kinds) == 0 {
		return fmt.Errorf("no source connectors available")
	}
	if b.Source.Kind, err = b.askKind(kinds); err != nil {
		return err
	}
	connector := sources[b.Source.Kind]
	if b.Source.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}
func (b *Binding) askTarget() error {
	var err error
	if b.Target.Name, err = NewName().
		RenderTarget(); err != nil {
		return err
	}
	var kinds []string
	targets := make(map[string]*Connector)
	for _, c := range b.targetsList {
		kinds = append(kinds, c.Kind)
		targets[c.Kind] = c
	}
	if len(kinds) == 0 {
		return fmt.Errorf("no targets connectors available")
	}

	if b.Target.Kind, err = b.askKind(kinds); err != nil {
		return err
	}
	connector := targets[b.Target.Kind]
	if b.Target.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}

func (b *Binding) Render() (*Binding, error) {
	var err error
	if b.Name, err = NewName().
		SetTakenNames(b.takenBindingNames).
		RenderBinding(); err != nil {
		return nil, err
	}
	if err := b.askSource(); err != nil {
		return nil, err
	}
	if err := b.askTarget(); err != nil {
		return nil, err
	}
	if b.Properties, err = NewProperties().
		Render(); err != nil {
		return nil, err
	}
	return b, nil
}

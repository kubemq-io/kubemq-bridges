package binding

import (
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/builder/connector/bridges/source"
	"github.com/kubemq-hub/builder/connector/bridges/target"
)

type Binding struct {
	Name              string            `json:"name"`
	Sources           *source.Source    `json:"sources"`
	Targets           *target.Target    `json:"targets"`
	Properties        map[string]string `json:"properties"`
	addressOptions    []string
	takenSourceNames  []string
	takenTargetsNames []string
	takenBindingNames []string
}

func NewBinding() *Binding {
	return &Binding{}
}
func (b *Binding) SetAddress(value []string) *Binding {
	b.addressOptions = value
	return b
}
func (b *Binding) SetTakenSourceNames(value []string) *Binding {
	b.takenSourceNames = value
	return b
}
func (b *Binding) SetTakenTargetsNames(value []string) *Binding {
	b.takenTargetsNames = value
	return b
}
func (b *Binding) SetTakenBindingNames(value []string) *Binding {
	b.takenBindingNames = value
	return b
}
func (b *Binding) SourceName() string {
	if b.Sources != nil {
		return b.Sources.Name
	}
	return ""
}
func (b *Binding) TargetName() string {
	if b.Targets != nil {
		return b.Targets.Name
	}
	return ""
}
func (b *Binding) BindingName() string {
	return b.Name
}
func (b *Binding) Render() (*Binding, error) {
	var err error
	if b.Name, err = NewName().
		SetTakenNames(b.takenBindingNames).
		Render(); err != nil {
		return nil, err
	}
	if b.Sources, err = source.NewSource().
		SetAddress(b.addressOptions).
		SetTakenNames(b.takenSourceNames).
		Render(); err != nil {
		return nil, err
	}
	if b.Targets, err = target.NewTarget().
		SetAddress(b.addressOptions).
		SetTakenNames(b.takenTargetsNames).
		Render(); err != nil {
		return nil, err
	}
	if b.Properties, err = common.NewProperties().
		Render(); err != nil {
		return nil, err
	}
	return b, nil
}

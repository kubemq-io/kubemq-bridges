package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"gopkg.in/yaml.v2"
)

type Builder struct {
	Bindings          []*Binding `json:"bindings"`
	manifest          *Manifest
	loadedOptions     DefaultOptions
	takenBindingNames []string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetManifest(value *Manifest) *Builder {
	b.manifest = value
	return b
}
func (b *Builder) SetOptions(value DefaultOptions) *Builder {
	b.loadedOptions = value
	return b
}
func (b *Builder) askAddBinding() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-binding").
		SetMessage("Would you like to add another binding").
		SetDefault("false").
		SetHelp("Add new binding").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false, err
	}
	return val, nil
}

func (b *Builder) addBinding() error {
	if bnd, err := NewBinding().
		SetDefaultOptions(b.loadedOptions).
		SetSourcesList(b.manifest.Sources).
		SetTargetsList(b.manifest.Targets).
		SetTakenBindingNames(b.takenBindingNames).
		Render(); err != nil {
		return err
	} else {
		b.Bindings = append(b.Bindings, bnd)
		b.takenBindingNames = append(b.takenBindingNames, bnd.Name)
	}
	return nil
}

func (b *Builder) Render() ([]byte, error) {
	if b.manifest == nil {
		return nil, fmt.Errorf("inavlid manifest")
	}
	err := b.addBinding()
	if err != nil {
		return nil, err
	}
	for {
		addMore, err := b.askAddBinding()
		if err != nil {
			return nil, err
		}
		if addMore {
			err := b.addBinding()
			if err != nil {
				return nil, err
			}
		} else {
			goto done
		}
	}
done:
	return yaml.Marshal(b)
}

func (b *Builder) Yaml() ([]byte, error) {
	return yaml.Marshal(b)
}

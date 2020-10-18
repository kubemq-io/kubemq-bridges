package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"gopkg.in/yaml.v2"
)

type Bindings struct {
	Bindings          []*Binding `json:"bindings"`
	manifest          *Manifest
	loadedOptions     DefaultOptions
	takenBindingNames []string
	defaultName       string
}

func NewBindings(defaultName string) *Bindings {
	return &Bindings{
		defaultName: defaultName,
	}
}

func (b *Bindings) SetManifest(value *Manifest) *Bindings {
	b.manifest = value
	return b
}
func (b *Bindings) SetOptions(value DefaultOptions) *Bindings {
	b.loadedOptions = value
	return b
}
func (b *Bindings) askAddBinding() (bool, error) {
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
func (b *Bindings) confirmBinding(bnd *Binding) bool {
	utils.Println(fmt.Sprintf(promptBindingConfirm, bnd.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println(promptBindingReconfigure)
	}
	return val
}
func (b *Bindings) addBinding() error {
	for {
		bnd := NewBinding(fmt.Sprintf("%s-binding-%d", b.defaultName, len(b.Bindings)+1))
		var err error
		if bnd, err = bnd.
			SetDefaultOptions(b.loadedOptions).
			SetSourcesList(b.manifest.Sources).
			SetTargetsList(b.manifest.Targets).
			SetTakenBindingNames(b.takenBindingNames).
			Render(); err != nil {
			return err
		}
		ok := b.confirmBinding(bnd)
		if ok {
			b.Bindings = append(b.Bindings, bnd)
			b.takenBindingNames = append(b.takenBindingNames, bnd.Name)
			break
		}
	}
	return nil
}

func (b *Bindings) Render() ([]byte, error) {
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

func (b *Bindings) Yaml() ([]byte, error) {
	return yaml.Marshal(b)
}

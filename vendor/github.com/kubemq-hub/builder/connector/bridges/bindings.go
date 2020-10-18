package bridges

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Bindings struct {
	Bindings          []*Binding `json:"bindings"`
	defaultOptions    common.DefaultOptions
	takenBindingNames []string
	takenSourceNames  []string
	takenTargetNames  []string
	addressOptions    []string
	defaultName       string
}

func NewBindings(defaultName string) *Bindings {
	return &Bindings{
		defaultName: defaultName,
	}
}
func (b *Bindings) SetDefaultOptions(value common.DefaultOptions) *Bindings {
	b.defaultOptions = value
	return b
}

func (b *Bindings) askAddBinding() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-binding").
		SetMessage("Would you like to add another bindings").
		SetDefault("false").
		SetHelp("Add new bindings bridge").
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
			SetAddress(b.addressOptions).
			SetTakenBindingNames(b.takenBindingNames).
			SetTakenSourceNames(b.takenSourceNames).
			SetTakenTargetsNames(b.takenTargetNames).
			Render(); err != nil {
			return err
		}
		ok := b.confirmBinding(bnd)
		if ok {
			b.Bindings = append(b.Bindings, bnd)
			b.takenBindingNames = append(b.takenBindingNames, bnd.BindingName())
			b.takenSourceNames = append(b.takenSourceNames, bnd.SourceName())
			b.takenTargetNames = append(b.takenTargetNames, bnd.TargetName())
			break
		}
	}

	return nil
}
func (b *Bindings) Render() ([]byte, error) {
	err := b.addBinding()
	if err != nil {
		return nil, err
	}
	for {
		addMore, err := b.askAddBinding()
		if err != nil {
			return yaml.Marshal(b)
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

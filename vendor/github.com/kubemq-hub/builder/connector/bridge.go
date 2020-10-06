package connector

import (
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/bridges/binding"
	"github.com/kubemq-hub/builder/survey"
)

type Bridge struct {
	Bindings          []*binding.Binding `json:"bindings"`
	takenBindingNames []string
	takenSourceNames  []string
	takenTargetNames  []string
	addressOptions    []string
}

func NewBridge() *Bridge {
	return &Bridge{}
}
func (b *Bridge) SetClusterAddress(value []string) *Bridge {
	b.addressOptions = value
	return b
}

func (b *Bridge) askAddBinding() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-binding").
		SetMessage("Would you like to add another bindings bridge").
		SetDefault("false").
		SetHelp("Add new bindings bridge").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false, err
	}
	return val, nil
}
func (b *Bridge) addBinding() error {
	if bnd, err := binding.NewBinding().
		SetAddress(b.addressOptions).
		SetTakenBindingNames(b.takenBindingNames).
		SetTakenSourceNames(b.takenSourceNames).
		SetTakenTargetsNames(b.takenTargetNames).
		Render(); err != nil {
		return err
	} else {
		b.Bindings = append(b.Bindings, bnd)
		b.takenBindingNames = append(b.takenBindingNames, bnd.BindingName())
		b.takenSourceNames = append(b.takenSourceNames, bnd.SourceName())
		b.takenTargetNames = append(b.takenTargetNames, bnd.TargetName())
	}
	return nil
}
func (b *Bridge) Render() ([]byte, error) {
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

func (b *Bridge) Yaml() ([]byte, error) {
	return yaml.Marshal(b)
}

package bridges

import (
	"github.com/kubemq-hub/builder/connector/common"
)

type Bridges struct {
	defaultOptions common.DefaultOptions
	defaultName    string
}

func NewBridges(defaultName string) *Bridges {
	return &Bridges{
		defaultName: defaultName,
	}
}
func (b *Bridges) SetDefaultOptions(value common.DefaultOptions) *Bridges {
	b.defaultOptions = value
	return b
}

func (b *Bridges) Render() ([]byte, error) {
	if bnd, err := NewBindings(b.defaultName).
		SetDefaultOptions(b.defaultOptions).
		Render(); err != nil {
		return nil, err
	} else {
		return bnd, nil
	}

}

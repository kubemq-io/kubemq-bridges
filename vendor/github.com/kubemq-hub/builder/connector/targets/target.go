package targets

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"io/ioutil"
)

type Target struct {
	manifestData   []byte
	defaultOptions common.DefaultOptions
	defaultName    string
}

func NewTarget(defaultName string) *Target {
	return &Target{
		defaultName: defaultName,
	}
}

func (t *Target) SetManifest(value []byte) *Target {
	t.manifestData = value
	return t
}
func (t *Target) SetManifestFile(filename string) *Target {
	t.manifestData, _ = ioutil.ReadFile(filename)
	return t
}
func (t *Target) SetDefaultOptions(value common.DefaultOptions) *Target {
	t.defaultOptions = value
	return t
}

func (t *Target) Render() ([]byte, error) {
	if t.manifestData == nil {
		return nil, fmt.Errorf("invalid manifest")
	}
	m, err := common.LoadManifest(t.manifestData)
	if err != nil {
		return nil, err
	}
	if m.Schema != "targets" {
		return nil, fmt.Errorf("invalid scheme, %s", m.Schema)
	}
	if b, err := common.NewBindings(t.defaultName).
		SetManifest(m).
		SetOptions(t.defaultOptions).
		Render(); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

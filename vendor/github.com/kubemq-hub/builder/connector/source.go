package connector

import (
	"fmt"
	"github.com/kubemq-hub/builder/common"
	"io/ioutil"
)

type Source struct {
	manifestData   []byte
	defaultOptions common.DefaultOptions
}

func NewSource() *Source {
	return &Source{}
}

func (s *Source) SetManifest(value []byte) *Source {
	s.manifestData = value
	return s
}
func (s *Source) SetManifestFile(filename string) *Source {
	s.manifestData, _ = ioutil.ReadFile(filename)
	return s
}

func (s *Source) SetDefaultOptions(value common.DefaultOptions) *Source {
	s.defaultOptions = value
	return s
}

func (s *Source) Render() ([]byte, error) {
	if s.manifestData == nil {
		return nil, fmt.Errorf("invalid manifest")
	}
	m, err := common.LoadManifest(s.manifestData)
	if err != nil {
		return nil, err
	}
	if m.Schema != "sources" {
		return nil, fmt.Errorf("invalid scheme")
	}
	if b, err := common.NewBuilder().
		SetManifest(m).
		SetOptions(s.defaultOptions).
		Render(); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"strings"
)

type Spec struct {
	Kind           string            `json:"kind"`
	Properties     map[string]string `json:"properties"`
	PropertiesSpec string            `json:"-" yaml:"-"`
}

func NewSpec() *Spec {
	return &Spec{}
}
func (s *Spec) ColoredYaml(template string) string {
	s.PropertiesSpec = utils.MapToYaml(s.Properties)
	tpl := utils.NewTemplate(template, &s)
	spec, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring spec,%s", err.Error())
	}
	return string(spec)
}
func (s *Spec) Clone() *Spec {
	newSpec := &Spec{
		Kind:           s.Kind,
		Properties:     map[string]string{},
		PropertiesSpec: "",
	}
	for key, val := range s.Properties {
		newSpec.Properties[key] = val
	}
	return newSpec
}
func (s *Spec) TableItemShort() string {
	return s.Kind
}

func (s *Spec) IsKubemqAddress(address string) bool {
	if val, ok := s.Properties["address"]; ok {
		return strings.HasPrefix(val, address)
	}
	return false
}

func (s *Spec) Title() string {
	channel, ok := s.Properties["channel"]
	if ok {
		return fmt.Sprintf("%s/%s", s.Kind, channel)
	} else {
		return s.Kind
	}
}

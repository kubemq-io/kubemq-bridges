package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
)

type Spec struct {
	Name           string            `json:"name"`
	Kind           string            `json:"kind"`
	Properties     map[string]string `json:"properties"`
	PropertiesSpec string            `json:"-" yaml:"-"`
}

func (s Spec) String(template string) string {
	s.PropertiesSpec = utils.MapToYaml(s.Properties)
	tpl := utils.NewTemplate(template, &s)
	spec, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring spec,%s", err.Error())
	}
	return string(spec)
}

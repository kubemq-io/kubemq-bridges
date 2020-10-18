package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

const propertiesTml = `
<red>properties:</>
{{ .ValuesSpec | indent 2}}
`

type Properties struct {
	values     map[string]string
	ValuesSpec string
}

func NewProperties() *Properties {
	return &Properties{
		values: map[string]string{},
	}
}

func (p *Properties) Render() (map[string]string, error) {
	boolVal := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-middleware").
		SetMessage("Would you to add middlewares to this binding").
		SetDefault("false").
		SetHelp("Add a middleware properties").
		SetRequired(true).
		Render(&boolVal)
	if err != nil {
		return nil, err
	}
	if !boolVal {
		return nil, nil
	}
	if values, err := NewLog().Render(); err != nil {
		return nil, err
	} else {
		for key, val := range values {
			p.values[key] = val
		}
	}
	if values, err := NewRateLimiter().Render(); err != nil {
		return nil, err
	} else {
		for key, val := range values {
			p.values[key] = val
		}
	}
	if values, err := NewRetry().Render(); err != nil {
		return nil, err
	} else {
		for key, val := range values {
			p.values[key] = val
		}
	}
	return p.values, nil
}
func (p *Properties) String() string {
	if len(p.values) == 0 {
		return "\n<red>properties: {}</>"
	}
	p.ValuesSpec = utils.MapToYaml(p.values)
	tpl := utils.NewTemplate(propertiesTml, p)
	b, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring properties spec,%s", err.Error())
	}
	return string(b)
}

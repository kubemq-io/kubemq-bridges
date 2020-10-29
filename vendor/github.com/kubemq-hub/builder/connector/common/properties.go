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
	Values     map[string]string
	ValuesSpec string
}

func NewProperties(current map[string]string) *Properties {

	p := &Properties{
		Values: map[string]string{},
	}
	for key, val := range current {
		p.Values[key] = val
	}
	return p
}

func (p *Properties) Render() (map[string]string, error) {
	if len(p.Values) == 0 {
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
	} else {
		boolVal := false
		err := survey.NewBool().
			SetKind("bool").
			SetName("add-middleware").
			SetMessage("Would you to change middlewares to this binding").
			SetDefault("false").
			SetHelp("Change a middleware properties").
			SetRequired(true).
			Render(&boolVal)
		if err != nil {
			return nil, err
		}
		if !boolVal {
			return p.Values, nil
		}
	}

	if values, err := NewLog().Render(p.Values); err != nil {
		return nil, err
	} else {
		if values == nil {
			delete(p.Values, "log_level")
		} else {
			for key, val := range values {
				p.Values[key] = val
			}
		}

	}
	if values, err := NewRateLimiter().Render(p.Values); err != nil {
		return nil, err
	} else {
		if values == nil {
			delete(p.Values, "rate_per_second")
		}
		for key, val := range values {
			p.Values[key] = val
		}
	}

	if values, err := NewRetry().Render(p.Values); err != nil {
		return nil, err
	} else {
		if values == nil {
			delete(p.Values, "retry_delay_type")
			delete(p.Values, "retry_attempts")
			delete(p.Values, "retry_delay_milliseconds")
			delete(p.Values, "retry_max_jitter_milliseconds")
		}
		for key, val := range values {
			p.Values[key] = val
		}
	}
	return p.Values, nil
}
func (p *Properties) ColoredYaml() string {
	if len(p.Values) == 0 {
		return "\n<red>properties: {}</>"
	}
	p.ValuesSpec = utils.MapToYaml(p.Values)
	tpl := utils.NewTemplate(propertiesTml, p)
	b, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring properties spec,%s", err.Error())
	}
	return string(b)
}

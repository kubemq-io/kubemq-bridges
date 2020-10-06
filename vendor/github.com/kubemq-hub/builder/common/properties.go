package common

import "github.com/kubemq-hub/builder/survey"

type Properties struct {
	values map[string]string
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

package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type RateLimiter struct {
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

func (r *RateLimiter) Render(values map[string]string) (map[string]string, error) {
	var result map[string]string
	menu := survey.NewMenu("Select Rate Limiter Middleware Options:")

	addEditFn := func(values map[string]string) (map[string]string, error) {
		defaultRate := values["rate_per_second"]
		if defaultRate == "" {
			defaultRate = "100"
		}
		val := 0
		err := survey.NewInt().
			SetKind("int").
			SetName("rate-limiter").
			SetMessage("Set rate request per second limiting").
			SetDefault(defaultRate).
			SetHelp("Set how many request per second to limit").
			SetRequired(true).
			Render(&val)
		if err != nil {
			return nil, err
		}
		return map[string]string{"rate_per_second": fmt.Sprintf("%d", val)}, nil
	}
	_, ok := values["rate_per_second"]
	if ok {
		menu.AddItem("Edit Rate Limiter Middleware", func() error {
			var err error
			values, err := addEditFn(values)
			if err != nil {
				return err
			}
			result = values
			return nil
		})
		menu.AddItem("Remove Rate Limiter Middleware", func() error {
			result = nil
			return nil
		})
	} else {
		menu.AddItem("Add Rate Limiter Middleware", func() error {
			var err error
			resultValues, err := addEditFn(values)
			if err != nil {
				return err
			}
			result = resultValues
			return nil
		})
	}
	menu.SetDisableLoop(true)
	menu.SetBackOption(true)
	menu.SetErrorHandler(survey.MenuShowErrorFn)

	if err := menu.Render(); err != nil {
		return nil, err
	}
	return result, nil
}

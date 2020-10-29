package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
)

type Retry struct {
	values map[string]string
}

func NewRetry() *Retry {
	return &Retry{
		values: map[string]string{},
	}
}
func (r *Retry) askDelayType(values map[string]string) error {

	defaultValue := values["retry_delay_type"]
	if defaultValue == "" {
		defaultValue = "fixed"
	}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("retry_delay_type").
		SetMessage("Set retry delay type").
		SetOptions([]string{"fixed", "back-off", "random"}).
		SetDefault(defaultValue).
		SetHelp("Set retry delay type").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_delay_type"] = val
	return nil
}
func (r *Retry) askAttempts(values map[string]string) error {
	defaultValue := values["retry_attempts"]
	if defaultValue == "" {
		defaultValue = "1"
	}
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("retry_attempts").
		SetMessage("Set retry max attempts type").
		SetDefault(defaultValue).
		SetHelp("Set retry max attempts type").
		SetRequired(true).
		SetRange(1, 1024).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_attempts"] = fmt.Sprintf("%d", val)
	return nil
}
func (r *Retry) askDelayMillisecond(values map[string]string) error {
	defaultValue := values["retry_delay_milliseconds"]
	if defaultValue == "" {
		defaultValue = "1"
	}
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("retry_delay_milliseconds").
		SetMessage("Set retry delay milliseconds").
		SetDefault(defaultValue).
		SetHelp("Set retry delay milliseconds").
		SetRequired(true).
		SetRange(0, math.MaxInt32).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_delay_milliseconds"] = fmt.Sprintf("%d", val)
	return nil
}
func (r *Retry) askDelayJitter(values map[string]string) error {
	defaultValue := values["retry_max_jitter_milliseconds"]
	if defaultValue == "" {
		defaultValue = "100"
	}
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("retry_max_jitter_milliseconds").
		SetMessage("Set retry delay milliseconds jitter").
		SetDefault(defaultValue).
		SetHelp("Set retry delay milliseconds jitter").
		SetRequired(true).
		SetRange(1, math.MaxInt32).
		Render(&val)
	if err != nil {
		return err
	}
	r.values["retry_max_jitter_milliseconds"] = fmt.Sprintf("%d", val)
	return nil
}
func (r *Retry) Render(values map[string]string) (map[string]string, error) {
	var result map[string]string
	menu := survey.NewMenu("Select Retries Middleware Options:")

	addEditFn := func(values map[string]string) (map[string]string, error) {
		if err := r.askDelayType(values); err != nil {
			return nil, err
		}
		if err := r.askAttempts(values); err != nil {
			return nil, err
		}
		if err := r.askDelayMillisecond(values); err != nil {
			return nil, err
		}
		if err := r.askDelayJitter(values); err != nil {
			return nil, err
		}
		return r.values, nil
	}

	_, ok := values["retry_delay_type"]
	if ok {
		menu.AddItem("Edit Retries Middleware", func() error {
			var err error
			values, err := addEditFn(values)
			if err != nil {
				return err
			}
			result = values
			return nil
		})
		menu.AddItem("Remove Retries Middleware", func() error {
			result = nil
			return nil
		})
	} else {
		menu.AddItem("Add Retries Middleware", func() error {
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

package common

import (
	"github.com/kubemq-hub/builder/survey"
)

type Log struct {
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Render(values map[string]string) (map[string]string, error) {
	var result map[string]string
	menu := survey.NewMenu("Select Logging Middleware Options:")

	addEditFn := func(values map[string]string) (map[string]string, error) {
		defaultLogLevel := values["log_level"]
		if defaultLogLevel == "" {
			defaultLogLevel = "info"
		}
		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("log-level").
			SetMessage("Set Log level").
			SetOptions([]string{"debug", "info", "error"}).
			SetDefault(defaultLogLevel).
			SetHelp("Set Log level printing").
			SetRequired(true).
			Render(&val)
		if err != nil {
			return nil, err
		}
		return map[string]string{"log_level": val}, nil
	}

	_, ok := values["log_level"]
	if ok {
		menu.AddItem("Edit Logging Middleware", func() error {
			var err error
			values, err := addEditFn(values)
			if err != nil {
				return err
			}
			result = values
			return nil
		})
		menu.AddItem("Remove Logging Middleware", func() error {
			result = nil
			return nil
		})
	} else {
		menu.AddItem("Add Logging Middleware", func() error {
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

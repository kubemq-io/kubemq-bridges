package common

import "github.com/kubemq-hub/builder/survey"

type Log struct {
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Render() (map[string]string, error) {
	boolVal := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-log-middleware").
		SetMessage("Would you to set a logging middleware").
		SetDefault("true").
		SetHelp("Add logging middleware properties").
		SetRequired(true).
		Render(&boolVal)
	if err != nil {
		return nil, err
	}
	if !boolVal {
		return nil, nil
	}
	val := ""
	err = survey.NewString().
		SetKind("string").
		SetName("log-level").
		SetMessage("Set Log level").
		SetOptions([]string{"debug", "info", "error"}).
		SetDefault("info").
		SetHelp("Set Log level printing").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return nil, err
	}
	return map[string]string{"log_level": val}, nil
}

package source

import "github.com/kubemq-hub/builder/survey"

type Kind struct {
	defaultKind string
}

func NewKind(defaultKind string) *Kind {
	return &Kind{
		defaultKind: defaultKind,
	}

}

func (k *Kind) Render() (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Set Source kind").
		SetDefault(k.defaultKind).
		SetHelp("Set sources kind entry").
		SetRequired(true).
		SetOptions([]string{"kubemq.queue", "kubemq.query", "kubemq.events", "kubemq.events-store", "kubemq.command"}).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

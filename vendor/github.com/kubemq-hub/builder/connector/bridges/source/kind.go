package source

import "github.com/kubemq-hub/builder/survey"

type Kind struct {
}

func NewKind() *Kind {
	return &Kind{}
}

func (k *Kind) Render() (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Set Source kind").
		SetDefault("source.queue").
		SetHelp("Sets sources kind entry").
		SetRequired(true).
		SetOptions([]string{"source.queue", "source.events", "source.events-store", "source.command", "source.query"}).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

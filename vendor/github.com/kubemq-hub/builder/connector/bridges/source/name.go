package source

import "github.com/kubemq-hub/builder/survey"

type Name struct {
	takenNames  []string
	defaultName string
}

func NewName(defaultName string) *Name {
	return &Name{
		defaultName: defaultName,
	}
}
func (n *Name) SetTakenNames(value []string) *Name {
	n.takenNames = value
	return n
}
func (n *Name) Render() (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("name").
		SetMessage("Set Sources name").
		SetDefault(n.defaultName).
		SetHelp("Set sources name entry").
		SetRequired(true).
		SetValidator(survey.ValidateNoneSpace).
		SetInvalidOptions(n.takenNames).
		SetInvalidOptionsMessage("source name must be unique").
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

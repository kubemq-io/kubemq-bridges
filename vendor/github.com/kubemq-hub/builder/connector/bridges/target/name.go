package target

import "github.com/kubemq-hub/builder/survey"

type Name struct {
	takenNames []string
}

func NewName() *Name {
	return &Name{}
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
		SetMessage("Set Target name").
		SetDefault("").
		SetHelp("Sets targets name entry").
		SetRequired(true).
		SetInvalidOptions(n.takenNames).
		SetInvalidOptionsMessage("target name must be unique").
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

package binding

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
		SetMessage("Set Binding name").
		SetDefault("").
		SetHelp("Sets binding name entry").
		SetRequired(true).
		SetInvalidOptions(n.takenNames).
		SetInvalidOptionsMessage("binding name must be unique").
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

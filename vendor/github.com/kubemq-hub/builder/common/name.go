package common

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
func (n *Name) RenderBinding() (string, error) {
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
func (n *Name) RenderSource() (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("name").
		SetMessage("Set Source name").
		SetDefault("").
		SetHelp("Sets source name entry").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}
func (n *Name) RenderTarget() (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("name").
		SetMessage("Set Target name").
		SetDefault("").
		SetHelp("Sets Target name entry").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

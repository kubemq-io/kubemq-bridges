package common

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
func (n *Name) RenderBinding() (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("name").
		SetMessage("Set Binding name").
		SetDefault(n.defaultName).
		SetHelp("Set binding name entry").
		SetRequired(true).
		SetInvalidOptions(n.takenNames).
		SetValidator(survey.ValidateNoneSpace).
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
		SetDefault(n.defaultName).
		SetHelp("Set source name entry").
		SetRequired(true).
		SetValidator(survey.ValidateNoneSpace).
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
		SetDefault(n.defaultName).
		SetHelp("Set Target name entry").
		SetRequired(true).
		SetValidator(survey.ValidateNoneSpace).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

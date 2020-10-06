package connector

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Namespace struct {
	Namespace  string
	namespaces []string
}

func NewNamespace() *Namespace {
	return &Namespace{}
}
func (n *Namespace) Validate() error {
	if n.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	return nil
}
func (n *Namespace) SetNamespaces(value []string) *Namespace {
	n.namespaces = value
	return n
}
func (n *Namespace) checkNonEmptyNamespace(val interface{}) error {
	str, _ := val.(string)
	if str == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	return nil
}
func (n *Namespace) Render() (*Namespace, error) {

	err := survey.NewString().
		SetKind("string").
		SetName("namespace").
		SetMessage("Set connector namespace").
		SetDefault("").
		SetHelp("Sets connector namespace").
		SetRequired(true).
		SetValidator(n.checkNonEmptyNamespace).
		Render(&n.Namespace)
	if err != nil {
		return nil, err
	}
	return n, nil
}

var _ Validator = NewNamespace()

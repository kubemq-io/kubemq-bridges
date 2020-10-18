package target

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Target struct {
	Name           string              `json:"name"`
	Kind           string              `json:"kind"`
	Connections    []map[string]string `json:"connections"`
	ConnectionSpec string              `json:"-" yaml:"-"`
	addressOptions []string
	takenNames     []string
	defaultName    string
}

func NewTarget(defaultName string) *Target {
	return &Target{
		addressOptions: nil,
		defaultName:    defaultName,
	}
}
func (t *Target) SetAddress(value []string) *Target {
	t.addressOptions = value
	return t
}
func (t *Target) SetTakenNames(value []string) *Target {
	t.takenNames = value
	return t
}
func (t *Target) askAddConnection() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-connection").
		SetMessage("Would you like to add another target connection").
		SetDefault("false").
		SetHelp("Add new target connection").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false, err
	}
	return val, nil
}
func (t *Target) addConnection() error {
	if connection, err := NewConnection().
		SetAddress(t.addressOptions).
		Render(t.Name, t.Kind); err != nil {
		return err
	} else {
		t.Connections = append(t.Connections, connection)
	}
	return nil
}

func (t *Target) Render() (*Target, error) {
	var err error
	if t.Name, err = NewName(t.defaultName).
		SetTakenNames(t.takenNames).
		Render(); err != nil {
		return nil, err
	}
	if t.Kind, err = NewKind().
		Render(); err != nil {
		return nil, err
	}
	utils.Println(promptTargetFirstConnection, t.Kind)
	err = t.addConnection()
	if err != nil {
		return nil, err
	}
	for {
		addMore, err := t.askAddConnection()
		if err != nil {
			return t, nil
		}
		if addMore {
			err = t.addConnection()
			if err != nil {
				return nil, err
			}
		} else {
			goto done
		}
	}
done:
	return t, nil
}
func (t *Target) String() string {
	t.ConnectionSpec = utils.MapArrayToYaml(t.Connections)
	tpl := utils.NewTemplate(targetTemplate, t)
	b, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring target  spec,%s", err.Error())
	}
	return string(b)
}

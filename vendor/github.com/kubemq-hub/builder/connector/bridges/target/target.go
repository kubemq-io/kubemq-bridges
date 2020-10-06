package target

import (
	"github.com/kubemq-hub/builder/survey"
)

type Target struct {
	Name           string
	Kind           string
	Connections    []map[string]string
	addressOptions []string
	takenNames     []string
}

func NewTarget() *Target {
	return &Target{
		addressOptions: nil,
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
	if t.Name, err = NewName().
		SetTakenNames(t.takenNames).
		Render(); err != nil {
		return nil, err
	}
	if t.Kind, err = NewKind().
		Render(); err != nil {
		return nil, err
	}
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

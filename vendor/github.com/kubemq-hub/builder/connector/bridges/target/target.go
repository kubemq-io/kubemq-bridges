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
	WasEdited      bool                `json:"-" yaml:"-"`
	addressOptions []string
	takenNames     []string
	defaultName    string
	isEdit         bool
}

func NewTarget(defaultName string) *Target {
	return &Target{
		addressOptions: nil,
		defaultName:    defaultName,
	}
}
func (t *Target) Clone() *Target {
	newTarget := &Target{
		Name:           t.Name,
		Kind:           t.Kind,
		Connections:    []map[string]string{},
		ConnectionSpec: t.ConnectionSpec,
		addressOptions: t.addressOptions,
		takenNames:     t.takenNames,
		defaultName:    t.Name,
	}
	for _, connection := range t.Connections {
		newConnection := map[string]string{}
		for Key, val := range connection {
			newConnection[Key] = val
		}
		newTarget.Connections = append(newTarget.Connections, newConnection)
	}
	return newTarget
}
func (t *Target) SetIsEdit(value bool) *Target {
	t.isEdit = value
	return t
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
func (t *Target) add() (*Target, error) {
	var err error
	if t.Name, err = NewName(t.defaultName).
		SetTakenNames(t.takenNames).
		Render(); err != nil {
		return nil, err
	}
	if t.Kind, err = NewKind("").
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
func (t *Target) editName() error {
	var err error
	if t.Name, err = NewName(t.Name).
		SetTakenNames(t.takenNames).
		Render(); err != nil {
		return err
	}
	return nil
}
func (t *Target) editKind() (bool, error) {
	var err error
	current := t.Kind
	if t.Kind, err = NewKind(t.Kind).
		Render(); err != nil {
		return false, err
	}
	return t.Kind != current, nil
}
func (t *Target) editConnections() error {
	t.Connections = []map[string]string{}
	utils.Println(promptTargetFirstConnection, t.Kind)
	err := t.addConnection()
	if err != nil {
		return err
	}
	for {
		addMore, err := t.askAddConnection()
		if err != nil {
			return nil
		}
		if addMore {
			err = t.addConnection()
			if err != nil {
				return err
			}
		} else {
			goto done
		}
	}
done:
	return nil
}
func (t *Target) showConfiguration() error {
	utils.Println(promptShowTarget, t.Name)
	utils.Println(t.ColoredYaml())
	return nil
}
func (t *Target) edit() (*Target, error) {
	for {
		ops := []string{
			"Edit Targets name",
			"Edit Targets kind",
			"Edit Targets connections",
			"Show Targets configuration",
			"Return",
		}

		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-operation").
			SetMessage("Select Edit Binding Targets operation").
			SetDefault(ops[0]).
			SetHelp("Select Edit Binding Targets operation").
			SetRequired(true).
			SetOptions(ops).
			Render(&val)
		if err != nil {
			return nil, err
		}
		switch val {
		case ops[0]:
			if err := t.editName(); err != nil {
				return nil, err
			}
			t.WasEdited = true
		case ops[1]:
			if changed, err := t.editKind(); err != nil {
				return nil, err
			} else {
				if changed {
					if err := t.editConnections(); err != nil {
						return nil, err
					}
				}
			}
			t.WasEdited = true
		case ops[2]:
			if err := t.editConnections(); err != nil {
				return nil, err
			}
			t.WasEdited = true
		case ops[3]:
			if err := t.showConfiguration(); err != nil {
				return nil, err
			}

		default:
			return t, nil
		}
	}
}
func (t *Target) Render() (*Target, error) {
	if t.isEdit {
		return t.edit()
	} else {
		return t.add()
	}
}
func (t *Target) ColoredYaml() string {
	t.ConnectionSpec = utils.MapArrayToYaml(t.Connections)
	tpl := utils.NewTemplate(targetTemplate, t)
	b, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring target  spec,%s", err.Error())
	}
	return string(b)
}

func (t *Target) TableItemShort() string {
	return fmt.Sprintf("%s/%s/%d", t.Name, t.Kind, len(t.Connections))
}

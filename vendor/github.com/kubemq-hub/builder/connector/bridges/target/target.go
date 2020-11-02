package target

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Target struct {
	Kind           string              `json:"kind"`
	Connections    []map[string]string `json:"connections"`
	ConnectionSpec string              `json:"-" yaml:"-"`
	WasEdited      bool                `json:"-" yaml:"-"`
	isEdit         bool
	kubemqAddress  []string
}

func NewTarget() *Target {
	return &Target{}
}
func (t *Target) Clone() *Target {
	newTarget := &Target{
		Kind:           t.Kind,
		Connections:    []map[string]string{},
		ConnectionSpec: t.ConnectionSpec,
		WasEdited:      t.WasEdited,
		isEdit:         t.isEdit,
		kubemqAddress:  t.kubemqAddress,
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
func (t *Target) SetKubemqAddress(values []string) *Target {
	t.kubemqAddress = values
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
		SetAddress(t.kubemqAddress).
		Render(t.Kind); err != nil {
		return err
	} else {
		t.Connections = append(t.Connections, connection)
	}
	return nil
}
func (t *Target) add() (*Target, error) {
	var err error

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
	t.WasEdited = true
	return nil
}

func (t *Target) edit() (*Target, error) {
	var result *Target
	edited := t.Clone()
	form := survey.NewForm("Select Edit Target Option:")

	ftKind := new(string)
	*ftKind = fmt.Sprintf("<k> Edit Target Kind (%s)", edited.Kind)
	ftConnections := new(string)
	*ftConnections = fmt.Sprintf("<c> Edit Target Connections (%s)", edited.Kind)

	form.AddItem(ftKind, func() error {
		if changed, err := edited.editKind(); err != nil {
			return err
		} else {
			if changed {
				if err := edited.editConnections(); err != nil {
					return err
				}
				edited.WasEdited = true
			}
		}
		*ftKind = fmt.Sprintf("<k> Edit Target Kind (%s)", edited.Kind)
		*ftConnections = fmt.Sprintf("<c> Edit Target Connections (%s)", edited.Kind)
		return nil
	})

	form.AddItem(ftConnections, func() error {
		if err := edited.editConnections(); err != nil {
			return err
		}
		*ftConnections = fmt.Sprintf("<c> Edit Target Connections (%s)", edited.Kind)
		return nil
	})

	form.AddItem("<s> Show Target Configuration", func() error {
		utils.Println(promptShowTarget)
		utils.Println("%s\n", edited.ColoredYaml())
		return nil
	})
	form.SetOnSaveFn(func() error {
		if err := edited.Validate(); err != nil {
			return err
		}
		result = edited
		return nil
	})

	form.SetOnCancelFn(func() error {
		result = t
		return nil
	})
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
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
	return fmt.Sprintf("%s/%d", t.Kind, len(t.Connections))
}

func (t *Target) Validate() error {
	return nil
}

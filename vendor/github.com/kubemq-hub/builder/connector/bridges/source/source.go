package source

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Source struct {
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

func NewSource(defaultName string) *Source {
	return &Source{
		addressOptions: nil,
		defaultName:    defaultName,
	}
}
func (s *Source) Clone() *Source {
	newSrc := &Source{
		Name:           s.Name,
		Kind:           s.Kind,
		Connections:    []map[string]string{},
		ConnectionSpec: s.ConnectionSpec,
		addressOptions: s.addressOptions,
		takenNames:     s.takenNames,
		defaultName:    s.Name,
	}
	for _, connection := range s.Connections {
		newConnection := map[string]string{}
		for Key, val := range connection {
			newConnection[Key] = val
		}
		newSrc.Connections = append(newSrc.Connections, newConnection)
	}
	return newSrc
}

func (s *Source) SetAddress(value []string) *Source {
	s.addressOptions = value
	return s
}
func (s *Source) SetTakenNames(value []string) *Source {
	s.takenNames = value
	return s
}
func (s *Source) SetIsEdit(value bool) *Source {
	s.isEdit = value
	return s
}
func (s *Source) askAddConnection() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-connection").
		SetMessage("Would you like to add another source connection").
		SetDefault("false").
		SetHelp("Add new source connection").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false, err
	}
	return val, nil
}

func (s *Source) addConnection() error {
	if connection, err := NewConnection().
		SetAddress(s.addressOptions).
		Render(s.Name, s.Kind); err != nil {
		return err
	} else {
		s.Connections = append(s.Connections, connection)
	}
	return nil
}
func (s *Source) add() (*Source, error) {
	var err error
	if s.Name, err = NewName(s.defaultName).
		SetTakenNames(s.takenNames).
		Render(); err != nil {
		return nil, err
	}
	if s.Kind, err = NewKind("").
		Render(); err != nil {
		return nil, err
	}
	utils.Println(promptSourceFirstConnection, s.Kind)
	err = s.addConnection()
	if err != nil {
		return nil, err
	}
	for {
		addMore, err := s.askAddConnection()
		if err != nil {
			return s, nil
		}
		if addMore {
			err = s.addConnection()
			if err != nil {
				return nil, err
			}
		} else {
			goto done
		}
	}
done:
	return s, nil
}
func (s *Source) editName() error {
	var err error
	if s.Name, err = NewName(s.Name).
		SetTakenNames(s.takenNames).
		Render(); err != nil {
		return err
	}
	return nil
}
func (s *Source) editKind() (bool, error) {
	var err error
	current := s.Kind
	if s.Kind, err = NewKind(s.Kind).
		Render(); err != nil {
		return false, err
	}
	return s.Kind != current, nil
}
func (s *Source) editConnections() error {
	s.Connections = []map[string]string{}
	utils.Println(promptSourceFirstConnection, s.Kind)
	err := s.addConnection()
	if err != nil {
		return err
	}
	for {
		addMore, err := s.askAddConnection()
		if err != nil {
			return nil
		}
		if addMore {
			err = s.addConnection()
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
func (s *Source) showConfiguration() error {
	utils.Println(promptShowSource, s.Name)
	utils.Println(s.ColoredYaml())
	return nil
}
func (s *Source) edit() (*Source, error) {
	for {
		ops := []string{
			"Edit Sources name",
			"Edit Sources kind",
			"Edit Sources connections",
			"Show Sources configuration",
			"Return",
		}

		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-operation").
			SetMessage("Select Edit Binding Sources operation").
			SetDefault(ops[0]).
			SetHelp("Select Edit Binding Sources operation").
			SetRequired(true).
			SetOptions(ops).
			Render(&val)
		if err != nil {
			return nil, err
		}
		switch val {
		case ops[0]:
			if err := s.editName(); err != nil {
				return nil, err
			}
			s.WasEdited = true
		case ops[1]:
			if changed, err := s.editKind(); err != nil {
				return nil, err
			} else {
				if changed {
					if err := s.editConnections(); err != nil {
						return nil, err
					}
				}
			}
			s.WasEdited = true
		case ops[2]:
			if err := s.editConnections(); err != nil {
				return nil, err
			}
			s.WasEdited = true
		case ops[3]:
			if err := s.showConfiguration(); err != nil {
				return nil, err
			}

		default:
			return s, nil
		}
	}
}
func (s *Source) Render() (*Source, error) {
	if s.isEdit {
		return s.edit()
	}
	return s.add()
}

func (s *Source) ColoredYaml() string {
	s.ConnectionSpec = utils.MapArrayToYaml(s.Connections)
	t := utils.NewTemplate(sourceTemplate, s)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring source  spec,%s", err.Error())
	}
	return string(b)
}
func (s *Source) TableItemShort() string {
	return fmt.Sprintf("%s/%s/%d", s.Name, s.Kind, len(s.Connections))
}

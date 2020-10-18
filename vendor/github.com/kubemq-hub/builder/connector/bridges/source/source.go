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
	addressOptions []string
	takenNames     []string
	defaultName    string
}

func NewSource(defaultName string) *Source {
	return &Source{
		addressOptions: nil,
		defaultName:    defaultName,
	}
}
func (s *Source) SetAddress(value []string) *Source {
	s.addressOptions = value
	return s
}
func (s *Source) SetTakenNames(value []string) *Source {
	s.takenNames = value
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

func (s *Source) Render() (*Source, error) {
	var err error
	if s.Name, err = NewName(s.defaultName).
		SetTakenNames(s.takenNames).
		Render(); err != nil {
		return nil, err
	}
	if s.Kind, err = NewKind().
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

func (s *Source) String() string {
	s.ConnectionSpec = utils.MapArrayToYaml(s.Connections)
	t := utils.NewTemplate(sourceTemplate, s)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring source  spec,%s", err.Error())
	}
	return string(b)
}

package source

import (
	"github.com/kubemq-hub/builder/survey"
)

type Source struct {
	Name           string
	Kind           string
	Connections    []map[string]string
	addressOptions []string
	takenNames     []string
}

func NewSource() *Source {
	return &Source{
		addressOptions: nil,
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
	if s.Name, err = NewName().
		SetTakenNames(s.takenNames).
		Render(); err != nil {
		return nil, err
	}
	if s.Kind, err = NewKind().
		Render(); err != nil {
		return nil, err
	}
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

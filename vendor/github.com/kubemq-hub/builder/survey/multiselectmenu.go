package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

type MultiSelectMenu struct {
	title        string
	fnMap        map[string]func() error
	fnItems      []string
	errorHandler func(err error) error
}

func NewMultiSelectMenu(title string) *MultiSelectMenu {
	return &MultiSelectMenu{
		title:   title,
		fnMap:   map[string]func() error{},
		fnItems: []string{},
	}
}

func (m *MultiSelectMenu) SetErrorHandler(value func(err error) error) *MultiSelectMenu {
	m.errorHandler = value
	return m
}

func (m *MultiSelectMenu) AddItem(title string, fn func() error) *MultiSelectMenu {
	m.fnMap[title] = fn
	m.fnItems = append(m.fnItems, title)
	return m
}

func (m *MultiSelectMenu) Render() error {
	if len(m.fnItems) == 0 {
		return fmt.Errorf("no items to select are available")
	}
	m.AddItem("<cancel>", nil)
	if m.errorHandler == nil {
		m.errorHandler = MenuShowErrorFn
	}
	itemsLength := len(m.fnItems) + 1
	pageSize := 7
	if itemsLength > pageSize {
		pageSize = len(m.fnItems) + 1
	}
	if pageSize > 25 {
		pageSize = 25
	}
	var values []string
	menu := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       m.title,
		Options:       m.fnItems,
		Default:       m.fnItems[0],
		PageSize:      pageSize,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}

	err := survey.AskOne(menu, &values)
	if err != nil {
		return err
	}
	for _, value := range values {
		if value == "<cancel>" {
			return nil
		}
	}

	for _, val := range values {
		fn, ok := m.fnMap[val]
		if !ok {
			return fmt.Errorf("menu function for %s not found", val)
		}
		if fn == nil {
			continue
		}
		if err := fn(); err != nil {
			err := m.errorHandler(err)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

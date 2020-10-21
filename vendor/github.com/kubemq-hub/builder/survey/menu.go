package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-hub/builder/pkg/utils"
)

type Menu struct {
	title        string
	fnMap        map[string]func() error
	fnItems      []string
	keepFilter   bool
	pageSize     int
	errorHandler func(err error) error
	disableLoop  bool
	back         bool
}

func NewMenu(title string) *Menu {
	return &Menu{
		title:   title,
		fnMap:   map[string]func() error{},
		fnItems: []string{},
	}
}
func (m *Menu) SetKeepFilter(value bool) *Menu {
	m.keepFilter = value
	return m
}
func (m *Menu) SetPageSize(value int) *Menu {
	m.pageSize = value
	return m
}
func (m *Menu) SetErrorHandler(value func(err error) error) *Menu {
	m.errorHandler = value
	return m
}
func (m *Menu) SetDisableLoop(value bool) *Menu {
	m.disableLoop = value
	return m
}
func (m *Menu) SetBackOption(value bool) *Menu {
	m.back = value
	return m
}
func (m *Menu) AddItem(title string, fn func() error) *Menu {
	m.fnMap[title] = fn
	m.fnItems = append(m.fnItems, title)
	return m
}

func (m *Menu) Render() error {
	if m.back {
		m.AddItem("<-back", nil)
	}
	val := ""
	menu := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       m.title,
		Options:       m.fnItems,
		Default:       m.fnItems[0],
		PageSize:      m.pageSize,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	for {
		err := survey.AskOne(menu, &val)
		if err != nil {
			return err
		}

		fn, ok := m.fnMap[val]
		if !ok {
			return fmt.Errorf("menu function for %s not found", val)
		}
		if fn == nil {
			return nil
		}
		if err := fn(); err != nil {
			if m.errorHandler != nil {
				err := m.errorHandler(err)
				if err != nil {
					return err
				}
				goto loop
			} else {
				return err
			}
		}
	loop:
		if m.disableLoop {
			return nil
		}
	}
}

func MenuShowErrorFn(err error) error {
	utils.Println("<red>%s</>", err.Error())
	return nil
}

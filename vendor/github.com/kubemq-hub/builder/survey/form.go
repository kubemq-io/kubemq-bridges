package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-hub/builder/pkg/utils"
)

const (
	FormSave     = "<save>"
	FormCancel   = "<cancel>"
	FormDefaults = "<defaults>"
)

type Form struct {
	title     string
	fnMap     map[int]func() error
	fnItems   []interface{}
	pageSize  int
	onError   func(err error) error
	onSave    func() error
	onDefault func() error
	onCancel  func() error
}

func NewForm(title string) *Form {
	return &Form{
		title:   title,
		fnMap:   map[int]func() error{},
		fnItems: []interface{}{},
	}
}
func (f *Form) SetOnErrorFn(fn func(err error) error) *Form {
	f.onError = fn
	return f
}
func (f *Form) SetOnSaveFn(fn func() error) *Form {
	f.onSave = fn
	return f
}
func (f *Form) SetOnDefaultFn(fn func() error) *Form {
	f.onDefault = fn
	return f
}
func (f *Form) SetOnCancelFn(fn func() error) *Form {
	f.onCancel = fn
	return f
}

func (f *Form) SetPageSize(value int) *Form {
	f.pageSize = value
	return f
}
func (f *Form) AddItem(title interface{}, fn func() error) *Form {
	f.fnMap[len(f.fnItems)] = fn
	f.fnItems = append(f.fnItems, title)
	return f
}
func (f *Form) buildOptions() ([]string, map[string]int) {
	m := map[string]int{}
	var list []string
	for i, item := range f.fnItems {
		switch v := item.(type) {
		case string:
			list = append(list, v)
			m[v] = i
		case *string:
			list = append(list, *v)
			m[*v] = i
		}
	}
	return list, m
}
func (f *Form) Render() error {

	if f.onSave != nil {
		f.AddItem(FormSave, f.onSave)
	}
	if f.onDefault != nil {
		f.AddItem(FormDefaults, f.onDefault)
	}
	if f.onCancel != nil {
		f.AddItem(FormCancel, f.onCancel)
	}
	lastIndex := 0
	for {
		options, selectionMap := f.buildOptions()
		val := ""
		menu := &survey.Select{
			Renderer:      survey.Renderer{},
			Message:       f.title,
			Options:       options,
			Default:       options[lastIndex],
			PageSize:      f.pageSize,
			VimMode:       false,
			FilterMessage: "",
			Filter:        nil,
		}
		err := survey.AskOne(menu, &val)
		if err != nil {
			return err
		}
		lastIndex = selectionMap[val]
		fn, ok := f.fnMap[lastIndex]
		if !ok {
			return fmt.Errorf("form function for %s not found", val)
		}
		if fn == nil {
			return nil
		}

		switch val {
		case FormSave:
			err := fn()
			if err == nil {
				return nil
			}
			if f.onError != nil {
				err := f.onError(err)
				if err != nil {
					return err
				}
				continue
			} else {
				return err
			}
		case FormCancel:
			_ = fn()
			return nil
		case FormDefaults:
			err := fn()
			if err != nil {
				if f.onError != nil {
					err := f.onError(err)
					if err != nil {
						return err
					}
					continue
				} else {
					return err
				}
			} else {
				continue
			}
		default:
			err := fn()
			if err == nil {
				continue
			}
			if f.onError != nil {
				err := f.onError(err)
				if err != nil {
					return err
				}
				continue
			} else {
				return err
			}
		}
	}
}

func FormShowErrorFn(err error) error {
	utils.Println("<red>%s</>", err.Error())
	return nil
}

package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

type List struct {
	*KindMeta
	*ObjectMeta
	Options []string

	askOpts    []survey.AskOpt
	validators []func(val interface{}) error
	pageSize   int
}

func NewList() *List {
	return &List{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
		Options:    nil,
		askOpts:    nil,
	}
}

func (l *List) NewKindMeta() *List {
	l.KindMeta = NewKindMeta()
	return l
}
func (l *List) NewObjectMeta() *List {
	l.ObjectMeta = NewObjectMeta()
	return l
}
func (l *List) SetKind(value string) *List {
	l.KindMeta.SetKind(value)
	return l
}
func (l *List) SetPageSize(value int) *List {
	l.pageSize = value
	return l
}
func (l *List) SetName(value string) *List {
	l.ObjectMeta.SetName(value)
	return l
}

func (l *List) SetMessage(value string) *List {
	l.ObjectMeta.SetMessage(value)
	return l
}

func (l *List) SetDefault(value string) *List {
	l.ObjectMeta.SetDefault(value)
	return l
}

func (l *List) SetHelp(value string) *List {
	l.ObjectMeta.SetHelp(value)
	return l
}

func (l *List) SetRequired(value bool) *List {
	l.ObjectMeta.SetRequired(value)
	return l
}

func (l *List) SetOptions(value []string) *List {
	l.Options = value
	return l
}

func (l *List) SetValidator(f func(val interface{}) error) *List {
	l.validators = append(l.validators, f)
	return l
}

func (l *List) Complete() error {
	if err := l.KindMeta.complete(); err != nil {
		return err
	}
	l.askOpts = append(l.askOpts, l.KindMeta.askOpts...)

	if err := l.ObjectMeta.complete(); err != nil {
		return err
	}
	if len(l.Options) == 0 {
		return fmt.Errorf("no options to select")
	}
	l.askOpts = append(l.askOpts, l.ObjectMeta.askOpts...)
	for _, validator := range l.validators {
		l.askOpts = append(l.askOpts, survey.WithValidator(validator))
	}
	return nil
}

func (l *List) Render(target interface{}) error {
	if err := l.Complete(); err != nil {
		return err
	}
	selectInput := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       l.Message,
		Options:       l.Options,
		Default:       l.Default,
		Help:          l.Help,
		PageSize:      l.pageSize,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err := survey.AskOne(selectInput, target, l.askOpts...)
	if err != nil {
		return err
	}
	return nil
}

var _ Question = NewList()

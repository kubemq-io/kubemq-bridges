package survey

import (
	"github.com/AlecAivazis/survey/v2"
)

type Multiline struct {
	*KindMeta
	*ObjectMeta
	askOpts    []survey.AskOpt
	validators []func(val interface{}) error
}

func NewMultiline() *Multiline {
	return &Multiline{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
		askOpts:    nil,
	}
}

func (e *Multiline) NewKindMeta() *Multiline {
	e.KindMeta = NewKindMeta()
	return e
}
func (e *Multiline) NewObjectMeta() *Multiline {
	e.ObjectMeta = NewObjectMeta()
	return e
}
func (e *Multiline) SetKind(value string) *Multiline {
	e.KindMeta.SetKind(value)
	return e
}

func (e *Multiline) SetName(value string) *Multiline {
	e.ObjectMeta.SetName(value)
	return e
}

func (e *Multiline) SetMessage(value string) *Multiline {
	e.ObjectMeta.SetMessage(value)
	return e
}

func (e *Multiline) SetDefault(value string) *Multiline {
	e.ObjectMeta.SetDefault(value)
	return e
}

func (e *Multiline) SetHelp(value string) *Multiline {
	e.ObjectMeta.SetHelp(value)
	return e
}

func (e *Multiline) SetRequired(value bool) *Multiline {
	e.ObjectMeta.SetRequired(value)
	return e
}

func (e *Multiline) SetValidator(f func(val interface{}) error) *Multiline {
	e.validators = append(e.validators, f)
	return e
}

func (e *Multiline) Complete() error {
	if err := e.KindMeta.complete(); err != nil {
		return err
	}
	e.askOpts = append(e.askOpts, e.KindMeta.askOpts...)

	if err := e.ObjectMeta.complete(); err != nil {
		return err
	}
	e.askOpts = append(e.askOpts, e.ObjectMeta.askOpts...)

	for _, validator := range e.validators {
		e.askOpts = append(e.askOpts, survey.WithValidator(validator))
	}
	return nil
}

func (e *Multiline) Render(target interface{}) error {
	if err := e.Complete(); err != nil {
		return err
	}
	selectInput := &survey.Multiline{
		Renderer: survey.Renderer{},
		Message:  e.Message,
		Default:  e.Default,
		Help:     e.Help,
	}
	err := survey.AskOne(selectInput, target, e.askOpts...)
	if err != nil {
		return err
	}
	return nil
}

var _ Question = NewMultiline()

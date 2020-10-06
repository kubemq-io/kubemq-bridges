package survey

import (
	"github.com/AlecAivazis/survey/v2"
)

type Editor struct {
	*KindMeta
	*ObjectMeta
	Pattern    string
	askOpts    []survey.AskOpt
	validators []func(val interface{}) error
}

func NewEditor() *Editor {
	return &Editor{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
		askOpts:    nil,
	}
}

func (e *Editor) NewKindMeta() *Editor {
	e.KindMeta = NewKindMeta()
	return e
}
func (e *Editor) NewObjectMeta() *Editor {
	e.ObjectMeta = NewObjectMeta()
	return e
}
func (e *Editor) SetKind(value string) *Editor {
	e.KindMeta.SetKind(value)
	return e
}

func (e *Editor) SetName(value string) *Editor {
	e.ObjectMeta.SetName(value)
	return e
}

func (e *Editor) SetMessage(value string) *Editor {
	e.ObjectMeta.SetMessage(value)
	return e
}

func (e *Editor) SetDefault(value string) *Editor {
	e.ObjectMeta.SetDefault(value)
	return e
}

func (e *Editor) SetHelp(value string) *Editor {
	e.ObjectMeta.SetHelp(value)
	return e
}
func (e *Editor) SetPattern(value string) *Editor {
	e.Pattern = value
	return e
}
func (e *Editor) SetRequired(value bool) *Editor {
	e.ObjectMeta.SetRequired(value)
	return e
}

func (e *Editor) SetValidator(f func(val interface{}) error) *Editor {
	e.validators = append(e.validators, f)
	return e
}

func (e *Editor) Complete() error {
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

func (e *Editor) Render(target interface{}) error {
	if err := e.Complete(); err != nil {
		return err
	}
	selectInput := &survey.Editor{
		Renderer:      survey.Renderer{},
		Message:       e.Message,
		Default:       e.Default,
		Help:          e.Help,
		Editor:        "",
		HideDefault:   false,
		AppendDefault: false,
		FileName:      e.Pattern,
	}

	err := survey.AskOne(selectInput, target, e.askOpts...)
	if err != nil {
		return err
	}
	return nil
}

var _ Question = NewEditor()

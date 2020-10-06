package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strconv"
)

type Int struct {
	*KindMeta
	*ObjectMeta
	InvalidOptions        []string
	InvalidOptionsMessage string
	Range                 bool
	Min                   int
	Max                   int
	askOpts               []survey.AskOpt
	validators            []func(val interface{}) error
}

func NewInt() *Int {
	return &Int{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
		Range:      false,
		Min:        0,
		Max:        0,
		askOpts:    nil,
	}
}

func (i *Int) NewKindMeta() *Int {
	i.KindMeta = NewKindMeta()
	return i
}
func (i *Int) NewObjectMeta() *Int {
	i.ObjectMeta = NewObjectMeta()
	return i
}
func (i *Int) SetKind(value string) *Int {
	i.KindMeta.SetKind(value)
	return i
}

func (i *Int) SetName(value string) *Int {
	i.ObjectMeta.SetName(value)
	return i
}

func (i *Int) SetMessage(value string) *Int {
	i.ObjectMeta.SetMessage(value)
	return i
}

func (i *Int) SetDefault(value string) *Int {
	i.ObjectMeta.SetDefault(value)
	return i
}

func (i *Int) SetHelp(value string) *Int {
	i.ObjectMeta.SetHelp(value)
	return i
}
func (i *Int) SetInvalidOptionsMessage(value string) *Int {
	i.InvalidOptionsMessage = value
	return i
}
func (i *Int) SetInvalidOptions(value []string) *Int {
	i.InvalidOptions = value
	return i
}
func (i *Int) SetRequired(value bool) *Int {
	i.ObjectMeta.SetRequired(value)
	return i
}

func (i *Int) SetRange(min, max int) *Int {
	i.Range = true
	i.Min = min
	i.Max = max
	return i
}
func (i *Int) SetValidator(f func(val interface{}) error) *Int {
	i.validators = append(i.validators, f)
	return i
}

func (i *Int) checkValue(val interface{}) error {
	if str, ok := val.(string); ok {
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		if i.Range {
			if val < i.Min {
				return fmt.Errorf("value cannot be lower than minimum %d", i.Min)
			}
			if val > i.Max {
				return fmt.Errorf("value cannot be higher than maximum %d", i.Max)
			}
		}

	}
	return nil
}
func (i *Int) invalidOptionValidator(val interface{}) error {
	if str, ok := val.(string); ok {
		for _, item := range i.InvalidOptions {
			if str == item {
				return fmt.Errorf("%s", i.InvalidOptionsMessage)
			}
		}
	}
	return nil
}
func (i *Int) Complete() error {
	if err := i.KindMeta.complete(); err != nil {
		return err
	}
	i.askOpts = append(i.askOpts, i.KindMeta.askOpts...)

	if err := i.ObjectMeta.complete(); err != nil {
		return err
	}
	i.askOpts = append(i.askOpts, i.ObjectMeta.askOpts...)

	i.askOpts = append(i.askOpts, survey.WithValidator(i.checkValue))
	if i.InvalidOptionsMessage == "" {
		i.InvalidOptionsMessage = "invalid option,"
	}
	if len(i.InvalidOptions) > 0 {
		i.askOpts = append(i.askOpts, survey.WithValidator(i.invalidOptionValidator))
	}
	for _, validator := range i.validators {
		i.askOpts = append(i.askOpts, survey.WithValidator(validator))
	}
	return nil
}

func (i *Int) Render(target interface{}) error {
	if err := i.Complete(); err != nil {
		return err
	}
	singleInput := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  i.Message,
		Default:  i.Default,
		Help:     i.Help,
	}
	return survey.AskOne(singleInput, target, i.askOpts...)
}

var _ Question = NewInt()

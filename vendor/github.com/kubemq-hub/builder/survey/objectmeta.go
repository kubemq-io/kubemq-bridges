package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

type ObjectMeta struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	Default  string `json:"default"`
	Help     string `json:"help"`
	Required bool   `json:"required"`
	askOpts  []survey.AskOpt
}

func NewObjectMeta() *ObjectMeta {
	return &ObjectMeta{}
}

func (o *ObjectMeta) SetName(value string) *ObjectMeta {
	o.Name = value
	return o
}

func (o *ObjectMeta) SetMessage(value string) *ObjectMeta {
	o.Message = value
	return o
}

func (o *ObjectMeta) SetDefault(value string) *ObjectMeta {
	o.Default = value
	return o
}

func (o *ObjectMeta) SetHelp(value string) *ObjectMeta {
	o.Help = value
	return o
}

func (o *ObjectMeta) SetRequired(value bool) *ObjectMeta {
	o.Required = value
	return o
}
func (o *ObjectMeta) complete() error {
	if o.Name == "" {
		return fmt.Errorf("input must have a name")
	}
	if o.Message == "" {
		return fmt.Errorf("input must have a message")
	}
	if o.Required {
		o.askOpts = append(o.askOpts, survey.WithValidator(survey.Required))
	}
	return nil
}

package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

type String struct {
	*KindMeta
	*ObjectMeta
	Options               []string
	InvalidOptions        []string
	InvalidOptionsMessage string
	askOpts               []survey.AskOpt
	validators            []func(val interface{}) error
	keepFilter            bool
	pageSize              int
}

func NewString() *String {
	return &String{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
		Options:    nil,
		askOpts:    nil,
	}
}

func (s *String) NewKindMeta() *String {
	s.KindMeta = NewKindMeta()
	return s
}
func (s *String) NewObjectMeta() *String {
	s.ObjectMeta = NewObjectMeta()
	return s
}
func (s *String) SetKind(value string) *String {
	s.KindMeta.SetKind(value)
	return s
}
func (s *String) SetKeepFilter(value bool) *String {
	s.keepFilter = value
	return s
}
func (s *String) SetPageSize(value int) *String {
	s.pageSize = value
	return s
}

func (s *String) SetName(value string) *String {
	s.ObjectMeta.SetName(value)
	return s
}

func (s *String) SetMessage(value string) *String {
	s.ObjectMeta.SetMessage(value)
	return s
}

func (s *String) SetDefault(value string) *String {
	s.ObjectMeta.SetDefault(value)
	return s
}

func (s *String) SetHelp(value string) *String {
	s.ObjectMeta.SetHelp(value)
	return s
}
func (s *String) SetInvalidOptionsMessage(value string) *String {
	s.InvalidOptionsMessage = value
	return s
}
func (s *String) SetInvalidOptions(value []string) *String {
	s.InvalidOptions = value
	return s
}
func (s *String) SetRequired(value bool) *String {
	s.ObjectMeta.SetRequired(value)
	return s
}

func (s *String) SetOptions(value []string) *String {
	s.Options = value
	return s
}

func (s *String) SetValidator(f func(val interface{}) error) *String {
	s.validators = append(s.validators, f)
	return s
}

func (s *String) invalidOptionValidator(val interface{}) error {
	if str, ok := val.(string); ok {
		for _, item := range s.InvalidOptions {
			if str == item {
				return fmt.Errorf("%s", s.InvalidOptionsMessage)
			}
		}
	}
	return nil
}
func (s *String) Complete() error {
	if err := s.KindMeta.complete(); err != nil {
		return err
	}
	s.askOpts = append(s.askOpts, s.KindMeta.askOpts...)

	if err := s.ObjectMeta.complete(); err != nil {
		return err
	}
	s.askOpts = append(s.askOpts, s.ObjectMeta.askOpts...)

	if s.InvalidOptionsMessage == "" {
		s.InvalidOptionsMessage = "invalid option,"
	}
	if len(s.InvalidOptions) > 0 {
		s.askOpts = append(s.askOpts, survey.WithValidator(s.invalidOptionValidator))
	}
	for _, validator := range s.validators {
		s.askOpts = append(s.askOpts, survey.WithValidator(validator))
	}
	s.askOpts = append(s.askOpts, survey.WithKeepFilter(true))
	if s.pageSize > 0 && len(s.Options) >= s.pageSize {
		s.askOpts = append(s.askOpts, survey.WithPageSize(s.pageSize))
	}
	return nil
}

func (s *String) Render(target interface{}) error {
	if err := s.Complete(); err != nil {
		return err
	}
	if len(s.Options) == 0 {
		singleInput := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  s.Message,
			Default:  s.Default,
			Help:     s.Help,
		}
		return survey.AskOne(singleInput, target, s.askOpts...)
	}

	if s.Default == "" {
		s.Default = s.Options[0]
	}
	selectInput := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       s.Message,
		Options:       s.Options,
		Default:       s.Default,
		Help:          s.Help,
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}

	err := survey.AskOne(selectInput, target, s.askOpts...)
	if err != nil {
		return err
	}
	val, _ := target.(*string)
	if *val == "Other" {
		singleInput := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  fmt.Sprintf("%s, Other", s.Message),
			Default:  "",
			Help:     s.Help,
		}
		return survey.AskOne(singleInput, target, s.askOpts...)
	}
	return nil
}

var _ Question = NewString()

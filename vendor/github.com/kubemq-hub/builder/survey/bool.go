package survey

import "github.com/AlecAivazis/survey/v2"

type Bool struct {
	*KindMeta
	*ObjectMeta
	askOpts []survey.AskOpt
}

func (c *Bool) NewKindMeta() *Bool {
	c.KindMeta = NewKindMeta()
	return c
}
func (c *Bool) NewObjectMeta() *Bool {
	c.ObjectMeta = NewObjectMeta()
	return c
}
func (c *Bool) SetKind(value string) *Bool {
	c.KindMeta.SetKind(value)
	return c
}

func (c *Bool) SetName(value string) *Bool {
	c.ObjectMeta.SetName(value)
	return c
}

func (c *Bool) SetMessage(value string) *Bool {
	c.ObjectMeta.SetMessage(value)
	return c
}

func (c *Bool) SetDefault(value string) *Bool {
	c.ObjectMeta.SetDefault(value)
	return c
}

func (c *Bool) SetHelp(value string) *Bool {
	c.ObjectMeta.SetHelp(value)
	return c
}
func (c *Bool) SetRequired(value bool) *Bool {
	c.ObjectMeta.SetRequired(value)
	return c
}

func (c *Bool) Complete() error {
	return nil
}

func (c *Bool) Render(target interface{}) error {
	if err := c.Complete(); err != nil {
		return err
	}
	defValue := false
	if c.Default == "true" {
		defValue = true
	}
	boolVal := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  c.Message,
		Default:  defValue,
		Help:     c.Help,
	}
	return survey.AskOne(boolVal, target, c.askOpts...)
}

func NewBool() *Bool {
	return &Bool{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
	}
}

var _ Question = NewBool()

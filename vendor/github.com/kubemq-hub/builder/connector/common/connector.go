package common

import (
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"strings"
)

type Connector struct {
	Kind           string      `json:"kind"`
	Description    string      `json:"description"`
	Properties     []*Property `json:"properties"`
	PropertiesSpec string
	loadedOptions  DefaultOptions
	values         map[string]string
}

func NewConnector() *Connector {
	return &Connector{
		Kind:          "",
		Description:   "",
		Properties:    nil,
		loadedOptions: nil,
	}
}

func (c *Connector) SetKind(value string) *Connector {
	c.Kind = value
	return c
}

func (c *Connector) SetDescription(value string) *Connector {
	c.Description = value
	return c
}
func (c *Connector) AddProperty(value *Property) *Connector {
	c.Properties = append(c.Properties, value)
	return c
}
func (c *Connector) Validate() error {
	if c.Kind == "" {
		return fmt.Errorf("connector kind cannot be empty")
	}
	if c.Description == "" {
		return fmt.Errorf("connector description cannot be empty")
	}
	if len(c.Properties) == 0 {
		return fmt.Errorf("connector must have at least one property")
	}
	for _, property := range c.Properties {
		if err := property.Validate(); err != nil {
			return err
		}
	}
	return nil
}
func (c *Connector) askString(p *Property) error {
	val := ""
	options := p.Options
	loaded, ok := c.loadedOptions[p.LoadedOptions]
	if ok {
		options = append(options, loaded...)
	}
	err := survey.NewString().
		SetKind("string").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetOptions(options).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	if val != "" {
		c.values[p.Name] = val
	}
	return nil
}
func (c *Connector) askInt(p *Property) error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetHelp(p.Description).
		SetRequired(p.Must).
		SetRange(p.Min, p.Max).
		Render(&val)
	if err != nil {
		return err
	}
	c.values[p.Name] = fmt.Sprintf("%d", val)
	return nil
}
func (c *Connector) askBool(p *Property) error {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	c.values[p.Name] = fmt.Sprintf("%t", val)
	return nil
}
func (c *Connector) askMultilines(p *Property) error {
	val := ""
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	c.values[p.Name] = val
	return nil
}
func (c *Connector) askMap(p *Property) error {
	values := map[string]string{}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetOptions(p.Options).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	for _, str := range strings.Split(val, ";") {
		keyValue := strings.Split(str, "=")
		if len(keyValue) == 2 {
			values[keyValue[0]] = keyValue[1]
		}
	}
	b, err := json.Marshal(values)
	if err != nil {
		return nil
	}
	c.values[p.Name] = string(b)
	return nil
}

func (c *Connector) askCondition(p *Property) error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetOptions(p.Options).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	list, ok := p.Conditional[val]
	if ok {
		if err := c.renderList(list); err != nil {
			return nil
		}
	}
	return nil
}
func (c *Connector) askNull(p *Property) error {
	c.values[p.Name] = p.Default
	return nil
}
func (c *Connector) renderList(list []*Property) error {
	for _, p := range list {
		switch p.Kind {
		case "string":
			if err := c.askString(p); err != nil {
				return err
			}
		case "int":
			if err := c.askInt(p); err != nil {
				return err
			}
		case "bool":
			if err := c.askBool(p); err != nil {
				return err
			}
		case "null":
			if err := c.askNull(p); err != nil {
				return err
			}
		case "multilines":
			if err := c.askMultilines(p); err != nil {
				return err
			}
		case "map":
			if err := c.askMap(p); err != nil {
				return err
			}
		case "condition":
			if err := c.askCondition(p); err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *Connector) Render(options DefaultOptions) (map[string]string, error) {
	c.values = map[string]string{}
	c.loadedOptions = options
	if err := c.renderList(c.Properties); err != nil {
		return nil, err
	}
	return c.values, nil
}
func (c *Connector) ColoredYaml() string {
	var propertiesMap []map[string]string
	for _, property := range c.Properties {
		m := property.Map()
		if m != nil {
			propertiesMap = append(propertiesMap, m)
		}
	}
	c.PropertiesSpec = utils.MapArrayToYaml(propertiesMap)
	t := utils.NewTemplate(connectorTemplate, c)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring connector spec,%s", err.Error())
	}
	return string(b)
}

type Connectors []*Connector

func (c Connectors) Validate() error {
	list := map[string]*Connector{}
	for _, connector := range c {
		_, ok := list[connector.Kind]
		if ok {
			return fmt.Errorf("duplicate connector kind: %s", connector.Kind)
		} else {
			list[connector.Kind] = connector
		}
		if err := connector.Validate(); err != nil {
			return err
		}
	}
	return nil
}

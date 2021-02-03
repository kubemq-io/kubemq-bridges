package common

import (
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"strings"
)

type Connector struct {
	Kind             string      `json:"kind"`
	Description      string      `json:"description"`
	Properties       []*Property `json:"properties"`
	Metadata         []*Metadata `json:"metadata"`
	PropertiesSpec   string
	loadedOptions    DefaultOptions
	propertiesValues map[string]string
	metadataValues   map[string]string
	defaultKeys      map[string]string
}

func NewConnector() *Connector {
	return &Connector{
		Kind:             "",
		Description:      "",
		Properties:       nil,
		Metadata:         nil,
		PropertiesSpec:   "",
		loadedOptions:    nil,
		propertiesValues: nil,
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
func (c *Connector) AddMetadata(value *Metadata) *Connector {
	c.Metadata = append(c.Metadata, value)
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
		return nil
	}
	for _, property := range c.Properties {
		if err := property.Validate(); err != nil {
			return err
		}
	}
	return nil
}
func (c *Connector) checkDefaultKey(p *Property) string {

	if c.defaultKeys != nil {
		if p.DefaultFromKey != "" {
			val, ok := c.defaultKeys[p.DefaultFromKey]
			if ok {
				return val
			}
		}
	}
	return p.Default
}
func (c *Connector) askString(p *Property, targetValues map[string]string) error {
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
		SetDefault(c.checkDefaultKey(p)).
		SetOptions(options).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	if val != "" {
		targetValues[p.Name] = val
	}
	return nil
}
func (c *Connector) askInt(p *Property, targetValues map[string]string) error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(c.checkDefaultKey(p)).
		SetHelp(p.Description).
		SetRequired(p.Must).
		SetRange(p.Min, p.Max).
		Render(&val)
	if err != nil {
		return err
	}
	targetValues[p.Name] = fmt.Sprintf("%d", val)
	return nil
}
func (c *Connector) askBool(p *Property, targetValues map[string]string) error {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(c.checkDefaultKey(p)).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	targetValues[p.Name] = fmt.Sprintf("%t", val)
	return nil
}
func (c *Connector) askMultilines(p *Property, targetValues map[string]string) error {
	val := ""
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(c.checkDefaultKey(p)).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	targetValues[p.Name] = val
	return nil
}
func (c *Connector) askMap(p *Property, targetValues map[string]string) error {
	values := map[string]string{}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(c.checkDefaultKey(p)).
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
	targetValues[p.Name] = string(b)
	return nil
}

func (c *Connector) askCondition(p *Property, targetValues map[string]string) error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(c.checkDefaultKey(p)).
		SetOptions(p.Options).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	list, ok := p.Conditional[val]
	if ok {
		if err := c.renderQuestionList(list, targetValues); err != nil {
			return nil
		}
	}
	return nil
}
func (c *Connector) askNull(p *Property, targetValues map[string]string) error {
	targetValues[p.Name] = p.Default
	return nil
}
func (c *Connector) renderQuestionList(list []*Property, targetValues map[string]string) error {
	for _, p := range list {
		switch p.Kind {
		case "string":
			if err := c.askString(p, targetValues); err != nil {
				return err
			}
		case "int":
			if err := c.askInt(p, targetValues); err != nil {
				return err
			}
		case "bool":
			if err := c.askBool(p, targetValues); err != nil {
				return err
			}
		case "null":
			if err := c.askNull(p, targetValues); err != nil {
				return err
			}
		case "multilines":
			if err := c.askMultilines(p, targetValues); err != nil {
				return err
			}
		case "map":
			if err := c.askMap(p, targetValues); err != nil {
				return err
			}
		case "condition":
			if err := c.askCondition(p, targetValues); err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *Connector) renderList(list []*Property, targetValue map[string]string) error {
	var requiredList []*Property
	var notRequiredList []*Property
	for _, p := range list {
		if p.Must {
			requiredList = append(requiredList, p)
		} else {
			notRequiredList = append(notRequiredList, p)
		}
	}
	if err := c.renderQuestionList(requiredList, targetValue); err != nil {
		return err
	}
	if len(notRequiredList) > 0 {
		options := []string{
			"Set them to defaults values",
			"Let me configure them",
		}
		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("check not required").
			SetMessage(fmt.Sprintf("There are %d values which are not mandatory to configure:", len(notRequiredList))).
			SetDefault(options[0]).
			SetOptions(options).
			SetRequired(true).
			Render(&val)
		if err != nil {
			return err
		}
		if val == options[1] {
			if err := c.renderQuestionList(notRequiredList, targetValue); err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *Connector) RenderProperties(options DefaultOptions, defaultKeys map[string]string) (map[string]string, error) {
	c.propertiesValues = map[string]string{}
	c.loadedOptions = options
	c.defaultKeys = defaultKeys
	if err := c.renderList(c.Properties, c.propertiesValues); err != nil {
		return nil, err
	}
	return c.propertiesValues, nil
}
func (c *Connector) RenderMetadata(options DefaultOptions, defaultKeys map[string]string) (map[string]string, error) {
	c.metadataValues = map[string]string{}
	c.loadedOptions = options
	c.defaultKeys = defaultKeys
	if err := c.renderList(c.Properties, c.metadataValues); err != nil {
		return nil, err
	}
	return c.metadataValues, nil
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
			return fmt.Errorf("copied connector kind: %s", connector.Kind)
		} else {
			list[connector.Kind] = connector
		}
		if err := connector.Validate(); err != nil {
			return err
		}
	}
	return nil
}

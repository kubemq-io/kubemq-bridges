package source

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Connection struct {
	addressOptions []string
	properties     map[string]string
	name           string
	kind           string
}

func NewConnection() *Connection {
	return &Connection{
		addressOptions: nil,
		properties:     map[string]string{},
	}
}
func (c *Connection) SetAddress(value []string) *Connection {
	c.addressOptions = value
	return c
}

func (c *Connection) askAddress() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("address").
		SetMessage("Set Kubemq connection address").
		SetDefault("").
		SetHelp("Set address of Kubemq cluster grpc endpoint").
		SetRequired(true).
		SetOptions(c.addressOptions).
		SetValidator(survey.ValidateNoneSpace).
		Render(&val)
	if err != nil {
		return err
	}
	c.properties["address"] = val
	return nil
}
func (c *Connection) askChannel() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("channel").
		SetMessage("Set source channel").
		SetDefault(fmt.Sprintf("source.%s.%s", c.name, c.kind)).
		SetHelp("Set source channel").
		SetRequired(true).
		SetValidator(survey.ValidateNoneSpace).
		Render(&val)
	if err != nil {
		return err
	}
	c.properties["channel"] = val
	return nil
}

func (c *Connection) askGroup() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("group").
		SetMessage("Set source channel group").
		SetDefault("").
		SetHelp("Set source channel group (load balancing)").
		SetRequired(false).
		SetValidator(survey.ValidateNoneSpace).
		Render(&val)
	if err != nil {
		return err
	}
	if val != "" {
		c.properties["group"] = val
	}
	return nil
}

func (c *Connection) askClientID() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("client_id").
		SetMessage("Set source connection client id").
		SetDefault("").
		SetHelp("Set source connection client id").
		SetRequired(false).
		SetValidator(survey.ValidateNoneSpace).
		Render(&val)
	if err != nil {
		return err
	}
	if val != "" {
		c.properties["client_id"] = val
	}
	return nil
}
func (c *Connection) askAuthToken() error {
	val := ""
	err := survey.NewMultiline().
		SetKind("multilines").
		SetName("auth_token").
		SetMessage("Set source connection authentication token").
		SetDefault("").
		SetHelp("Set JWT source connection authentication token").
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val != "" {
		c.properties["auth_token"] = val
	}
	return nil
}
func (c *Connection) askBatchSize() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("batch_size").
		SetMessage("Set source channel queue polling size").
		SetDefault("1").
		SetHelp("Set source channel queue polling size").
		SetRange(1, 1024).
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val > 1 {
		c.properties["batch_size"] = fmt.Sprintf("%d", val)
	}
	return nil
}
func (c *Connection) askWaitTimeout() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("wait_timeout").
		SetMessage("Set source channel queue polling interval in seconds").
		SetDefault("60").
		SetHelp("Set source channel queue polling interval in seconds").
		SetRange(1, 24*60*60).
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val != 60 {
		c.properties["wait_timeout"] = fmt.Sprintf("%d", val)
	}
	return nil
}
func (c *Connection) renderEventsKind() (map[string]string, error) {
	if err := c.askAddress(); err != nil {
		return nil, err
	}
	if err := c.askChannel(); err != nil {
		return nil, err
	}
	if err := c.askGroup(); err != nil {
		return nil, err
	}

	if err := c.askClientID(); err != nil {
		return nil, err
	}
	if err := c.askAuthToken(); err != nil {
		return nil, err
	}
	return c.properties, nil
}
func (c *Connection) renderRPCKind() (map[string]string, error) {
	if err := c.askAddress(); err != nil {
		return nil, err
	}
	if err := c.askChannel(); err != nil {
		return nil, err
	}
	if err := c.askGroup(); err != nil {
		return nil, err
	}
	if err := c.askClientID(); err != nil {
		return nil, err
	}
	if err := c.askAuthToken(); err != nil {
		return nil, err
	}
	return c.properties, nil
}
func (c *Connection) renderQueueKind() (map[string]string, error) {
	if err := c.askAddress(); err != nil {
		return nil, err
	}
	if err := c.askChannel(); err != nil {
		return nil, err
	}
	if err := c.askClientID(); err != nil {
		return nil, err
	}
	if err := c.askAuthToken(); err != nil {
		return nil, err
	}

	if err := c.askBatchSize(); err != nil {
		return nil, err
	}

	if err := c.askWaitTimeout(); err != nil {
		return nil, err
	}

	return c.properties, nil
}

func (c *Connection) Render(name, kind string) (map[string]string, error) {
	c.name = name

	switch kind {
	case "source.queue":
		c.kind = "queues"
		return c.renderQueueKind()
	case "source.events":
		c.kind = "events"
		return c.renderEventsKind()
	case "source.events-store":
		c.kind = "events-store"
		return c.renderEventsKind()
	case "source.command":
		c.kind = "commands"
		return c.renderRPCKind()
	case "source.query":
		c.kind = "queries"
		return c.renderRPCKind()
	default:
		return nil, fmt.Errorf("invalid kind")
	}
}

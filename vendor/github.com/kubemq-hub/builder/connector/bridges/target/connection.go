package target

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
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
		SetMessage("Sets Kubemq connection address").
		SetDefault("").
		SetHelp("Sets address of Kubemq cluster grpc endpoint").
		SetRequired(true).
		SetOptions(c.addressOptions).
		Render(&val)
	if err != nil {
		return err
	}
	c.properties["address"] = val
	return nil
}
func (c *Connection) askDefaultChannel() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("default_channel").
		SetMessage("Sets target default channel").
		SetDefault(fmt.Sprintf("target.%s.%s", c.name, c.kind)).
		SetHelp("Sets target channel").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return err
	}
	c.properties["default_channel"] = val
	return nil
}
func (c *Connection) askChannelList() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("channels").
		SetMessage("Sets target channels list separated by comma").
		SetDefault(fmt.Sprintf("target.%s.%s", c.name, c.kind)).
		SetHelp("Sets target channels list ").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return err
	}
	c.properties["channels"] = val
	return nil
}

func (c *Connection) askExpirationSeconds() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("expiration_seconds").
		SetMessage("Sets target queue message expiration seconds (0 - not expired)").
		SetDefault("0").
		SetHelp("Sets target queue message expiration seconds (0 - not expired)").
		SetRange(0, math.MaxInt32).
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val > 0 {
		c.properties["expiration_seconds"] = fmt.Sprintf("%d", val)
	}

	return nil
}

func (c *Connection) askDelaySeconds() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("delay_seconds").
		SetMessage("Sets target queue message delay seconds (0 - no delay)").
		SetDefault("0").
		SetHelp("Sets target queue message delay seconds (0 - no delay)").
		SetRange(0, math.MaxInt32).
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val > 0 {
		c.properties["delay_seconds"] = fmt.Sprintf("%d", val)
	}
	return nil
}

func (c *Connection) askDeadLetterQueue() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("dead_letter_queue").
		SetMessage("Sets target dead letter queue routing channel").
		SetDefault("").
		SetHelp("Sets target queue dead letter routing max receive fails (0 - no routing").
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val != "" {
		c.properties["dead_letter_queue"] = val
		maxVal := 0
		err := survey.NewInt().
			SetKind("int").
			SetName("max_receive_count").
			SetMessage("Sets target queue dead letter routing max receive fails").
			SetDefault("1").
			SetHelp("Sets target queue dead letter routing max receive fails").
			SetRange(1, math.MaxInt32).
			SetRequired(true).
			Render(&maxVal)
		if err != nil {
			return err
		}
		c.properties["max_receive_count"] = fmt.Sprintf("%d", maxVal)
	}
	return nil
}

func (c *Connection) askTimeoutSeconds() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("timeout_seconds").
		SetMessage("Sets target response timeout seconds").
		SetDefault("60").
		SetHelp("Sets target  response timeout seconds").
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val > 0 {
		c.properties["timeout_seconds"] = fmt.Sprintf("%d", val)
	}

	return nil
}
func (c *Connection) askClientID() error {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("client_id").
		SetMessage("Sets target connection client id").
		SetDefault("").
		SetHelp("Sets target connection client id").
		SetRequired(false).
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
		SetMessage("Sets source connection authentication token").
		SetDefault("").
		SetHelp("Sets JWT source connection authentication token").
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

func (c *Connection) renderQueueKind() (map[string]string, error) {
	if err := c.askAddress(); err != nil {
		return nil, err
	}
	if err := c.askChannelList(); err != nil {
		return nil, err
	}
	if err := c.askClientID(); err != nil {
		return nil, err
	}
	if err := c.askAuthToken(); err != nil {
		return nil, err
	}
	if err := c.askExpirationSeconds(); err != nil {
		return nil, err
	}

	if err := c.askDelaySeconds(); err != nil {
		return nil, err
	}

	if err := c.askDeadLetterQueue(); err != nil {
		return nil, err
	}

	return c.properties, nil
}
func (c *Connection) renderRPCKinds() (map[string]string, error) {
	if err := c.askAddress(); err != nil {
		return nil, err
	}
	if err := c.askDefaultChannel(); err != nil {
		return nil, err
	}
	if err := c.askClientID(); err != nil {
		return nil, err
	}
	if err := c.askAuthToken(); err != nil {
		return nil, err
	}
	if err := c.askTimeoutSeconds(); err != nil {
		return nil, err
	}
	return c.properties, nil
}
func (c *Connection) renderEventsKind() (map[string]string, error) {
	if err := c.askAddress(); err != nil {
		return nil, err
	}
	if err := c.askChannelList(); err != nil {
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

func (c *Connection) Render(name, kind string) (map[string]string, error) {
	c.name = name
	switch kind {
	case "target.queue":
		c.kind = "queues"
		return c.renderQueueKind()
	case "target.events":
		c.kind = "events"
		return c.renderEventsKind()
	case "target.events-store":
		c.kind = "events-store"
		return c.renderEventsKind()
	case "target.command":
		c.kind = "commands"
		return c.renderRPCKinds()
	case "target.query":
		c.kind = "queries"
		return c.renderRPCKinds()
	default:
		return nil, fmt.Errorf("invalid kind")
	}

}

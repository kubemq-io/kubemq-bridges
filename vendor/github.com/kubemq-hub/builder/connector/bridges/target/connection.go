package target

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
)

type Connection struct {
	addressOptions []string
	properties     map[string]string
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
		SetValidator(survey.ValidateNoneSpace).
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
		SetMessage("Set target default channel").
		SetDefault(fmt.Sprintf("target.%s", c.kind)).
		SetHelp("Set target channel").
		SetValidator(survey.ValidateNoneSpace).
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
		SetMessage("Set target channels list separated by comma").
		SetDefault(fmt.Sprintf("target.%s", c.kind)).
		SetHelp("Set target channels list ").
		SetRequired(true).
		SetValidator(survey.ValidateNoneSpace).
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
		SetMessage("Set target queue message expiration seconds (0 - not expired)").
		SetDefault("0").
		SetHelp("Set target queue message expiration seconds (0 - not expired)").
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
		SetMessage("Set target queue message delay seconds (0 - no delay)").
		SetDefault("0").
		SetHelp("Set target queue message delay seconds (0 - no delay)").
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
		SetMessage("Set target dead letter queue routing channel").
		SetDefault("").
		SetHelp("Set target queue dead letter routing max receive fails (0 - no routing").
		SetRequired(false).
		SetValidator(survey.ValidateNoneSpace).
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
			SetMessage("Set target queue dead letter routing max receive fails").
			SetDefault("1").
			SetHelp("Set target queue dead letter routing max receive fails").
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
		SetMessage("Set target response timeout seconds").
		SetDefault("60").
		SetHelp("Set target  response timeout seconds").
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
		SetMessage("Set target connection client id").
		SetDefault("").
		SetHelp("Set target connection client id").
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

func (c *Connection) Render(kind string) (map[string]string, error) {

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

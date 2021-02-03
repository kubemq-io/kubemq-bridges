package source

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Connection struct {
	addressOptions []string
	properties     map[string]string
	kind           string
	bindingName    string
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
		SetHelp("Set source channel").
		SetRequired(true).
		SetValidator(survey.ValidateNoneSpace).
		SetDefault(fmt.Sprintf("%s.%s", c.kind, c.bindingName)).
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
func (c *Connection) askSources() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("sources").
		SetMessage("Set how many sources to subscribe").
		SetDefault("1").
		SetHelp("Set how many sources to subscribe").
		SetRange(1, 1024).
		SetRequired(false).
		Render(&val)
	if err != nil {
		return err
	}
	if val != 1 {
		c.properties["sources"] = fmt.Sprintf("%d", val)
	}
	return nil
}
func (c *Connection) askMaxRequeue() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("max_requeue").
		SetMessage("Set how many times to requeue a request due to target error").
		SetDefault("0").
		SetHelp("Set how many times to requeue a request due to target error").
		SetRange(0, 1024).
		SetRequired(true).
		Render(&val)
	if err != nil {
		return err
	}
	if val != 0 {
		c.properties["max_requeue"] = fmt.Sprintf("%d", val)
	}
	return nil
}
func (c *Connection) askVisibilityTimeout() error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName("visibility_timeout_seconds").
		SetMessage("Set how long to keep current queue message for target processing").
		SetDefault("60").
		SetHelp("Set how long to keep current queue message for target processing (visibility)").
		SetRange(1, 24*60*60).
		SetRequired(true).
		Render(&val)
	if err != nil {
		return err
	}
	if val != 60 {
		c.properties["visibility_timeout_seconds"] = fmt.Sprintf("%d", val)
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
	options := []string{
		"Set them to defaults values",
		"Let me configure them",
	}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("check not required").
		SetMessage("There are 4 values which are not mandatory to configure:").
		SetDefault(options[0]).
		SetOptions(options).
		SetRequired(true).
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val == options[0] {
		return c.properties, nil
	}
	if err := c.askGroup(); err != nil {
		return nil, err
	}
	if err := c.askSources(); err != nil {
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
	options := []string{
		"Set them to defaults values",
		"Let me configure them",
	}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("check not required").
		SetMessage("There are 4 values which are not mandatory to configure:").
		SetDefault(options[0]).
		SetOptions(options).
		SetRequired(true).
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val == options[0] {
		return c.properties, nil
	}
	if err := c.askGroup(); err != nil {
		return nil, err
	}
	if err := c.askSources(); err != nil {
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
	options := []string{
		"Set them to defaults values",
		"Let me configure them",
	}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("check not required").
		SetMessage("There are 6 values which are not mandatory to configure:").
		SetDefault(options[0]).
		SetOptions(options).
		SetRequired(true).
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val == options[0] {
		return c.properties, nil
	}
	if err := c.askSources(); err != nil {
		return nil, err
	}
	if err := c.askMaxRequeue(); err != nil {
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
func (c *Connection) renderQueueStreamKind() (map[string]string, error) {
	if err := c.askAddress(); err != nil {
		return nil, err
	}
	if err := c.askChannel(); err != nil {
		return nil, err
	}
	options := []string{
		"Set them to defaults values",
		"Let me configure them",
	}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("check not required").
		SetMessage("There are 5 values which are not mandatory to configure:").
		SetDefault(options[0]).
		SetOptions(options).
		SetRequired(true).
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val == options[0] {
		return c.properties, nil
	}
	if err := c.askSources(); err != nil {
		return nil, err
	}

	if err := c.askClientID(); err != nil {
		return nil, err
	}
	if err := c.askAuthToken(); err != nil {
		return nil, err
	}

	if err := c.askVisibilityTimeout(); err != nil {
		return nil, err
	}

	if err := c.askWaitTimeout(); err != nil {
		return nil, err
	}

	return c.properties, nil
}
func (c *Connection) Render(kind string, bindingName string) (map[string]string, error) {
	c.bindingName = bindingName
	switch kind {
	case "kubemq.queue":
		c.kind = "queue"
		return c.renderQueueKind()
	case "kubemq.queue-stream":
		c.kind = "queue-stream"
		return c.renderQueueStreamKind()
	case "kubemq.events":
		c.kind = "events"
		return c.renderEventsKind()
	case "kubemq.events-store":
		c.kind = "events-store"
		return c.renderEventsKind()
	case "kubemq.command":
		c.kind = "command"
		return c.renderRPCKind()
	case "kubemq.query":
		c.kind = "query"
		return c.renderRPCKind()
	default:
		return nil, fmt.Errorf("invalid kind")
	}
}

package bridges

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/bridges/source"
	"github.com/kubemq-hub/builder/connector/bridges/target"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Binding struct {
	Name              string            `json:"name"`
	Sources           *source.Source    `json:"sources"`
	Targets           *target.Target    `json:"targets"`
	Properties        map[string]string `json:"properties"`
	SourcesSpec       string            `json:"-" yaml:"-"`
	TargetsSpec       string            `json:"-" yaml:"-"`
	PropertiesSpec    string            `json:"-" yaml:"-"`
	addressOptions    []string
	takenSourceNames  []string
	takenTargetsNames []string
	takenBindingNames []string
	defaultName       string
}

func NewBinding(defaultName string) *Binding {
	return &Binding{
		defaultName: defaultName,
	}
}
func (b *Binding) SetAddress(value []string) *Binding {
	b.addressOptions = value
	return b
}
func (b *Binding) SetTakenSourceNames(value []string) *Binding {
	b.takenSourceNames = value
	return b
}
func (b *Binding) SetTakenTargetsNames(value []string) *Binding {
	b.takenTargetsNames = value
	return b
}
func (b *Binding) SetTakenBindingNames(value []string) *Binding {
	b.takenBindingNames = value
	return b
}
func (b *Binding) SourceName() string {
	if b.Sources != nil {
		return b.Sources.Name
	}
	return ""
}
func (b *Binding) TargetName() string {
	if b.Targets != nil {
		return b.Targets.Name
	}
	return ""
}
func (b *Binding) BindingName() string {
	return b.Name
}
func (b *Binding) confirmSource() bool {
	utils.Println(fmt.Sprintf(promptSourceConfirm, b.Sources.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println(promptSourceReconfigure)
	}
	return val
}
func (b *Binding) confirmTarget() bool {
	utils.Println(fmt.Sprintf(promptTargetConfirm, b.Targets.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println(promptTargetReconfigure)
	}
	return val
}
func (b *Binding) confirmProperties(p *common.Properties) bool {
	utils.Println(fmt.Sprintf(promptPropertiesConfirm, p.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println(promptPropertiesReconfigure)
	}
	return val
}
func (b *Binding) Render() (*Binding, error) {
	var err error
	if b.Name, err = NewName(b.defaultName).
		SetTakenNames(b.takenBindingNames).
		Render(); err != nil {
		return nil, err
	}
	utils.Println(promptSourceStart)
	for {
		if b.Sources, err = source.NewSource(fmt.Sprintf("%s-source", b.defaultName)).
			SetAddress(b.addressOptions).
			SetTakenNames(b.takenSourceNames).
			Render(); err != nil {
			return nil, err
		}
		ok := b.confirmSource()
		if ok {
			b.SourcesSpec = b.Sources.String()
			break
		}
	}
	utils.Println(promptTargetStart)

	for {
		if b.Targets, err = target.NewTarget(fmt.Sprintf("%s-target", b.defaultName)).
			SetAddress(b.addressOptions).
			SetTakenNames(b.takenTargetsNames).
			Render(); err != nil {
			return nil, err
		}
		ok := b.confirmTarget()
		if ok {
			b.TargetsSpec = b.Targets.String()
			break
		}
	}
	utils.Println(promptBindingComplete)
	for {
		p := common.NewProperties()
		if b.Properties, err = p.
			Render(); err != nil {
			return nil, err
		}
		ok := b.confirmProperties(p)
		if ok {
			b.PropertiesSpec = p.String()
			break
		}

	}
	return b, nil
}

func (b *Binding) String() string {
	tpl := utils.NewTemplate(bindingTemplate, b)
	bnd, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring binding spec,%s", err.Error())
	}
	return string(bnd)
}

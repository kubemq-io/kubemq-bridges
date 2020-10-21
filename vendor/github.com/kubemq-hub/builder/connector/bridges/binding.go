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
	isEditMode        bool
	wasEdited         bool
}

func NewBinding(defaultName string) *Binding {
	return &Binding{
		defaultName: defaultName,
	}
}
func (b *Binding) Clone() *Binding {
	newBnd := &Binding{
		Name:              b.Name,
		Sources:           b.Sources.Clone(),
		Targets:           b.Targets.Clone(),
		Properties:        map[string]string{},
		SourcesSpec:       b.SourcesSpec,
		TargetsSpec:       b.TargetsSpec,
		PropertiesSpec:    b.PropertiesSpec,
		addressOptions:    b.addressOptions,
		takenSourceNames:  b.takenSourceNames,
		takenTargetsNames: b.takenTargetsNames,
		takenBindingNames: b.takenBindingNames,
		defaultName:       b.Name,
	}
	for key, val := range b.Properties {
		newBnd.Properties[key] = val
	}

	return newBnd
}
func (b *Binding) SetAddress(value []string) *Binding {
	b.addressOptions = value
	return b
}
func (b *Binding) SetEditMode(value bool) *Binding {
	b.isEditMode = value
	return b
}
func (b *Binding) SetDefaultName(value string) *Binding {
	b.defaultName = value
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
	utils.Println(fmt.Sprintf(promptSourceConfirm, b.Sources.ColoredYaml()))
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
	utils.Println(fmt.Sprintf(promptTargetConfirm, b.Targets.ColoredYaml()))
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
	utils.Println(fmt.Sprintf(promptPropertiesConfirm, p.ColoredYaml()))
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
func (b *Binding) setSource() error {
	if !b.isEditMode {
		utils.Println(promptSourceStart)
		b.Sources = source.NewSource(fmt.Sprintf("%s-source", b.defaultName))
	}

	var err error
	for {
		if b.Sources, err = b.Sources.
			SetAddress(b.addressOptions).
			SetIsEdit(b.isEditMode).
			SetTakenNames(b.takenSourceNames).
			Render(); err != nil {
			return err
		}
		if !b.Sources.WasEdited {
			return nil
		}
		ok := b.confirmSource()
		if ok {
			b.SourcesSpec = b.Sources.ColoredYaml()
			break
		}
	}
	b.wasEdited = b.Sources.WasEdited
	return nil
}
func (b *Binding) setTarget() error {

	if !b.isEditMode {
		utils.Println(promptTargetStart)
		b.Targets = target.NewTarget(fmt.Sprintf("%s-target", b.defaultName))
	}
	var err error
	for {
		if b.Targets, err = b.Targets.
			SetAddress(b.addressOptions).
			SetIsEdit(b.isEditMode).
			SetTakenNames(b.takenSourceNames).
			Render(); err != nil {
			return err
		}
		if !b.Targets.WasEdited {
			return nil
		}
		ok := b.confirmTarget()
		if ok {
			b.TargetsSpec = b.Targets.ColoredYaml()
			break
		}
	}
	b.wasEdited = b.Targets.WasEdited
	return nil
}
func (b *Binding) setProperties() error {
	var err error
	for {
		p := common.NewProperties()
		if b.Properties, err = p.
			Render(); err != nil {
			return err
		}
		if len(b.Properties) == 0 {
			break
		}
		ok := b.confirmProperties(p)
		if ok {
			b.PropertiesSpec = p.ColoredYaml()
			break
		}

	}
	return nil
}
func (b *Binding) showConfiguration() error {
	utils.Println(promptShowBinding, b.Name)
	utils.Println(b.ColoredYaml())
	return nil
}
func (b *Binding) setName() error {
	var err error
	if b.Name, err = NewName(b.defaultName).
		SetTakenNames(b.takenBindingNames).
		Render(); err != nil {
		return err
	}
	b.wasEdited = true
	return nil
}
func (b *Binding) add() (*Binding, error) {
	if err := b.setName(); err != nil {
		return nil, err
	}

	if err := b.setSource(); err != nil {
		return nil, err
	}

	if err := b.setTarget(); err != nil {
		return nil, err
	}
	utils.Println(promptBindingComplete)
	if err := b.setProperties(); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Binding) edit() (*Binding, error) {
	menu := survey.NewMenu("Select Edit Binding operation").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	menu.AddItem("Edit Binding Name", b.setName)
	menu.AddItem("Edit Binding Sources", b.setSource)
	menu.AddItem("Edit Binding Targets", b.setTarget)
	menu.AddItem("Edit Binding Middlewares", b.setProperties)
	menu.AddItem("Show Binding Configuration", b.showConfiguration)
	if err := menu.Render(); err != nil {
		return nil, err
	}
	return b, nil

}
func (b *Binding) Render() (*Binding, error) {
	if b.isEditMode {
		return b.edit()
	}
	return b.add()
}

func (b *Binding) ColoredYaml() string {
	b.SourcesSpec = b.Sources.ColoredYaml()
	b.TargetsSpec = b.Targets.ColoredYaml()
	b.PropertiesSpec = utils.MapToYaml(b.Properties)
	tpl := utils.NewTemplate(bindingTemplate, b)
	bnd, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring binding spec,%s", err.Error())
	}
	return string(bnd)
}

package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Binding struct {
	Name              string            `json:"name" yaml:"name"`
	Source            *Spec             `json:"source" yaml:"source"`
	Target            *Spec             `json:"target" yaml:"target"`
	Properties        map[string]string `json:"properties" yaml:"properties"`
	SourceSpec        string            `json:"-" yaml:"-"`
	TargetSpec        string            `json:"-" yaml:"-"`
	PropertiesSpec    string            `json:"-" yaml:"-"`
	loadedOptions     DefaultOptions
	targetsList       Connectors
	sourcesList       Connectors
	takenBindingNames []string
	defaultName       string
	isEditMode        bool
	wasEdited         bool
}

func NewBinding(defaultName string) *Binding {
	return &Binding{
		Name:              "",
		Source:            NewSpec(),
		Target:            NewSpec(),
		Properties:        map[string]string{},
		SourceSpec:        "",
		TargetSpec:        "",
		PropertiesSpec:    "",
		loadedOptions:     nil,
		targetsList:       nil,
		sourcesList:       nil,
		takenBindingNames: nil,
		defaultName:       defaultName,
		isEditMode:        false,
	}
}
func (b *Binding) SetDefaultOptions(value DefaultOptions) *Binding {
	b.loadedOptions = value
	return b
}
func (b *Binding) Clone() *Binding {
	newBinding := &Binding{
		Name:              b.Name,
		Source:            b.Source.Clone(),
		Target:            b.Target.Clone(),
		Properties:        map[string]string{},
		SourceSpec:        b.SourceSpec,
		TargetSpec:        b.TargetSpec,
		PropertiesSpec:    b.PropertiesSpec,
		loadedOptions:     nil,
		targetsList:       nil,
		sourcesList:       nil,
		takenBindingNames: nil,
		defaultName:       "",
		isEditMode:        false,
	}
	for key, val := range b.Properties {
		newBinding.Properties[key] = val
	}
	return newBinding
}
func (b *Binding) SetTargetsList(value Connectors) *Binding {
	b.targetsList = value
	return b
}
func (b *Binding) SetSourcesList(value Connectors) *Binding {
	b.sourcesList = value
	return b
}
func (b *Binding) SetEditMode(value bool) *Binding {
	b.isEditMode = value
	return b
}
func (b *Binding) SetTakenBindingNames(value []string) *Binding {
	b.takenBindingNames = value
	return b
}
func (b *Binding) SourceName() string {
	return b.Source.Name
}
func (b *Binding) TargetName() string {
	return b.Target.Name
}
func (b *Binding) askKind(kinds []string, currentKind string) (string, error) {
	defaultKind := ""
	if b.isEditMode {
		defaultKind = currentKind
	} else {
		defaultKind = kinds[0]
	}
	if defaultKind == "" {
		defaultKind = kinds[0]
	}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Select Connector Kind").
		SetDefault(defaultKind).
		SetOptions(kinds).
		SetHelp("Select Connector Kind").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (b *Binding) confirmSource() bool {
	utils.Println(fmt.Sprintf(promptSourceConfirm, b.Source.ColoredYaml(sourceSpecTemplate)))
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

	return val
}
func (b *Binding) confirmTarget() bool {
	utils.Println(fmt.Sprintf(promptTargetConfirm, b.Target.ColoredYaml(targetSpecTemplate)))
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

	return val
}
func (b *Binding) confirmProperties(p *Properties) bool {
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
func (b *Binding) addSource(defaultName string) error {
	utils.Println(promptSourceStart)
	var err error
	sourceDefaultName := ""
	if b.isEditMode {
		sourceDefaultName = b.Source.Name
	} else {
		sourceDefaultName = defaultName
	}
	if b.Source.Name, err = NewName(sourceDefaultName).
		RenderSource(); err != nil {
		return err
	}
	var kinds []string
	sources := make(map[string]*Connector)
	for _, c := range b.sourcesList {
		kinds = append(kinds, c.Kind)
		sources[c.Kind] = c
	}
	if len(kinds) == 0 {
		return fmt.Errorf("no source connectors available")
	}

	if b.Source.Kind, err = b.askKind(kinds, b.Source.Kind); err != nil {
		return err
	}
	connector := sources[b.Source.Kind]
	if b.Source.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}

func (b *Binding) editSource() error {
	for {
		ops := []string{
			"Edit Source name",
			"Edit Source kind",
			"Edit Source properties",
			"Show Source configuration",
			"Return",
		}

		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-operation").
			SetMessage("Select Edit Binding Source operation").
			SetDefault(ops[0]).
			SetHelp("Select Edit Binding Source operation").
			SetRequired(true).
			SetOptions(ops).
			Render(&val)
		if err != nil {
			return err
		}
		switch val {
		case ops[0]:
			if b.Source.Name, err = NewName(b.Source.Name).
				RenderSource(); err != nil {
				return err
			}
			b.wasEdited = true
		case ops[1]:
			var kinds []string
			sources := make(map[string]*Connector)
			for _, c := range b.sourcesList {
				kinds = append(kinds, c.Kind)
				sources[c.Kind] = c
			}
			if len(kinds) == 0 {
				return fmt.Errorf("no source connectors available")
			}
			lastKind := b.Source.Kind
			if b.Source.Kind, err = b.askKind(kinds, b.Source.Kind); err != nil {
				return err
			}
			if lastKind != b.Source.Kind {
				connector := sources[b.Source.Kind]
				if b.Source.Properties, err = connector.Render(b.loadedOptions); err != nil {
					return err
				}
			}
			b.wasEdited = true
		case ops[2]:
			var kinds []string
			sources := make(map[string]*Connector)
			for _, c := range b.sourcesList {
				kinds = append(kinds, c.Kind)
				sources[c.Kind] = c
			}
			if len(kinds) == 0 {
				return fmt.Errorf("no source connectors available")
			}
			connector := sources[b.Source.Kind]
			if b.Source.Properties, err = connector.Render(b.loadedOptions); err != nil {
				return err
			}
			b.wasEdited = true
		case ops[3]:
			utils.Println(promptShowSource, b.Source.Name)
			utils.Println(b.Source.ColoredYaml(sourceSpecTemplate))
		default:
			return nil
		}
	}

}

func (b *Binding) addTarget(defaultName string) error {
	utils.Println(promptTargetStart)
	var err error
	targetDefaultName := ""
	if b.isEditMode {
		targetDefaultName = b.Target.Name
	} else {
		targetDefaultName = defaultName
	}
	if b.Target.Name, err = NewName(targetDefaultName).
		RenderTarget(); err != nil {
		return err
	}
	var kinds []string
	targets := make(map[string]*Connector)
	for _, c := range b.targetsList {
		kinds = append(kinds, c.Kind)
		targets[c.Kind] = c
	}
	if len(kinds) == 0 {
		return fmt.Errorf("no targets connectors available")
	}

	if b.Target.Kind, err = b.askKind(kinds, b.Target.Kind); err != nil {
		return err
	}
	connector := targets[b.Target.Kind]
	if b.Target.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}
func (b *Binding) editTarget() error {
	for {
		ops := []string{
			"Edit Target name",
			"Edit Target kind",
			"Edit Target properties",
			"Show Target configuration",
			"Return",
		}

		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-operation").
			SetMessage("Select Edit Binding Target operation").
			SetDefault(ops[0]).
			SetHelp("Select Edit Binding Target operation").
			SetRequired(true).
			SetOptions(ops).
			Render(&val)
		if err != nil {
			return err
		}
		switch val {
		case ops[0]:
			if b.Target.Name, err = NewName(b.Target.Name).
				RenderTarget(); err != nil {
				return err
			}
			b.wasEdited = true
		case ops[1]:
			var kinds []string
			targets := make(map[string]*Connector)
			for _, c := range b.targetsList {
				kinds = append(kinds, c.Kind)
				targets[c.Kind] = c
			}
			if len(kinds) == 0 {
				return fmt.Errorf("no target connectors available")
			}
			lastKind := b.Target.Kind
			if b.Target.Kind, err = b.askKind(kinds, b.Target.Kind); err != nil {
				return err
			}
			if lastKind != b.Target.Kind {
				connector := targets[b.Target.Kind]
				if b.Target.Properties, err = connector.Render(b.loadedOptions); err != nil {
					return err
				}
			}
			b.wasEdited = true
		case ops[2]:
			var kinds []string
			sources := make(map[string]*Connector)
			for _, c := range b.sourcesList {
				kinds = append(kinds, c.Kind)
				sources[c.Kind] = c
			}
			if len(kinds) == 0 {
				return fmt.Errorf("no source connectors available")
			}
			connector := sources[b.Target.Kind]
			if b.Target.Properties, err = connector.Render(b.loadedOptions); err != nil {
				return err
			}
			b.wasEdited = true
		case ops[3]:
			utils.Println(promptShowTarget, b.Target.Name)
			utils.Println(b.Target.ColoredYaml(sourceSpecTemplate))
		default:
			return nil
		}
	}

}

func (b *Binding) setName() error {
	var err error
	if b.Name, err = NewName(b.defaultName).
		SetTakenNames(b.takenBindingNames).
		RenderBinding(); err != nil {
		return err
	}
	return nil
}
func (b *Binding) showConfiguration() error {
	utils.Println(promptShowBinding, b.Name)
	utils.Println(b.ColoredYaml())

	return nil
}
func (b *Binding) setProperties() error {
	var err error
	for {
		p := NewProperties()
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
func (b *Binding) edit() (*Binding, error) {
	for {
		ops := []string{
			"Edit binding name",
			"Edit binding Source",
			"Edit binding Target",
			"Edit binding Middlewares",
			"Show binding configuration",
			"Done",
		}

		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-operation").
			SetMessage("Select Edit Binding operation").
			SetDefault(ops[0]).
			SetHelp("Select Edit Binding operation").
			SetRequired(true).
			SetOptions(ops).
			Render(&val)
		if err != nil {
			return nil, err
		}
		switch val {
		case ops[0]:
			if err := b.setName(); err != nil {
				return nil, err
			}
			b.wasEdited = true
		case ops[1]:
			for {
				if err := b.editSource(); err != nil {
					return nil, err
				}
				if b.confirmSource() {
					break
				}
			}

		case ops[2]:
			if err := b.editTarget(); err != nil {
				return nil, err
			}
			if b.confirmTarget() {
				break
			}

		case ops[3]:
			if err := b.setProperties(); err != nil {
				return nil, err
			}
			b.wasEdited = true
		case ops[4]:
			if err := b.showConfiguration(); err != nil {
				return nil, err
			}
		default:
			return b, nil
		}
	}

}
func (b *Binding) add() (*Binding, error) {
	if err := b.setName(); err != nil {
		return nil, err
	}
	for {
		if err := b.addSource(fmt.Sprintf("%s-source", b.Name)); err != nil {
			return nil, err
		}
		if b.confirmSource() {
			break
		}
	}

	for {
		if err := b.addTarget(fmt.Sprintf("%s-target", b.Name)); err != nil {
			return nil, err
		}
		if b.confirmTarget() {
			break
		}
	}
	utils.Println(promptBindingComplete)
	var err error
	for {
		p := NewProperties()
		if b.Properties, err = p.
			Render(); err != nil {
			return nil, err
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
	return b, nil
}

func (b *Binding) Render() (*Binding, error) {
	if b.isEditMode {
		return b.edit()
	}
	return b.add()

}

func (b *Binding) ColoredYaml() string {
	tpl := utils.NewTemplate(bindingTemplate, b)
	b.TargetSpec = b.Target.ColoredYaml(targetSpecTemplate)
	b.SourceSpec = b.Source.ColoredYaml(sourceSpecTemplate)
	b.PropertiesSpec = utils.MapToYaml(b.Properties)
	bnd, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring binding spec,%s", err.Error())
	}
	return string(bnd)
}
func (b *Binding) TableRowShort() []interface{} {
	var list []interface{}
	ms := utils.MapFlatten(b.Properties)
	if ms == "" {
		ms = "none"
	}
	list = append(list, b.Name, b.Source.TableItemShort(), b.Target.TableItemShort(), ms)
	return list
}

package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Binding struct {
	Name           string            `json:"name" yaml:"name"`
	Source         *Spec             `json:"source" yaml:"source"`
	Target         *Spec             `json:"target" yaml:"target"`
	Properties     map[string]string `json:"properties" yaml:"properties"`
	SourceSpec     string            `json:"-" yaml:"-"`
	TargetSpec     string            `json:"-" yaml:"-"`
	PropertiesSpec string            `json:"-" yaml:"-"`
	Side           string            `json:"-" yaml:"-"`
	loadedOptions  DefaultOptions
	targetsList    Connectors
	sourcesList    Connectors
	defaultName    string
	isEditMode     bool
}

func NewBinding(defaultName string, side string, loadedOptions DefaultOptions, targetList, sourcesList Connectors) *Binding {
	return &Binding{
		Name:           "",
		Source:         NewSpec(),
		Target:         NewSpec(),
		Properties:     map[string]string{},
		SourceSpec:     "",
		TargetSpec:     "",
		PropertiesSpec: "",
		loadedOptions:  loadedOptions,
		targetsList:    targetList,
		sourcesList:    sourcesList,
		defaultName:    defaultName,
		isEditMode:     false,
		Side:           side,
	}
}

func (b *Binding) Clone() *Binding {
	newBinding := &Binding{
		Name:           b.Name,
		Source:         b.Source.Clone(),
		Target:         b.Target.Clone(),
		Properties:     map[string]string{},
		SourceSpec:     b.SourceSpec,
		TargetSpec:     b.TargetSpec,
		PropertiesSpec: b.PropertiesSpec,
		loadedOptions:  b.loadedOptions,
		targetsList:    b.targetsList,
		sourcesList:    b.sourcesList,
		defaultName:    b.defaultName,
		isEditMode:     false,
		Side:           b.Side,
	}
	for key, val := range b.Properties {
		newBinding.Properties[key] = val
	}
	return newBinding
}
func (b *Binding) Validate() error {
	return nil
}

func (b *Binding) SetEditMode(value bool) *Binding {
	b.isEditMode = value
	return b
}

func (b *Binding) askKind(connector string, kinds []string, currentKind string) (string, error) {
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
		SetMessage(fmt.Sprintf("Select %s Kind", connector)).
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

func (b *Binding) addSource(defaultName string) error {
	var err error
	var kinds []string
	sources := make(map[string]*Connector)
	for _, c := range b.sourcesList {
		kinds = append(kinds, c.Kind)
		sources[c.Kind] = c
	}

	if len(kinds) == 0 {
		return fmt.Errorf("no source connectors available")
	}

	if b.Source.Kind, err = b.askKind("Source", kinds, b.Source.Kind); err != nil {
		return err
	}
	connector := sources[b.Source.Kind]
	if b.Source.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}

func (b *Binding) editSource() (*Spec, error) {
	var result *Spec
	edited := b.Clone()
	form := survey.NewForm("Select Edit Source Option:")

	ftKind := new(string)
	*ftKind = fmt.Sprintf("<k> Edit Source Kind (%s)", edited.Source.Kind)
	ftProperties := new(string)
	*ftProperties = fmt.Sprintf("<p> Edit Source Properties (%s)", edited.Source.Kind)

	form.AddItem(ftKind, func() error {
		var kinds []string
		sources := make(map[string]*Connector)
		for _, c := range edited.sourcesList {
			kinds = append(kinds, c.Kind)
			sources[c.Kind] = c
		}
		kinds = append(kinds, "<back>")
		lastKind := edited.Source.Kind
		selected := ""

		var err error
		if selected, err = edited.askKind("Source", kinds, edited.Source.Kind); err != nil {
			return err
		}
		if selected == "<back>" {
			edited.Source.Kind = lastKind
			return nil
		} else {
			edited.Source.Kind = selected
		}
		if lastKind != edited.Source.Kind {
			connector := sources[edited.Source.Kind]
			if edited.Source.Properties, err = connector.Render(edited.loadedOptions); err != nil {
				return err
			}
		}
		*ftKind = fmt.Sprintf("<k> Edit Source Kind (%s)", edited.Source.Kind)
		*ftProperties = fmt.Sprintf("<p> Edit Source Properties (%s)", edited.Source.Kind)
		return nil
	})

	form.AddItem(ftProperties, func() error {
		var kinds []string
		sources := make(map[string]*Connector)
		for _, c := range edited.sourcesList {
			kinds = append(kinds, c.Kind)
			sources[c.Kind] = c
		}
		if len(kinds) == 0 {
			return fmt.Errorf("no source connectors available")
		}
		var err error
		connector := sources[edited.Source.Kind]
		if edited.Source.Properties, err = connector.Render(edited.loadedOptions); err != nil {
			return err
		}
		*ftProperties = fmt.Sprintf("<p> Edit Source Properties (%s)", edited.Source.Kind)
		return nil
	})

	form.AddItem("<s> Show Source Configuration", func() error {
		utils.Println(promptShowSource)
		utils.Println("%s\n", edited.Source.ColoredYaml(sourceSpecTemplate))
		return nil
	})
	form.SetOnSaveFn(func() error {
		if err := edited.Validate(); err != nil {
			return err
		}
		result = edited.Source
		return nil
	})

	form.SetOnCancelFn(func() error {
		result = b.Source
		return nil
	})
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Binding) addTarget(defaultName string) error {
	var err error
	var kinds []string
	targets := make(map[string]*Connector)
	for _, c := range b.targetsList {
		kinds = append(kinds, c.Kind)
		targets[c.Kind] = c
	}
	if len(kinds) == 0 {
		return fmt.Errorf("no targets connectors available")
	}

	if b.Target.Kind, err = b.askKind("Target", kinds, b.Target.Kind); err != nil {
		return err
	}
	connector := targets[b.Target.Kind]
	if b.Target.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}
func (b *Binding) editTarget() (*Spec, error) {
	var result *Spec
	edited := b.Clone()
	form := survey.NewForm("Select Edit Target Option:")

	ftKind := new(string)
	*ftKind = fmt.Sprintf("<k> Edit Target Kind (%s)", edited.Target.Kind)
	ftProperties := new(string)
	*ftProperties = fmt.Sprintf("<p> Edit Target Properties (%s)", edited.Target.Kind)

	form.AddItem(ftKind, func() error {
		var kinds []string
		targets := make(map[string]*Connector)
		for _, c := range edited.targetsList {
			kinds = append(kinds, c.Kind)
			targets[c.Kind] = c
		}
		kinds = append(kinds, "<back>")
		lastKind := edited.Target.Kind
		selected := ""
		var err error
		if selected, err = edited.askKind("Target", kinds, edited.Target.Kind); err != nil {
			return err
		}
		if selected == "<back>" {
			edited.Target.Kind = lastKind
			return nil

		} else {
			edited.Target.Kind = selected
		}

		if lastKind != edited.Target.Kind {
			connector := targets[edited.Target.Kind]
			if edited.Target.Properties, err = connector.Render(edited.loadedOptions); err != nil {
				return err
			}
		}
		*ftKind = fmt.Sprintf("<k> Edit Target Kind (%s)", edited.Target.Kind)
		*ftProperties = fmt.Sprintf("<p> Edit Target Properties (%s)", edited.Target.Kind)
		return nil
	})

	form.AddItem(ftProperties, func() error {
		var kinds []string
		targets := make(map[string]*Connector)
		for _, c := range edited.targetsList {
			kinds = append(kinds, c.Kind)
			targets[c.Kind] = c
		}
		if len(kinds) == 0 {
			return fmt.Errorf("no target connectors available")
		}
		var err error
		connector := targets[edited.Target.Kind]
		if edited.Target.Properties, err = connector.Render(edited.loadedOptions); err != nil {
			return err
		}
		*ftProperties = fmt.Sprintf("<p> Edit Target Properties (%s)", edited.Target.Kind)
		return nil
	})

	form.AddItem("<s> Show Target Configuration", func() error {
		utils.Println(promptShowTarget)
		utils.Println("%s\n", edited.Target.ColoredYaml(targetSpecTemplate))
		return nil
	})
	form.SetOnSaveFn(func() error {
		if err := edited.Validate(); err != nil {
			return err
		}
		result = edited.Target
		return nil
	})

	form.SetOnCancelFn(func() error {
		result = b.Target
		return nil
	})
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Binding) setName() error {
	var err error
	if b.Name, err = NewName(b.defaultName).
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
	p := NewProperties(b.Properties)
	if b.Properties, err = p.
		Render(); err != nil {
		return err
	}
	b.PropertiesSpec = p.ColoredYaml()
	return nil
}
func (b *Binding) edit() (*Binding, error) {
	var result *Binding
	edited := b.Clone().
		SetEditMode(true)
	form := survey.NewForm(fmt.Sprintf("Select Edit %s Binding Option:", edited.Name))

	ftName := new(string)
	*ftName = fmt.Sprintf("<n> Edit Binding's Name (%s)", edited.Name)
	form.AddItem(ftName, func() error {
		if err := edited.setName(); err != nil {
			return err
		}
		*ftName = fmt.Sprintf("<n> Edit Binding's Name (%s)", edited.Name)
		return nil
	})

	ftSource := new(string)
	*ftSource = fmt.Sprintf("<s> Edit Binding's Source (%s)", edited.Source.Kind)
	form.AddItem(ftSource, func() error {
		var err error
		if edited.Source, err = edited.editSource(); err != nil {
			return err
		}
		*ftSource = fmt.Sprintf("<s> Edit Binding's Source (%s)", edited.Source.Kind)
		return nil
	})

	ftTarget := new(string)
	*ftTarget = fmt.Sprintf("<t> Edit Binding's Target (%s)", edited.Target.Kind)
	form.AddItem(ftTarget, func() error {
		var err error
		if edited.Target, err = edited.editTarget(); err != nil {
			return err
		}
		*ftTarget = fmt.Sprintf("<t> Edit Binding's Target (%s)", edited.Target.Kind)
		return nil
	})

	form.AddItem("<m> Edit Binding's Middlewares", edited.setProperties)

	form.AddItem("<c> Show Binding Configuration", edited.showConfiguration)

	form.SetOnSaveFn(func() error {
		if err := edited.Validate(); err != nil {
			return err
		}
		result = edited
		return nil
	})

	form.SetOnCancelFn(func() error {
		result = b
		return nil
	})
	if err := form.Render(); err != nil {
		return nil, err
	}

	return result, nil

}
func (b *Binding) add() (*Binding, error) {
	if err := b.setName(); err != nil {
		return nil, err
	}

	if err := b.addSource(fmt.Sprintf("%s-source", b.Name)); err != nil {
		return nil, err
	}

	if err := b.addTarget(fmt.Sprintf("%s-target", b.Name)); err != nil {
		return nil, err
	}

	utils.Println(promptBindingComplete)
	var err error

	p := NewProperties(b.Properties)
	if b.Properties, err = p.
		Render(); err != nil {
		return nil, err
	}
	b.PropertiesSpec = p.ColoredYaml()

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
func (b *Binding) BelongToClusterAddress(address string, side string) bool {
	switch side {
	case "sources":
		return b.Target.IsKubemqAddress(address)
	case "targets":
		return b.Source.IsKubemqAddress(address)
	default:
		return false
	}
}
func generateUniqueIntegrationName(takenNames []string) string {
	for i := len(takenNames) + 1; i < 10000000; i++ {
		name := fmt.Sprintf("integration-%d", i)
		found := false
		for _, taken := range takenNames {
			if taken == name {
				found = true
				break
			}
		}
		if !found {
			return name
		}
	}
	return ""
}
func AddSourceIntegration(takenNames []string, connectorsManifest []byte, loadedOptions DefaultOptions) (*Binding, error) {

	b := &Binding{
		Name:           "",
		Source:         nil,
		Target:         nil,
		Properties:     nil,
		SourceSpec:     "",
		TargetSpec:     "",
		PropertiesSpec: "",
		Side:           "sources",
		loadedOptions:  loadedOptions,
		targetsList:    nil,
		sourcesList:    nil,
		defaultName:    "",
		isEditMode:     false,
	}
	m, err := LoadManifest(connectorsManifest)
	if err != nil {
		return nil, err
	}
	b.targetsList = m.Targets
	b.sourcesList = m.Sources
	// Setting Name
	err = survey.NewString().
		SetKind("string").
		SetName("name").
		SetMessage("Set Unique Integration Name:").
		SetDefault(generateUniqueIntegrationName(takenNames)).
		SetRequired(true).
		SetInvalidOptions(takenNames).
		SetValidator(survey.ValidateNoneSpace).
		SetInvalidOptionsMessage("Integration name must be unique").
		Render(&b.Name)
	if err != nil {
		return nil, err
	}

	// Configuring Source Kinds
	b.Source = &Spec{
		Kind:       "",
		Properties: nil,
	}

	var sourcesKinds []string
	sources := make(map[string]*Connector)
	for _, con := range b.sourcesList {
		sourcesKinds = append(sourcesKinds, con.Kind)
		sources[con.Kind] = con
	}

	if len(sourcesKinds) == 0 {
		return nil, fmt.Errorf("no source connectors available")
	}

	err = survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Select Source Kind:").
		SetDefault(sourcesKinds[0]).
		SetOptions(sourcesKinds).
		SetRequired(true).
		Render(&b.Source.Kind)
	if err != nil {
		return nil, err
	}

	if b.Source.Properties, err = sources[b.Source.Kind].Render(b.loadedOptions); err != nil {
		return nil, err
	}

	// Configuring Targets Kinds
	b.Target = &Spec{
		Kind:       "",
		Properties: nil,
	}

	var targetsKinds []string
	targets := make(map[string]*Connector)
	for _, con := range b.targetsList {
		targetsKinds = append(targetsKinds, con.Kind)
		targets[con.Kind] = con
	}

	if len(targetsKinds) == 0 {
		return nil, fmt.Errorf("no target connectors available")
	}

	err = survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Select Kubemq Target Kind:").
		SetDefault(targetsKinds[0]).
		SetOptions(targetsKinds).
		SetRequired(true).
		Render(&b.Target.Kind)
	if err != nil {
		return nil, err
	}

	if b.Target.Properties, err = targets[b.Target.Kind].Render(b.loadedOptions); err != nil {
		return nil, err
	}

	// Configuring Middlewares
	p := NewProperties(b.Properties)
	if b.Properties, err = p.
		Render(); err != nil {
		return nil, err
	}

	utils.Println("<cyan>Here is the configuration of %s Integration:</>%s", b.Name, b.ColoredYaml())
	val := true
	err = survey.NewBool().
		SetKind("bool").
		SetName("confirmation").
		SetMessage("Would you like to save this Integration").
		SetDefault("true").
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val {
		return b, nil
	}
	return nil, nil
}

func AddTargetIntegration(takenNames []string, connectorsManifest []byte, loadedOptions DefaultOptions) (*Binding, error) {

	b := &Binding{
		Name:           "",
		Source:         nil,
		Target:         nil,
		Properties:     nil,
		SourceSpec:     "",
		TargetSpec:     "",
		PropertiesSpec: "",
		Side:           "targets",
		loadedOptions:  loadedOptions,
		targetsList:    nil,
		sourcesList:    nil,
		defaultName:    "",
		isEditMode:     false,
	}
	m, err := LoadManifest(connectorsManifest)
	if err != nil {
		return nil, err
	}
	b.targetsList = m.Targets
	b.sourcesList = m.Sources
	// Setting Name
	err = survey.NewString().
		SetKind("string").
		SetName("name").
		SetMessage("Set Unique Integration Name:").
		SetDefault(generateUniqueIntegrationName(takenNames)).
		SetRequired(true).
		SetInvalidOptions(takenNames).
		SetValidator(survey.ValidateNoneSpace).
		SetInvalidOptionsMessage("Integration name must be unique").
		Render(&b.Name)
	if err != nil {
		return nil, err
	}

	// Configuring Targets Kinds
	b.Target = &Spec{
		Kind:       "",
		Properties: nil,
	}

	var targetsKinds []string
	targets := make(map[string]*Connector)
	for _, con := range b.targetsList {
		targetsKinds = append(targetsKinds, con.Kind)
		targets[con.Kind] = con
	}

	if len(targetsKinds) == 0 {
		return nil, fmt.Errorf("no target connectors available")
	}

	err = survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Select Target Kind (Scroll down for more):").
		SetDefault(targetsKinds[0]).
		SetOptions(targetsKinds).
		SetRequired(true).
		Render(&b.Target.Kind)
	if err != nil {
		return nil, err
	}

	if b.Target.Properties, err = targets[b.Target.Kind].Render(b.loadedOptions); err != nil {
		return nil, err
	}

	// Configuring Source Kinds
	b.Source = &Spec{
		Kind:       "",
		Properties: nil,
	}

	var sourcesKinds []string
	sources := make(map[string]*Connector)
	for _, con := range b.sourcesList {
		sourcesKinds = append(sourcesKinds, con.Kind)
		sources[con.Kind] = con
	}

	if len(sourcesKinds) == 0 {
		return nil, fmt.Errorf("no source connectors available")
	}

	err = survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Select Kubemq Source Kind:").
		SetDefault(sourcesKinds[0]).
		SetOptions(sourcesKinds).
		SetRequired(true).
		Render(&b.Source.Kind)
	if err != nil {
		return nil, err
	}

	if b.Source.Properties, err = sources[b.Source.Kind].Render(b.loadedOptions); err != nil {
		return nil, err
	}

	// Configuring Middlewares
	p := NewProperties(b.Properties)
	if b.Properties, err = p.
		Render(); err != nil {
		return nil, err
	}

	utils.Println("<cyan>Here is the configuration of %s Integration:</>%s", b.Name, b.ColoredYaml())
	val := true
	err = survey.NewBool().
		SetKind("bool").
		SetName("confirmation").
		SetMessage("Would you like to save this Integration").
		SetDefault("true").
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val {
		return b, nil
	}
	return nil, nil
}

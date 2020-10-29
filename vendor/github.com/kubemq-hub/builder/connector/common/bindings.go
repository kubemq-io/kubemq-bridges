package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"gopkg.in/yaml.v2"
	"sort"
)

type Bindings struct {
	Bindings      []*Binding `json:"bindings" yaml:"bindings"`
	Side          string     `json:"-" yaml:"-"`
	manifest      *Manifest
	loadedOptions DefaultOptions
	defaultName   string
}

func NewBindings(defaultName string, bindings []*Binding, side string, loadedOptions DefaultOptions, manifest *Manifest) *Bindings {
	return &Bindings{
		Bindings:      bindings,
		Side:          side,
		manifest:      manifest,
		loadedOptions: loadedOptions,
		defaultName:   defaultName,
	}
}
func (b *Bindings) Clone() *Bindings {
	cloned := &Bindings{
		Bindings:      nil,
		manifest:      b.manifest,
		loadedOptions: b.loadedOptions,
		defaultName:   b.defaultName,
		Side:          b.Side,
	}
	for _, binding := range b.Bindings {
		cloned.Bindings = append(cloned.Bindings, binding.Clone())
	}
	return cloned
}
func (b *Bindings) Update(manifest *Manifest, loadedOptions DefaultOptions) *Bindings {
	b.manifest = manifest
	b.loadedOptions = loadedOptions
	for _, binding := range b.Bindings {
		binding.loadedOptions = loadedOptions
		binding.targetsList = b.manifest.Targets
		binding.sourcesList = b.manifest.Sources
	}
	return b
}
func (b *Bindings) Sort() {
	sort.Slice(b.Bindings, func(i, j int) bool {
		return b.Bindings[i].Name < b.Bindings[j].Name
	})
}
func (b *Bindings) addBinding() error {
	bnd := NewBinding(b.GenerateNewBindingName(), b.Side, b.loadedOptions, b.manifest.Targets, b.manifest.Sources)
	var err error
	if bnd, err = bnd.
		Render(); err != nil {
		return err
	}
	for _, binding := range b.Bindings {
		if bnd.Name == binding.Name {
			return fmt.Errorf("added binding name it not unique, binding %s was not added", bnd.Name)
		}
	}
	utils.Println(promptBindingAddConfirmation, bnd.Name)
	b.Bindings = append(b.Bindings, bnd)
	b.Sort()
	return nil
}

func (b *Bindings) SwitchOrRemove(old, new *Binding) {
	var newBindingList []*Binding
	for _, binding := range b.Bindings {
		if old.Name != binding.Name {
			newBindingList = append(newBindingList, binding)
		}
	}
	if new != nil {
		newBindingList = append(newBindingList, new)
	}
	b.Bindings = newBindingList
	b.Sort()
}

func (b *Bindings) editBinding() error {
	menu := survey.NewMenu("Select Binding to edit:").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		SetDisableLoop(true)
	for _, binding := range b.Bindings {
		edited := binding.Clone()
		origin := binding
		editFn := func() error {
			var err error
			if edited, err = edited.
				SetEditMode(true).
				Render(); err != nil {
				return err
			}
			if origin.Name != edited.Name {
				if origin.Name != edited.Name {
					for _, binding := range b.Bindings {
						if edited.Name == binding.Name {
							return fmt.Errorf("binding name %s is not unique, binding %s was not edited", edited.Name, origin.Name)
						}
					}
				}
			}
			b.SwitchOrRemove(origin, edited)
			utils.Println(promptBindingEditedConfirmation, edited.Name)
			return nil
		}
		menu.AddItem(binding.Name, editFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}
func (b *Bindings) deleteBinding() error {
	menu := survey.NewMultiSelectMenu("Select Binding to delete:")
	for _, binding := range b.Bindings {
		deleted := binding
		deleteFn := func() error {
			val := false
			if err := survey.NewBool().
				SetName("confirm-delete").
				SetMessage(fmt.Sprintf("Are you sure you want to delete %s binding", deleted.Name)).
				SetRequired(true).
				SetDefault("false").
				Render(&val); err != nil {
				return err
			}
			if val {
				b.SwitchOrRemove(deleted, nil)
				utils.Println(promptBindingDeleteConfirmation, deleted.Name)
				return nil
			}
			return nil
		}
		menu.AddItem(binding.Name, deleteFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	b.Sort()
	return nil
}

func (b *Bindings) copyBinding() error {
	menu := survey.NewMultiSelectMenu("Select Binding to copy:")
	for _, binding := range b.Bindings {
		cloned := binding.Clone()
		origin := binding
		copyFn := func() error {
			if err := cloned.setName(); err != nil {
				return err
			}
			for _, binding := range b.Bindings {
				if cloned.Name == binding.Name {
					return fmt.Errorf("copied binding name (%s) must be unique\n", cloned.Name)
				}
			}
			checkEdit := false
			if err := survey.NewBool().
				SetKind("bool").
				SetMessage("Would you like to edit the copied binding before saving").
				SetRequired(true).
				SetDefault("false").
				Render(&checkEdit); err != nil {
				return err
			}
			if checkEdit {
				var err error
				cloned, err = cloned.edit()
				if err != nil {
					return err
				}
			}
			if origin.Name != cloned.Name {
				for _, binding := range b.Bindings {
					if cloned.Name == binding.Name {
						return fmt.Errorf("binding name %s is not unique, binding %s was not edited", cloned.Name, origin.Name)
					}
				}
			}

			b.Bindings = append(b.Bindings, cloned)
			return nil
		}
		menu.AddItem(binding.Name, copyFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	b.Sort()
	return nil
}
func (b *Bindings) listBindings() error {

	menu := survey.NewMenu("Select Binding to show configuration:").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, binding := range b.Bindings {
		selected := binding
		showFn := func() error {
			utils.Println(promptShowBinding, selected.Name)
			utils.Println("%s\n", selected.ColoredYaml())
			utils.WaitForEnter()
			return nil
		}
		menu.AddItem(binding.Name, showFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (b *Bindings) Render() ([]byte, error) {
	var result *Bindings
	clone := b.Clone()
	for _, binding := range clone.Bindings {
		binding.loadedOptions = b.loadedOptions
		binding.sourcesList = b.manifest.Sources
		binding.targetsList = b.manifest.Targets
	}
	form := survey.NewForm("Select Manage Bindings Option:")

	form.AddItem("<a> Add Binding", clone.addBinding)
	form.AddItem("<e> Edit Binding", clone.editBinding)
	form.AddItem("<c> Copy Binding", clone.copyBinding)
	form.AddItem("<d> Delete Binding", clone.deleteBinding)
	form.AddItem("<l> List of Bindings", clone.listBindings)

	form.SetOnSaveFn(func() error {
		if err := clone.Validate(); err != nil {
			return err
		}
		result = clone

		return nil
	})
	form.SetOnCancelFn(func() error {
		result = b
		return nil
	})
	form.SetOnErrorFn(survey.FormShowErrorFn)
	if err := form.Render(); err != nil {
		return nil, err
	}
	result.Sort()
	return yaml.Marshal(result)
}

func (b *Bindings) Yaml() ([]byte, error) {
	return yaml.Marshal(b)
}

func Unmarshal(data []byte) (*Bindings, error) {
	bnd := &Bindings{}
	err := yaml.Unmarshal(data, bnd)
	if err != nil {
		return nil, err
	}
	return bnd, nil
}

func (b *Bindings) Validate() error {

	return nil
}
func (b *Bindings) AddIntegration(integration *Binding) error {
	if err := integration.Validate(); err != nil {
		return err
	}
	b.Bindings = append(b.Bindings, integration)
	return b.Validate()
}
func (b *Bindings) RemoveIntegration(integration *Binding) {
	b.SwitchOrRemove(integration, nil)
}

func (b *Bindings) GetBindingsForCluster(address string) []*Binding {
	var list []*Binding

	for _, binding := range b.Bindings {
		binding.Side = b.Side
		if binding.BelongToClusterAddress(address, b.Side) {
			list = append(list, binding)
		}
	}
	return list
}
func (b *Bindings) checkUniqueBindingName(name string) bool {
	for _, binding := range b.Bindings {
		if binding.Name == name {
			return false
		}
	}
	return true
}

func (b *Bindings) GenerateNewBindingName() string {
	for i := len(b.Bindings) + 1; i < 10000000; i++ {
		name := fmt.Sprintf("binding-%d", i)
		if b.checkUniqueBindingName(name) {
			return name
		}
	}
	return ""
}

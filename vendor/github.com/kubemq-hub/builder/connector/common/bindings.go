package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"gopkg.in/yaml.v2"
	"sort"
)

type Bindings struct {
	Bindings          []*Binding `json:"bindings" yaml:"bindings"`
	manifest          *Manifest
	loadedOptions     DefaultOptions
	takenBindingNames []string
	defaultName       string
}

func NewBindings(defaultName string) *Bindings {
	return &Bindings{
		defaultName: defaultName,
	}
}
func (b *Bindings) Clone() *Bindings {
	cloned := &Bindings{
		Bindings:          nil,
		manifest:          b.manifest,
		loadedOptions:     b.loadedOptions,
		takenBindingNames: b.takenBindingNames,
		defaultName:       b.defaultName,
	}
	for _, binding := range b.Bindings {
		cloned.Bindings = append(cloned.Bindings, binding.Clone())
	}
	return cloned
}
func (b *Bindings) SetBindings(value []*Binding) *Bindings {
	b.Bindings = value
	return b
}
func (b *Bindings) SetManifest(value *Manifest) *Bindings {
	b.manifest = value
	return b
}
func (b *Bindings) SetDefaultOptions(value DefaultOptions) *Bindings {
	b.loadedOptions = value
	return b
}
func (b *Bindings) SetDefaultName(value string) *Bindings {
	b.defaultName = value
	return b
}

func (b *Bindings) sort() {
	sort.Slice(b.Bindings, func(i, j int) bool {
		return b.Bindings[i].Name < b.Bindings[j].Name
	})
}
func (b *Bindings) addBinding() error {
	bnd := NewBinding(fmt.Sprintf("binding-%d", len(b.Bindings)+1))
	var err error
	if bnd, err = bnd.
		SetDefaultOptions(b.loadedOptions).
		SetSourcesList(b.manifest.Sources).
		SetTargetsList(b.manifest.Targets).
		SetTakenBindingNames(b.takenBindingNames).
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
	b.sort()
	return nil
}

func (b *Bindings) switchOrRemove(old, new *Binding) {
	var newBindingList []*Binding
	var newTakenBindingNames []string

	for _, binding := range b.Bindings {
		if old.Name != binding.Name {
			newBindingList = append(newBindingList, binding)
			newTakenBindingNames = append(newTakenBindingNames, binding.Name)
		}
	}
	if new != nil {
		newBindingList = append(newBindingList, new)
		newTakenBindingNames = append(newTakenBindingNames, new.Name)
	}
	b.Bindings = newBindingList
	b.takenBindingNames = newTakenBindingNames
	b.sort()
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
				SetDefaultOptions(b.loadedOptions).
				SetSourcesList(b.manifest.Sources).
				SetTargetsList(b.manifest.Targets).
				SetTakenBindingNames(b.takenBindingNames).
				Render(); err != nil {
				return err
			}
			for _, binding := range b.Bindings {
				if edited.Name == binding.Name {
					return fmt.Errorf("binding name %s is not unique, binding %s was not edited", edited.Name, origin.Name)
				}
			}

			b.switchOrRemove(origin, edited)
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
	menu := survey.NewMenu("Select Binding to delete:").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		SetDisableLoop(true)
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
				b.switchOrRemove(deleted, nil)
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
	b.sort()
	return nil
}

func (b *Bindings) copyBinding() error {
	menu := survey.NewMenu("Select Binding to copy:").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		SetDisableLoop(true)
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
			for _, binding := range b.Bindings {
				if cloned.Name == binding.Name {
					return fmt.Errorf("binding name %s is not unique, binding %s was not edited", cloned.Name, origin.Name)
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
	b.sort()
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
	form := survey.NewForm("Select Manage Bindings Option:")

	form.AddItem("<a> Add Binding", clone.addBinding)
	form.AddItem("<e> Edit Bindings", clone.editBinding)
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
	result.sort()
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
	if len(b.Bindings) == 0 {
		return fmt.Errorf("at least one binding must be configured")
	}
	return nil
}

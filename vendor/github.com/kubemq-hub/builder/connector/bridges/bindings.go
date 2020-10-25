package bridges

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Bindings struct {
	Bindings          []*Binding `json:"bindings"`
	defaultOptions    common.DefaultOptions
	takenBindingNames []string
	takenSourceNames  []string
	takenTargetNames  []string
	addressOptions    []string
	defaultName       string
}

func NewBindings(defaultName string) *Bindings {
	return &Bindings{
		defaultName: defaultName,
	}
}

func (b *Bindings) SetDefaultOptions(value common.DefaultOptions) *Bindings {
	b.defaultOptions = value
	return b
}
func (b *Bindings) SetBindings(value []*Binding) *Bindings {
	b.Bindings = value
	return b
}
func (b *Bindings) SetDefaultName(value string) *Bindings {
	b.defaultName = value
	return b
}
func (b *Bindings) confirmBinding(bnd *Binding) bool {
	utils.Println(fmt.Sprintf(promptBindingConfirm, bnd.ColoredYaml()))
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
func (b *Bindings) addBinding() error {

	bnd := NewBinding(fmt.Sprintf("%s-binding-%d", b.defaultName, len(b.Bindings)+1))
	var err error
	if bnd, err = bnd.
		SetAddress(b.addressOptions).
		SetTakenBindingNames(b.takenBindingNames).
		SetTakenSourceNames(b.takenSourceNames).
		SetTakenTargetsNames(b.takenTargetNames).
		Render(); err != nil {
		return err
	}
	ok := b.confirmBinding(bnd)
	if ok {
		b.Bindings = append(b.Bindings, bnd)
		b.takenBindingNames = append(b.takenBindingNames, bnd.BindingName())
		b.takenSourceNames = append(b.takenSourceNames, bnd.SourceName())
		b.takenTargetNames = append(b.takenTargetNames, bnd.TargetName())

	}
	return nil
}
func (b *Bindings) switchOrRemove(old, new *Binding) {
	var newBindingList []*Binding
	var newTakenBindingNames []string
	var newTakenSourceNames []string
	var newTakenTargetNames []string

	for _, binding := range b.Bindings {
		if old.Name != binding.Name {
			newBindingList = append(newBindingList, binding)
			newTakenBindingNames = append(newTakenBindingNames, binding.Name)
			newTakenSourceNames = append(newTakenSourceNames, binding.Sources.Name)
			newTakenTargetNames = append(newTakenTargetNames, binding.Targets.Name)
		}
	}
	if new != nil {
		newBindingList = append(newBindingList, new)
		newTakenBindingNames = append(newTakenBindingNames, new.Name)
		newTakenSourceNames = append(newTakenSourceNames, new.Sources.Name)
		newTakenTargetNames = append(newTakenTargetNames, new.Targets.Name)
	}
	b.Bindings = newBindingList
	b.takenBindingNames = newTakenBindingNames
	b.takenSourceNames = newTakenSourceNames
	b.takenTargetNames = newTakenTargetNames

}
func (b *Bindings) editBinding() error {
	menu := survey.NewMenu("Select Binding to edit").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, binding := range b.Bindings {
		editFn := func() error {
			var err error
			defaultName := binding.Name
			edited := binding.Clone()
			if edited, err = edited.
				SetEditMode(true).
				SetDefaultName(defaultName).
				SetAddress(b.addressOptions).
				SetTakenBindingNames(b.takenBindingNames).
				SetTakenSourceNames(b.takenSourceNames).
				SetTakenTargetsNames(b.takenTargetNames).
				Render(); err != nil {
				return err
			}
			if !edited.wasEdited {
				return nil
			}
			ok := b.confirmBinding(edited)
			if ok {
				b.switchOrRemove(binding, edited)
				utils.Println(promptBindingEditedConfirmation, binding.Name)
			} else {
				utils.Println(promptBindingEditedNoSave, binding.Name)
			}
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
	menu := survey.NewMenu("Select Binding to delete").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		SetDisableLoop(true)
	for _, binding := range b.Bindings {
		deleteFn := func() error {
			bindingName := binding.Name
			val := false
			if err := survey.NewBool().
				SetName("confirm-delete").
				SetMessage(fmt.Sprintf("Are you sure you want to delete %s binding", bindingName)).
				SetRequired(true).
				SetDefault("false").
				Render(&val); err != nil {
				return err
			}
			if val {
				b.switchOrRemove(binding, nil)
				utils.Println(promptBindingDeleteConfirmation, binding.Name)
				return nil
			}
			return nil
		}
		menu.AddItem(binding.Name, deleteFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (b *Bindings) listBindings() error {
	menu := survey.NewMenu("Select Binding to show configuration").
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, binding := range b.Bindings {
		showFn := func() error {
			utils.Println(promptShowBinding, binding.Name)
			utils.Println("%s\n", binding.ColoredYaml())
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
	utils.Println(promptBindingStartMenu)
	for {
		menu := survey.NewMenu("Select Bindings operation").
			SetBackOption(true).
			SetErrorHandler(survey.MenuShowErrorFn)

		menu.AddItem("Add binding", b.addBinding)
		menu.AddItem("Edit binding", b.editBinding)
		menu.AddItem("Delete binding", b.deleteBinding)
		menu.AddItem("List bindings", b.listBindings)
		if err := menu.Render(); err != nil {
			return nil, err
		}

		if len(b.Bindings) == 0 {
			utils.Println(promptBindingEmptyError)
		} else {
			break
		}
	}
	return yaml.Marshal(b)
}

func (b *Bindings) Marshal() ([]byte, error) {
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

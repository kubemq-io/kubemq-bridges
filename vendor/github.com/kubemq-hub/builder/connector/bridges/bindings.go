package bridges

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/uitable"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"sort"
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
func (b *Bindings) askSelectBinding(op string) (*Binding, error) {
	var bindingList []string
	for _, bnd := range b.Bindings {
		bindingList = append(bindingList, bnd.Name)
	}
	bindingList = append(bindingList, "Return")
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("select-binding").
		SetMessage(fmt.Sprintf("Select Binding to %s", op)).
		SetDefault(bindingList[0]).
		SetHelp("Select Binding name to edit, show or delete").
		SetRequired(true).
		SetOptions(bindingList).
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val == "Return" {
		return nil, nil
	}
	for _, binding := range b.Bindings {
		if val == binding.Name {
			return binding, nil
		}
	}
	return nil, nil
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
	bnd, err := b.askSelectBinding("edit")
	if err != nil {
		return err
	}

	if bnd == nil {
		utils.Println(promptBindingEditCanceled)
		return nil
	}

	edited := bnd.Clone()
	if edited, err = edited.
		SetEditMode(true).
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
		b.switchOrRemove(bnd, edited)
		utils.Println(promptBindingEditedConfirmation, bnd.Name)
	} else {
		utils.Println(promptBindingEditedNoSave, bnd.Name)
	}

	return nil
}
func (b *Bindings) deleteBinding() error {
	bnd, err := b.askSelectBinding("delete")
	if err != nil {
		return err
	}
	if bnd == nil {
		utils.Println(promptBindingDeleteCanceled)
		return nil
	}
	b.switchOrRemove(bnd, nil)
	utils.Println(promptBindingDeleteConfirmation, bnd.Name)
	return nil
}
func (b *Bindings) showBinding() error {
	bnd, err := b.askSelectBinding("show")
	if err != nil {
		return err
	}
	if bnd == nil {
		utils.Println(promptBindingShowCanceled)
		return nil
	}
	utils.Println(promptShowBinding, bnd.Name)
	utils.Println(bnd.ColoredYaml())
	return nil
}
func (b *Bindings) showList() error {
	utils.Println(promptShowList)
	table := uitable.New()
	table.MaxColWidth = 80
	rows := b.TableShort()
	for i := 0; i < len(rows); i++ {
		table.AddRow(rows[i]...)
	}
	utils.Println(fmt.Sprintf("%s\n", table.String()))
	return nil
}
func (b *Bindings) askMenu() error {
	utils.Println(promptBindingStartMenu)
	for {
		sort.Slice(b.Bindings, func(i, j int) bool {
			return b.Bindings[i].Name < b.Bindings[j].Name
		})
		var ops []string
		if len(b.Bindings) == 0 {
			ops = []string{
				"Add",
				"RETURN",
			}
		} else {
			ops = []string{
				"Add new binding",
				"Edit existed binding",
				"Show existed binding",
				"Delete existed binding",
				"List of bindings",
				"RETURN",
			}
		}
		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-operation").
			SetMessage("Select Binding operation").
			SetDefault(ops[0]).
			SetHelp("Select Binding operation").
			SetRequired(true).
			SetOptions(ops).
			Render(&val)
		if err != nil {
			return err
		}
		switch val {
		case ops[0]:
			if err := b.addBinding(); err != nil {
				return err
			}
		case ops[1]:
			if err := b.editBinding(); err != nil {
				return err
			}
		case ops[2]:
			if err := b.showBinding(); err != nil {
				return err
			}
		case ops[3]:
			if err := b.deleteBinding(); err != nil {
				return err
			}

		case ops[4]:
			if err := b.showList(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}
func (b *Bindings) Render() ([]byte, error) {
	if err := b.askMenu(); err != nil {
		return nil, err
	}
	if len(b.Bindings) == 0 {
		return nil, fmt.Errorf("at least one binding must be configured")
	}
	return yaml.Marshal(b)
}

func (b *Bindings) Marshal() ([]byte, error) {
	return yaml.Marshal(b)
}
func (b *Bindings) Unmarshal(data []byte) *Bindings {
	bnd := &Bindings{}
	err := yaml.Unmarshal(data, bnd)
	if err != nil {
		return b
	}
	return bnd
}
func (b *Bindings) TableShort() [][]interface{} {
	var rows [][]interface{}
	headers := [][]interface{}{
		{
			"NAME",
			"SOURCES (Name/Kind/Connections)",
			"TARGETS (Name/Kind/Connections)",
			"MIDDLEWARES",
		},
		{
			"----",
			"-------------------------------",
			"-------------------------------",
			"------------",
		},
	}
	rows = append(rows, headers...)
	for _, bnd := range b.Bindings {
		rows = append(rows, bnd.TableRowShort())
	}
	return rows
}

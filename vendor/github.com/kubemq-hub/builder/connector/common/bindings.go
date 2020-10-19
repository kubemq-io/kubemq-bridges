package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/uitable"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"gopkg.in/yaml.v2"
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

func (b *Bindings) SetBindings(value []*Binding) *Bindings {
	b.Bindings = value
	return b
}
func (b *Bindings) SetManifest(value *Manifest) *Bindings {
	b.manifest = value
	return b
}
func (b *Bindings) SetOptions(value DefaultOptions) *Bindings {
	b.loadedOptions = value
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
		SetDefaultOptions(b.loadedOptions).
		SetSourcesList(b.manifest.Sources).
		SetTargetsList(b.manifest.Targets).
		SetTakenBindingNames(b.takenBindingNames).
		Render(); err != nil {
		return err
	}
	ok := b.confirmBinding(bnd)
	if ok {
		b.Bindings = append(b.Bindings, bnd)
		b.takenBindingNames = append(b.takenBindingNames, bnd.Name)
	}

	return nil
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
		SetDefaultOptions(b.loadedOptions).
		SetSourcesList(b.manifest.Sources).
		SetTargetsList(b.manifest.Targets).
		SetTakenBindingNames(b.takenBindingNames).
		Render(); err != nil {
		return err
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
		SetMessage(fmt.Sprintf("Select Binding name to %s", op)).
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
	if b.manifest == nil {
		return nil, fmt.Errorf("inavlid manifest")
	}
	if err := b.askMenu(); err != nil {
		return nil, err
	}

	if len(b.Bindings) == 0 {
		return nil, fmt.Errorf("at least one binding must be configured")
	}
	return yaml.Marshal(b)
}

func (b *Bindings) Yaml() ([]byte, error) {
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
			"SOURCE (Name/Kind)",
			"TARGET (Name/Kind)",
			"MIDDLEWARES",
		},
		{
			"----",
			"-----------------",
			"-----------------",
			"------------",
		},
	}
	rows = append(rows, headers...)
	for _, bnd := range b.Bindings {
		rows = append(rows, bnd.TableRowShort())
	}
	return rows
}

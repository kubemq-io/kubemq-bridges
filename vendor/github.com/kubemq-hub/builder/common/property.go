package common

import "fmt"

type Property struct {
	Name          string   `json:"name"`
	Kind          string   `json:"kind"`
	Description   string   `json:"description"`
	Default       string   `json:"default"`
	Options       []string `json:"options"`
	Must          bool     `json:"must"`
	Min           int      `json:"min"`
	Max           int      `json:"max"`
	Conditional   map[string][]*Property
	LoadedOptions string
}

func NewProperty() *Property {
	return &Property{}
}

func (p *Property) SetName(value string) *Property {
	p.Name = value
	return p
}

func (p *Property) SetKind(value string) *Property {
	p.Kind = value
	return p
}
func (p *Property) SetDescription(value string) *Property {
	p.Description = value
	return p
}

func (p *Property) SetDefault(value string) *Property {
	p.Default = value
	return p
}
func (p *Property) SetLoadedOptions(value string) *Property {
	p.LoadedOptions = value
	return p
}
func (p *Property) SetOptions(value []string) *Property {
	p.Options = value
	return p
}
func (p *Property) SetMust(value bool) *Property {
	p.Must = value
	return p
}

func (p *Property) SetMin(value int) *Property {
	p.Min = value
	return p
}
func (p *Property) SetMax(value int) *Property {
	p.Max = value
	return p
}
func (p *Property) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("property kind %s name cannot be empty", p.Kind)
	}
	if p.Kind == "" {
		return fmt.Errorf("property %s kind cannot be empty", p.Name)
	}

	if p.Description == "" {
		return fmt.Errorf("property %s description cannot be empty", p.Name)
	}

	for cond, pList := range p.Conditional {
		if len(pList) == 0 {
			return fmt.Errorf("condtion %s must have proerties", cond)
		}
		for _, p := range pList {
			if err := p.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
func (p *Property) NewCondition(condition string, properties []*Property) *Property {
	if p.Conditional == nil {
		p.Conditional = map[string][]*Property{}
	}
	p.Conditional[condition] = properties
	return p
}

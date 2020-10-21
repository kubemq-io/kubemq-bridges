package survey

import "fmt"

type Selector interface {
	Key() string
}

type ListSelector struct {
	title    string
	items    []string
	itemsMap map[string]Selector
}

func NewListSelector(title string) *ListSelector {
	return &ListSelector{
		title:    title,
		items:    []string{},
		itemsMap: map[string]Selector{},
	}
}

func (ls *ListSelector) AddItems(items ...Selector) *ListSelector {

	for _, item := range items {
		ls.items = append(ls.items, item.Key())
		ls.itemsMap[item.Key()] = item
	}
	return ls
}

func (ls *ListSelector) Render() (Selector, error) {
	if len(ls.items) == 0 {
		return nil, nil
	}
	val := ""
	err := NewString().
		SetKind("string").
		SetName("menu").
		SetMessage(ls.title).
		SetDefault(ls.items[0]).
		SetRequired(true).
		SetOptions(ls.items).
		Render(&val)
	if err != nil {
		return nil, err
	}
	s, ok := ls.itemsMap[val]
	if !ok {
		return nil, fmt.Errorf("obejct for %s not found", val)
	}
	return s, nil
}

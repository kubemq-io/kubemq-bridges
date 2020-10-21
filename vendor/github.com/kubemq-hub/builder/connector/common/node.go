package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Node struct {
	Name      string
	Childes   []*Node
	Connector *Connector
}

func NewNode() *Node {
	return &Node{}
}
func BuildNextNode(node *Node, kinds []string, connector *Connector) *Node {
	if len(kinds) == 0 {
		return nil
	}
	if len(kinds) == 1 {
		return &Node{
			Name:      kinds[0],
			Childes:   nil,
			Connector: connector,
		}
	}

	var child *Node
	for _, currentChild := range node.Childes {
		if currentChild.Name == kinds[0] {
			child = currentChild
			break
		}
	}
	if child == nil {
		child = &Node{
			Name:      kinds[0],
			Childes:   []*Node{},
			Connector: nil,
		}
		child.Childes = append(child.Childes, BuildNextNode(child, kinds[1:], connector))
	} else {
		child.Childes = append(child.Childes, BuildNextNode(child, kinds[1:], connector))
	}
	return child
}
func (n *Node) Render() error {
	if len(n.Childes) == 0 && n.Connector != nil {
		utils.Println("<red>%s</>", n.Connector.Kind)
		return nil
	}
	menu := survey.NewMenu(fmt.Sprintf("Browse %s catalog", n.Name)).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, child := range n.Childes {
		menu.AddItem(child.Name, child.Render)
	}
	if err := menu.Render(); err != nil {
		return nil
	}
	return nil
}

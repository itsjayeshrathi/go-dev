package exercise

import (
	"fmt"
	"slices"
)

type Item struct {
	Name string
	Type string
}

type Player struct {
	Name      string
	Inventory []Item
}

func (p *Player) PickUpItem(item Item) {
	p.Inventory = append(p.Inventory, item)
}

func (p *Player) RemoveItem(name string) {
	index := -1
	for i, value := range p.Inventory {
		if value.Name == name {
			index = i
			break
		}
	}
	if index != -1 {
		p.Inventory = slices.Delete(p.Inventory, index, index+1)
	}

}

func (p *Player) UseItem(name string) {
	for _, value := range p.Inventory {
		if value.Name == name {
			fmt.Printf("You are using %s of type %v\n", value.Name, value.Type)
			break
		}
	}
}

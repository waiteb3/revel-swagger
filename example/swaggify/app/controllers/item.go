package controllers

import "github.com/revel/revel"

type ItemController struct {
	*revel.Controller
}

type Item struct {
	Name string
}

var Items = make(map[string]Item)

// List all items
func (c ItemController) List() revel.Result {
	return c.RenderJson(Items)
}

// Create a new item
func (c ItemController) Create(name string) revel.Result {
	item := Item{name}
	Items[item.Name] = item
	return c.RenderJson(item)
}

// Get an item
func (c ItemController) Get(name string) revel.Result {
	if item, exists := Items[name]; exists {
		return c.RenderJson(item)
	} else {
		return c.NotFound("item not found")
	}
}

// Update an item by appending an exlamation point to the name
func (c ItemController) Update(name string) revel.Result {
	if item, exists := Items[name]; exists {
		delete(Items, name)
		name += "!"
		item.Name = name
		Items[name] = item
		return c.RenderJson(item)
	} else {
		return c.NotFound("item not found")
	}
}

// Delete an item from the pool
func (c ItemController) Delete(name string) revel.Result {
	if item, exists := Items[name]; exists {
		delete(Items, name)
		return c.RenderJson(item)
	} else {
		return c.NotFound("item not found")
	}
}

// Modify an item by appending an exlamation point to the name
func (c ItemController) Modify(name string) revel.Result {
	if item, exists := Items[name]; exists {
		delete(Items, name)
		name += "!"
		item.Name = name
		Items[name] = item
		return c.RenderJson(item)
	} else {
		return c.NotFound("item not found")
	}
}

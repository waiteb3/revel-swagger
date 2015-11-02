package controllers

import (
	"github.com/revel/revel"
	"sort"
)

type Items struct {
	*revel.Controller
}

type Item struct {
	Id   int    `json:"id"`
	Name string `json:"id"`
}

type ById []*Item

func (list ById) Len() int { return len(list) }

func (list ById) Less(i, j int) bool { return list[i].Id < list[j].Id }

func (list ById) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

// Add an item to the list
func (list ById) Add(item *Item) {
	if len(list) == cap(list) {
		var newList ById = make([]*Item, len(list), cap(list)*2)
		copy(newList, list)
		list = newList
	}
	list = append(list, item)
}

// Get an item from the list with Id `id`
// returns nil if not found
func (list ById) Get(id int) *Item {
	i := sort.Search(len(list), func(i int) bool { return list[i].Id == id })
	if i < len(list) && list[i].Id == id {
		return list[i]
	}
	return nil
}

// Delete an item from the list with Id `id`
// returns the deleted item or nil if not found
func (list ById) Delete(id int) *Item {
	i := sort.Search(len(list), func(i int) bool { return list[i].Id == id })
	if i < len(list) && list[i].Id == id {
		item := list[i]
		list = append(list[:i], list[i+1:]...)
		return item
	}
	return nil
}

var ItemList ById = make([]*Item, 0, 16)

// List the current set of items
func (c Items) List() revel.Result {
	return c.RenderJson(ItemList)
}

// Create an item
func (c Items) Create(id int, name string) revel.Result {
	ItemList.Add(&Item{id, name})
	sort.Sort(ById(ItemList))
	return c.RenderJson(ItemList.Get(id))
}

// Read an item
func (c Items) Read(id int) revel.Result {
	item := ItemList.Get(id)
	if item == nil {
		return c.NotFound("Item with id '%d' not found.", id)
	}
	return c.RenderJson(item)
}

// Update an item
func (c Items) Update(id int, name string) revel.Result {
	item := ItemList.Get(id)
	if item == nil {
		return c.NotFound("Item with id '%d' not found.", id)
	}
	item.Name = name
	return c.RenderJson(item)
}

// Delete an item
func (c Items) Delete(id int) revel.Result {
	// since this is an array, deleting an object doesn't affect sorting
	item := ItemList.Delete(id)
	if item == nil {
		return c.NotFound("Item with id '%d' not found.", id)
	}
	return c.RenderJson(item)
}

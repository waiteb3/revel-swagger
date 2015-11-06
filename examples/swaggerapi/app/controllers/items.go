package controllers

import (
	"sort"

	"github.com/revel/revel"
)

type Items struct {
	*revel.Controller
}

type Item struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ById []Item

func (list ById) Len() int { return len(list) }

func (list ById) Less(i, j int) bool { return list[i].Id < list[j].Id }

func (list ById) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

// Find an item from the list with Id `id` and the index
// -1 for not found
func (list ById) Find(id int) (Item, int) {
	i := list.Search(id)
	if i < len(list) && list[i].Id == id {
		return list[i], i
	}
	return Item{}, -1
}

// Search for an item by id and return the index
func (list ById) Search(id int) int {
	return sort.Search(len(list), func(i int) bool {
		return list[i].Id == id
	})
}

var ItemList ById = make([]Item, 0, 1)

// List the current set of items
func (c Items) List() revel.Result {
	return c.RenderJson(ItemList)
}

// Create an item
func (c Items) Create(item *Item) revel.Result {
	// 	var item Item
	// 	err := json.NewDecoder(c.Request.Body).Decode(&item)
	// if err != nil {
	// 	c.RenderError(err)
	// }

	ItemList = append(ItemList, *item)
	sort.Sort(ById(ItemList))
	return c.RenderJson(item)
}

// Read an item
func (c Items) Read(id int) revel.Result {
	item, i := ItemList.Find(id)
	if i < 0 {
		return c.NotFound("Item with id %d not found.", id)
	}
	return c.RenderJson(item)
}

// Update an item
func (c Items) Update(id int, name string) revel.Result {
	item, i := ItemList.Find(id)
	if i < 0 {
		return c.NotFound("Item with id %d not found.", id)
	}
	item.Name = name
	ItemList[i] = item
	sort.Sort(ById(ItemList))
	return c.RenderJson(item)
}

// Delete an item
func (c Items) Delete(id int) revel.Result {
	// since this is an array, deleting an object doesn't affect sorting
	if item, i := ItemList.Find(id); i < 0 {
		return c.NotFound("Item with id %d not found.", id)
	} else {
		ItemList = append(ItemList[:i], ItemList[i+1:]...)
		return c.RenderJson(item)
	}
}

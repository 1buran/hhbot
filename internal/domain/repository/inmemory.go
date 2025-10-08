package repository

import (
	"sync"
)

// In memory repository holds objects in memory.
type inmemory[T Entities] struct {
	items map[string]T // key value storage: id -> item

	sync.RWMutex
}

func (i *inmemory[T]) Add(v T) {
	i.Lock()
	defer i.Unlock()

	i.items[v.GetID()] = v
}

func (i *inmemory[T]) List() (items []T) {
	i.RLock()
	defer i.RUnlock()

	for _, v := range i.items {
		items = append(items, v)
	}
	return
}

func NewInMemoryRepository[T Entities]() Repository[T] {
	return &inmemory[T]{items: make(map[string]T)}
}

package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type entry struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, ok := cache.items[key]

	if ok {
		item.Value = entry{key: key, value: value}
		cache.queue.MoveToFront(item)
	} else {
		cache.queue.PushFront(entry{key: key, value: value})
		cache.items[key] = cache.queue.Front()
	}

	if cache.queue.Len() > cache.capacity {
		back := cache.queue.Back()
		cache.queue.Remove(back)
		delete(cache.items, back.Value.(entry).key)
	}

	return ok
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if item, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(item)
		return item.Value.(entry).value, true
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.queue = NewList()
	clear(cache.items)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

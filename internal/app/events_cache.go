package app

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

//кэш для хранения типа незавершенного события и его айди
type Cache struct {
	mu sync.Mutex
	m  map[string]primitive.ObjectID
}

func NewCache() *Cache {
	var c Cache
	c.m = make(map[string]primitive.ObjectID)
	c.mu = sync.Mutex{}
	return &c
}

func (c *Cache) Set(key string, value primitive.ObjectID) {
	c.mu.Lock()
	c.m[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (primitive.ObjectID, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.m[key]
	return value, ok
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.m, key)
	c.mu.Unlock()
}

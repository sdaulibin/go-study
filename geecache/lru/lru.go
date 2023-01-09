package lru

import "container/list"

type Cache struct {
	maxBytes int64
	nbytes int64
	ll *list.List
	cahce map[string]*list.Element
	OnEvicted func(key string,value Value)
}

type entry struct {
	key string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64,onEvicted func(string,Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cahce:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) Add(key string,value Value) {
	if element,ok := c.cahce[key]; ok {
		c.ll.MoveToBack(element)
		kv := element.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
		return
	}
	element := c.ll.PushBack(&entry{key,value})
	c.cahce[key] = element
	c.nbytes += int64(len(key)) + int64(value.Len())

	for c.maxBytes !=0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value,ok bool) {
	if element,ok := c.cahce[key] ; ok {
		c.ll.MoveToBack(element)
		kv := element.Value.(*entry)
		return kv.value,true
	}
	return
}

func (c *Cache) RemoveOldest() {
	element := c.ll.Front()
	if element!= nil {
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cahce,kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key,kv.value)
		}
	}
}

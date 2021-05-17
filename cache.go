//Cache implemented by LRU
package lrucache

import (
	"container/list"
	"sync"
)

/**
*@Author lyer
*@Date 3/17/21 20:10
*@Describe
**/

//Cache
type Cache struct {
	cap      int64
	size     int64
	ll       *list.List
	m        map[string]*list.Element
	OnRemove func(e *entry) //callback when remove key
	rwm      sync.RWMutex
}

type entry struct {
	key string
	val string
}

//New create a cache
func New(cap int64, f func(e *entry)) *Cache {
	return &Cache{
		cap:      cap,
		ll:       list.New(),
		m:        make(map[string]*list.Element),
		OnRemove: f,
	}
}

//RemoveOldest delete the least used element
func (c *Cache) removeOldest() {
	if ele := c.ll.Back(); ele != nil {
		c.removeElement(ele)
	}
}

//Add key and value
func (c *Cache) Add(key string, val string) {
	//如果不存在
	if ok := c.Update(key, val); !ok {
		c.rwm.Lock()
		defer c.rwm.Unlock()
		ele := c.ll.PushFront(&entry{key, val})
		c.size += int64(len(val)) + int64(len(key))
		c.m[key] = ele
	}
	for c.cap > 0 && c.cap < c.size {
		c.removeOldest()
	}
}

func (c *Cache) removeElement(ele *list.Element) {
	c.ll.Remove(ele)
	kv := ele.Value.(*entry)
	delete(c.m, kv.key)
	c.size -= int64(len(kv.key) + len(kv.val))
	if c.OnRemove != nil {
		c.OnRemove(kv)
	}
}

//Remove key
func (c *Cache) Remove(key string) bool {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	if ele, ok := c.m[key]; ok {
		c.removeElement(ele)
	}
	return false
}

//Update key
func (c *Cache) Update(key string, newVal string) bool {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	if ele, ok := c.m[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.size = c.size - int64(len(kv.val)) + int64(len(newVal))
		kv.val = newVal
		return true
	}
	return false
}

//Get key
func (c *Cache) Get(key string) (val string, ok bool) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	if ele, ok := c.m[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.val, true
	}
	return
}

//Keys get all keys
func (c *Cache) Keys() []string {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	keys := make([]string, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	return keys
}

//Values get all values
func (c *Cache) Values() []string {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	values := make([]string, 0, len(c.m))
	for _, ele := range c.m {
		val := ele.Value.(*entry).val
		values = append(values, val)
	}
	return values
}

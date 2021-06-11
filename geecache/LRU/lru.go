package lru

import (
	"container/list"
)

//valueList 双向链表,cache map保存映射关系
//Cache 维护缓存大小和映射关系
type Cache struct {
	maxBytes   	 int64 //最大可用内存
	nBytes		 int64
	valueList 	*list.List
	cache 	 	map[string]*list.Element
	OnEvited 	func(key string, value Value)
}

//entry 是具体的链表中的值 Value
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

//实例化Cache
func CacheNew(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		valueList: list.New(), // create linked list
		cache: make(map[string]*list.Element),
		OnEvited: onEvicted,
	}
}

//获取cache指定key的记录
func (c *Cache)CacheGet(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key];ok{
		c.valueList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//淘汰cache记录
func (c *Cache)RemoveOldest(){
	ele := c.valueList.Back()
	if ele != nil {
		c.valueList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvited != nil {
			c.OnEvited(kv.key, kv.value)
		}
	}
}

//增加记录到cache
func (c *Cache)CacheAdd(key string, value Value) {
	if ele, ok := c.cache[key];ok {
		//移动已经存在的对应节点值到队列尾部
		c.valueList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.valueList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

//获取双向链表的长度
func (c *Cache)Len() int {
	return c.valueList.Len()
}
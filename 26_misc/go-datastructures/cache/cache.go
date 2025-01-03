// 允许在内存中存储若干 Item（带有 Size() 方法，表示其占用的字节数）。
// 它有一个容量上限 (cap)，在新插入的 Item 总大小即将超出上限时，
// 会根据指定的缓存淘汰策略（Eviction Policy）进行自动驱逐（Eviction）。
//
// 支持的淘汰策略：
//
// LRA (Least Recently Added)：最早加入缓存的项目优先被驱逐。
// LRU (Least Recently Used)：最近最少使用的项目被驱逐（默认策略）

package main

import (
	"container/list"
	"fmt"
	"sync"
)

// 自定义一个简单的 Item 类型，以字符串长度作为其 Size()
type stringItem struct {
	s string
}

// 实现 Item 接口
func (si stringItem) Size() uint64 {
	return uint64(len(si.s))
}

func main() {
	// 1. 创建缓存，容量 20 字节，默认 LRU
	c := NewCache(20)

	// 2. 放入一些数据
	c.Put("apple", stringItem{"apple"})   // Size=5
	c.Put("banana", stringItem{"banana"}) // Size=6
	c.Put("foo", stringItem{"foo"})       // Size=3

	// 缓存此时已有 items: {apple, banana, foo}, size = 14

	// 3. Get 访问
	got := c.Get("banana", "apple", "missing") // "missing" 不存在
	fmt.Printf("Get banana, apple, missing => %+v\n", got)
	// got[0] = stringItem{"banana"}
	// got[1] = stringItem{"apple"}
	// got[2] = nil

	// 4. 再 Put 一个值，看看是否会触发淘汰
	//    c.Size() 现在是14，再放入 size=10 的 "abcdefghij" => total=24 > cap(20)
	c.Put("big", stringItem{"abcdefghij"})

	// ensureCapacity 会驱逐最久未使用的，基于 LRU => "foo" (因为最新访问过 apple/banana)
	//   "foo" 占3字节，腾出后还有 21 > 20 -> 依然不够
	//   继续驱逐下一个最久未使用 => "banana" (6字节), 21-6=15 <= 20 -> OK
	//
	// 剩余 apple(5字节), big(10字节) => size=15

	// 5. 检查最终有哪些 key
	items := c.Get("apple", "banana", "foo", "big")
	//   apple => still exist, size=5
	//   banana => nil (被驱逐)
	//   foo => nil (被驱逐)
	//   big => size=10

	fmt.Println("After inserting 'big':", items)

	// 6. 打印当前总大小
	fmt.Println("Cache size =", c.Size()) // 15

	// 7. 删除某些 key
	c.Remove("apple", "big")
	fmt.Println("Cache size after remove =", c.Size()) // 0
}

// Cache is a bounded-size in-memory cache of sized items with a configurable eviction policy
type Cache interface {
	// Get retrieves items from the cache by key.
	// If an item for a particular key is not found, its position in the result will be nil.
	Get(keys ...string) []Item

	// Put adds an item to the cache.
	Put(key string, item Item)

	// Remove clears items with the given keys from the cache
	Remove(keys ...string)

	// Size returns the size of all items currently in the cache.
	Size() uint64
}

// Item is an item in a cache
type Item interface {
	// Size returns the item's size, in bytes
	Size() uint64
}

// A tuple tracking a cached item and a reference to its node in the eviction list
type cached struct {
	item    Item
	element *list.Element
}

// Sets the provided list element on the cached item if it is not nil
func (c *cached) setElementIfNotNil(element *list.Element) {
	if element != nil {
		c.element = element
	}
}

// Private cache implementation
type cache struct {
	sync.Mutex                                  // Lock for synchronizing Get, Put, Remove
	cap          uint64                         // Capacity bound
	size         uint64                         // Cumulative size
	items        map[string]*cached             // Map from keys to cached items
	keyList      *list.List                     // List of cached items in order of increasing evictability
	recordAdd    func(key string) *list.Element // Function called to indicate that an item with the given key was added
	recordAccess func(key string) *list.Element // Function called to indicate that an item with the given key was accessed
}

// CacheOption configures a cache.
type CacheOption func(*cache)

// Policy is a cache eviction policy for use with the EvictionPolicy CacheOption.
type Policy uint8

const (
	// LeastRecentlyAdded indicates a least-recently-added eviction policy.
	LeastRecentlyAdded Policy = iota
	// LeastRecentlyUsed indicates a least-recently-used eviction policy.
	LeastRecentlyUsed
)

// EvictionPolicy sets the eviction policy to be used to make room for new items.
// If not provided, default is LeastRecentlyUsed.
func EvictionPolicy(policy Policy) CacheOption {
	return func(c *cache) {
		switch policy {
		case LeastRecentlyAdded:
			c.recordAccess = c.noop // 访问时不改变顺序
			c.recordAdd = c.record  // 添加时将其移动到链表头
		case LeastRecentlyUsed:
			c.recordAccess = c.record // 访问时将其移动到链表头
			c.recordAdd = c.noop      // 添加时不改变顺序（实际立即放在头也可，只是函数机制不同）
		}
	}
}

// NewCache returns a cache with the requested options configured.
// The cache consumes memory bounded by a fixed capacity,
// plus tracking overhead linear in the number of items.
func NewCache(capacity uint64, options ...CacheOption) Cache {
	c := &cache{
		cap:     capacity,
		keyList: list.New(),
		items:   map[string]*cached{},
	}
	// Default LRU eviction policy
	EvictionPolicy(LeastRecentlyUsed)(c)

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *cache) Get(keys ...string) []Item {
	c.Lock()
	defer c.Unlock()

	items := make([]Item, len(keys))
	for i, key := range keys {
		cached := c.items[key]
		if cached == nil {
			items[i] = nil
		} else {
			c.recordAccess(key)
			items[i] = cached.item
		}
	}

	return items
}

func (c *cache) Put(key string, item Item) {
	c.Lock()
	defer c.Unlock()

	// Remove the item currently with this key (if any)
	c.remove(key)

	// Make sure there's room to add this item
	c.ensureCapacity(item.Size())

	// Actually add the new item
	cached := &cached{item: item}
	cached.setElementIfNotNil(c.recordAdd(key))
	cached.setElementIfNotNil(c.recordAccess(key))
	c.items[key] = cached
	c.size += item.Size()
}

func (c *cache) Remove(keys ...string) {
	c.Lock()
	defer c.Unlock()

	for _, key := range keys {
		c.remove(key)
	}
}

func (c *cache) Size() uint64 {
	c.Lock()
	defer c.Unlock()

	return c.size
}

// Given the need to add some number of new bytes to the cache,
// evict items according to the eviction policy until there is room.
// The caller should hold the cache lock.
func (c *cache) ensureCapacity(toAdd uint64) {
	mustRemove := int64(c.size+toAdd) - int64(c.cap)
	for mustRemove > 0 {
		key := c.keyList.Back().Value.(string)
		mustRemove -= int64(c.items[key].item.Size())
		c.remove(key)
	}
}

// Remove the item associated with the given key.
// The caller should hold the cache lock.
func (c *cache) remove(key string) {
	if cached, ok := c.items[key]; ok {
		delete(c.items, key)
		c.size -= cached.item.Size()
		c.keyList.Remove(cached.element)
	}
}

// A no-op function that does nothing for the provided key
func (c *cache) noop(string) *list.Element { return nil }

// A function to record the given key and mark it as last to be evicted
func (c *cache) record(key string) *list.Element {
	if item, ok := c.items[key]; ok {
		c.keyList.MoveToFront(item.element)
		return item.element
	}
	return c.keyList.PushFront(key)
}

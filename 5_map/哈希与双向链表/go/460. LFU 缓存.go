// 432. 全 O(1) 的数据结构
// https://leetcode.cn/problems/all-oone-data-structure/description/?envType=problem-list-v2&envId=design
// 请你设计一个用于存储字符串计数的数据结构，并能够返回计数最小和最大的字符串。
//
// 实现 AllOne 类：
//
// AllOne() 初始化数据结构的对象。
// inc(String key) 字符串 key 的计数增加 1 。如果数据结构中尚不存在 key ，那么插入计数为 1 的 key 。
// dec(String key) 字符串 key 的计数减少 1 。如果 key 的计数在减少后为 0 ，那么需要将这个 key 从数据结构中删除。测试用例保证：在减少计数前，key 存在于数据结构中。
// getMaxKey() 返回任意一个计数最大的字符串。如果没有元素存在，返回一个空字符串 "" 。
// getMinKey() 返回任意一个计数最小的字符串。如果没有元素存在，返回一个空字符串 "" 。
// 注意：每个函数都应当满足 O(1) 平均时间复杂度。

package main

import "container/list"

type entry struct {
	key, value, freq int
}

type LFUCache struct {
	capacity   int
	minFreq    int
	keyToNode  map[int]*list.Element
	freqToList map[int]*list.List
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		capacity:   capacity,
		keyToNode:  make(map[int]*list.Element),
		freqToList: make(map[int]*list.List),
	}
}

func (c *LFUCache) Get(key int) int {
	if e := c.getEntry(key); e != nil {
		return e.value
	}
	return -1
}

func (c *LFUCache) Put(key int, value int) {
	if e := c.getEntry(key); e != nil {
		e.value = value
		return
	}
	if len(c.keyToNode) == c.capacity {
		lst := c.freqToList[c.minFreq]
		delete(c.keyToNode, lst.Remove(lst.Back()).(*entry).key)
		if lst.Len() == 0 {
			delete(c.freqToList, c.minFreq)
		}
	}
	c.pushFront(&entry{key: key, value: value, freq: 1}) // 新书放在「看过 1 次」的最上面
	c.minFreq = 1
}

func (c *LFUCache) pushFront(e *entry) {
	if _, ok := c.freqToList[e.freq]; !ok {
		c.freqToList[e.freq] = list.New()
	}
	c.keyToNode[e.key] = c.freqToList[e.freq].PushFront(e)
}

func (c *LFUCache) getEntry(key int) *entry {
	node := c.keyToNode[key]
	if node == nil {
		return nil
	}
	e := node.Value.(*entry)
	lst := c.freqToList[e.freq]
	lst.Remove(node)
	if lst.Len() == 0 {
		delete(c.freqToList, e.freq)
		if c.minFreq == e.freq {
			c.minFreq++
		}
	}
	e.freq++
	c.pushFront(e) // 放在右边这摞书的最上面
	return e
}

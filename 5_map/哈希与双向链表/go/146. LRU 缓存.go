// 146. LRU 缓存
// https://leetcode.cn/problems/lru-cache/description/
// 请你设计并实现一个满足  LRU (最近最少使用，Least Recently Used) 缓存 约束的数据结构。
// 实现 LRUCache 类：
// LRUCache(int capacity) 以 正整数 作为容量 capacity 初始化 LRU 缓存
// int get(int key) 如果关键字 key 存在于缓存中，则返回关键字的值，否则返回 -1 。
// void put(int key, int value) 如果关键字 key 已经存在，则变更其数据值 value ；如果不存在，则向缓存中插入该组 key-value 。如果插入操作导致关键字数量超过 capacity ，则应该 逐出 最久未使用的关键字。
// 函数 get 和 put 必须以 O(1) 的平均时间复杂度运行。

package main

import "container/list"

type entry struct {
	key, value int
}

type LRUCache struct {
	capacity  int
	list      *list.List
	keyToNode map[int]*list.Element
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity:  capacity,
		list:      list.New(),
		keyToNode: make(map[int]*list.Element),
	}
}

func (c *LRUCache) Get(key int) int {
	node := c.keyToNode[key]
	if node == nil {
		return -1
	}
	c.list.MoveToFront(node)
	return node.Value.(entry).value
}

func (c *LRUCache) Put(key int, value int) {
	if node := c.keyToNode[key]; node != nil {
		node.Value = (entry{key, value})
		c.list.MoveToFront(node)
		return
	}
	c.keyToNode[key] = c.list.PushFront(entry{key, value})
	if len(c.keyToNode) > c.capacity {
		delete(c.keyToNode, c.list.Remove(c.list.Back()).(entry).key)
	}
}

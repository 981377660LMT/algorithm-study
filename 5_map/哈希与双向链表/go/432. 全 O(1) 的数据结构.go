// 432. 全 O(1) 的数据结构
// https://leetcode.cn/problems/all-oone-data-structure/description/?envType=problem-list-v2&envId=design
// !请你设计一个用于存储字符串计数的数据结构，并能够返回计数最小和最大的字符串。
// 实现 AllOne 类：
// AllOne() 初始化数据结构的对象。
// inc(String key) 字符串 key 的计数增加 1 。如果数据结构中尚不存在 key ，那么插入计数为 1 的 key 。
// dec(String key) 字符串 key 的计数减少 1 。如果 key 的计数在减少后为 0 ，那么需要将这个 key 从数据结构中删除。测试用例保证：在减少计数前，key 存在于数据结构中。
// getMaxKey() 返回任意一个计数最大的字符串。如果没有元素存在，返回一个空字符串 "" 。
// getMinKey() 返回任意一个计数最小的字符串。如果没有元素存在，返回一个空字符串 "" 。
// 注意：每个函数都应当满足 O(1) 平均时间复杂度。
//
// !链表的每个节点维护count.

package main

import "container/list"

type bucket struct {
	keys  map[string]struct{}
	count int
}

type AllOne struct {
	*list.List
	keyToNode map[string]*list.Element
}

func Constructor() AllOne {
	return AllOne{List: list.New(), keyToNode: make(map[string]*list.Element)}
}

func (l *AllOne) Inc(key string) {
	if curNode := l.keyToNode[key]; curNode != nil {
		curBucket := curNode.Value.(*bucket)
		if nextNode := curNode.Next(); nextNode == nil || nextNode.Value.(*bucket).count > curBucket.count+1 {
			l.keyToNode[key] = l.InsertAfter(&bucket{keys: map[string]struct{}{key: {}}, count: curBucket.count + 1}, curNode)
		} else {
			nextNode.Value.(*bucket).keys[key] = struct{}{}
			l.keyToNode[key] = nextNode
		}
		delete(curBucket.keys, key)
		if len(curBucket.keys) == 0 {
			l.Remove(curNode)
		}
	} else {
		if l.Front() == nil || l.Front().Value.(*bucket).count > 1 {
			l.keyToNode[key] = l.PushFront(&bucket{keys: map[string]struct{}{key: {}}, count: 1})
		} else {
			l.Front().Value.(*bucket).keys[key] = struct{}{}
			l.keyToNode[key] = l.Front()
		}
	}
}

func (l *AllOne) Dec(key string) {
	curNode := l.keyToNode[key]
	curBucket := curNode.Value.(*bucket)
	if curBucket.count > 1 {
		if pre := curNode.Prev(); pre == nil || pre.Value.(*bucket).count < curBucket.count-1 {
			l.keyToNode[key] = l.InsertBefore(&bucket{keys: map[string]struct{}{key: {}}, count: curBucket.count - 1}, curNode)
		} else {
			pre.Value.(*bucket).keys[key] = struct{}{}
			l.keyToNode[key] = pre
		}
	} else {
		delete(l.keyToNode, key)
	}
	delete(curBucket.keys, key)
	if len(curBucket.keys) == 0 {
		l.Remove(curNode)
	}
}

func (l *AllOne) GetMaxKey() string {
	if b := l.Back(); b != nil {
		for k := range b.Value.(*bucket).keys {
			return k
		}
	}
	return ""
}

func (l *AllOne) GetMinKey() string {
	if f := l.Front(); f != nil {
		for k := range f.Value.(*bucket).keys {
			return k
		}
	}
	return ""
}

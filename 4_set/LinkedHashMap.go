package main

import (
	"container/list"
	"fmt"
	"strings"
)

func main() {
	mp := NewLinkedHashMap(10)
	mp.Set(1, 2)
	mp.Set(2, 3)
	mp.Set(3, 4)
	mp.Delete(2)
	fmt.Println(mp)
	mp.Set(2, 5)
	fmt.Println(mp.Get(2))
}

type listItem = struct{ key, value interface{} }
type LinkedHashMap struct {
	mp   map[interface{}]*list.Element
	list *list.List
}

func NewLinkedHashMap(initCapacity int) *LinkedHashMap {
	return &LinkedHashMap{make(map[interface{}]*list.Element, initCapacity), list.New()}
}

func (s *LinkedHashMap) Get(key interface{}) (res interface{}, ok bool) {
	if ele, hit := s.mp[key]; hit {
		return ele.Value.(listItem).value, true
	}
	return
}

func (s *LinkedHashMap) Set(key, value interface{}) *LinkedHashMap {
	if ele, hit := s.mp[key]; hit {
		ele.Value = listItem{key, value}
	} else {
		s.mp[key] = s.list.PushBack(listItem{key, value})
	}
	return s
}

func (s *LinkedHashMap) Delete(key interface{}) bool {
	if ele, hit := s.mp[key]; hit {
		s.list.Remove(ele)
		delete(s.mp, key)
		return true
	}
	return false
}

func (s *LinkedHashMap) Has(key interface{}) bool {
	_, ok := s.mp[key]
	return ok
}

func (s *LinkedHashMap) Size() int {
	return len(s.mp)
}

// 按照插入顺序遍历哈希表中的元素
// 当 f 返回 true 时停止遍历
func (s *LinkedHashMap) ForEach(f func(key interface{}, value interface{}) bool) {
	for node := s.list.Front(); node != nil; node = node.Next() {
		if f(node.Value.(listItem).key, node.Value.(listItem).value) {
			break
		}
	}
}

func (s *LinkedHashMap) String() string {
	res := []string{"LinkedHashMap{"}
	content := []string{}
	s.ForEach(func(key interface{}, value interface{}) bool {
		content = append(content, fmt.Sprintf("%v: %v", key, value))
		return false
	})
	res = append(res, strings.Join(content, ", "))
	res = append(res, "}")
	return strings.Join(res, "")
}

package main

import (
	"container/list"
	"fmt"
	"strings"
)

func main() {
	s := NewLinkedHashSet(10)
	s.Add(1)
	s.Add(2)
	fmt.Println(s.Has(1))
	s.Delete(1)
	fmt.Println(s)
}

type LinkedHashSet struct {
	set  map[interface{}]*list.Element
	list *list.List
}

func NewLinkedHashSet(initCapacity int) *LinkedHashSet {
	return &LinkedHashSet{make(map[interface{}]*list.Element, initCapacity), list.New()}
}

func (s *LinkedHashSet) Add(v interface{}) *LinkedHashSet {
	s.set[v] = s.list.PushBack(v)
	return s
}

func (s *LinkedHashSet) Delete(v interface{}) bool {
	if _, ok := s.set[v]; ok {
		s.list.Remove(s.set[v])
		return true
	}
	return false
}

func (s *LinkedHashSet) Has(v interface{}) bool {
	_, ok := s.set[v]
	return ok
}

func (s *LinkedHashSet) Size() int {
	return len(s.set)
}

// 按照插入顺序遍历哈希表中的元素
// 当 f 返回 true 时停止遍历
func (s *LinkedHashSet) ForEach(f func(key interface{}) bool) {
	for node := s.list.Front(); node != nil; node = node.Next() {
		if f(node.Value) {
			break
		}
	}
}

func (s *LinkedHashSet) String() string {
	res := []string{"LinkedHashSet{"}
	content := []string{}
	s.ForEach(func(key interface{}) bool {
		content = append(content, fmt.Sprintf("%v", key))
		return false
	})
	res = append(res, strings.Join(content, ", "))
	res = append(res, "}")
	return strings.Join(res, "")
}

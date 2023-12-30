package main

import (
	"fmt"
	"sort"
)

func main() {
	arr := NewPartialPersistentArray([]int{1, 2, 3, 4, 5})
	v1 := arr.Set(0, 10)
	fmt.Println(arr.Get(v1, 0))
}

type item = struct{ version, value int }

// 部分可持久化数组.
// 只有最新的版本可以修改，旧版本只能访问.
type PartialPersistentArray struct {
	actions [][]*item
	version int
}

func NewPartialPersistentArray(nums []int) *PartialPersistentArray {
	actions := make([][]*item, len(nums))
	for i := 0; i < len(nums); i++ {
		actions[i] = append(actions[i], &item{version: 0, value: nums[i]})
	}
	return &PartialPersistentArray{actions: actions, version: 0}
}

func (p *PartialPersistentArray) Set(index, value int) (newVersion int) {
	p.version++
	p.actions[index] = append(p.actions[index], &item{version: p.version, value: value})
	return p.version
}

func (p *PartialPersistentArray) Get(version, index int) int {
	pos := sort.Search(len(p.actions[index]), func(i int) bool {
		return p.actions[index][i].version > version
	})
	return p.actions[index][pos-1].value
}

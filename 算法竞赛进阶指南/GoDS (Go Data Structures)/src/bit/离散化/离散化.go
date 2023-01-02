package main

import (
	"fmt"
	"sort"
)

func main() {
	D := NewDiscretizer()
	D.Add(1)
	D.Add(2)
	D.Add(30)
	D.Add(4)
	D.Build()
	fmt.Println(D.Get(11))
	fmt.Println(D.Get(1), D.Len())
}

// 离散化模板
type Discretizer struct {
	set     map[int]struct{}
	mp      map[int]int // num -> index (1-based)
	allNums []int
}

func NewDiscretizer() *Discretizer {
	return &Discretizer{
		set: make(map[int]struct{}, 16),
		mp:  make(map[int]int, 16),
	}
}

func (d *Discretizer) Add(num int) {
	d.set[num] = struct{}{}
}

func (d *Discretizer) Build() {
	allNums := make([]int, 0, len(d.set))
	for num := range d.set {
		allNums = append(allNums, num)
	}
	sort.Ints(allNums)
	d.allNums = allNums
	for i, num := range allNums {
		d.mp[num] = i + 1
	}
}

func (d *Discretizer) Get(num int) int {
	res, ok := d.mp[num]
	if ok {
		return res
	}
	// bisect right
	return sort.Search(len(d.allNums), func(i int) bool { return d.allNums[i] > num })
}

func (d *Discretizer) Len() int {
	return len(d.allNums)
}

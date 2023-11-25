// https://nyaannyaan.github.io/library/math/sweep-restore.hpp

package main

import (
	"fmt"
	"sort"
)

func main() {
	b := NewLinearBaseRestore([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(b.Restore(13))
}

type pair struct {
	num int
	ids map[int]struct{}
}

type LinearBaseRestore struct {
	Bases []*pair
}

func NewLinearBaseRestore(nums []int) *LinearBaseRestore {
	s := &LinearBaseRestore{}
	for i := 0; i < len(nums); i++ {
		s.Add(nums[i], i)
	}
	return s
}

// 将num加入基底,返回是否添加成功.
func (s *LinearBaseRestore) Add(num int, id int) bool {
	v := &pair{num: num, ids: map[int]struct{}{id: {}}}
	for _, b := range s.Bases {
		if v.num > (v.num ^ b.num) {
			s._apply(v, b)
		}
	}
	if v.num != 0 {
		s.Bases = append(s.Bases, v)
		return true
	}
	return false
}

// 返回:x是否可以由基底线性组合得到,以及组合的基底编号.
func (s *LinearBaseRestore) Restore(x int) (ids []int, ok bool) {
	v := &pair{num: x, ids: map[int]struct{}{}}
	for _, b := range s.Bases {
		if v.num > (v.num ^ b.num) {
			s._apply(v, b)
		}
	}
	if v.num != 0 {
		return
	}
	ids = make([]int, 0, len(v.ids))
	for n := range v.ids {
		ids = append(ids, n)
	}
	sort.Ints(ids)
	ok = true
	return
}

func (s *LinearBaseRestore) _apply(p *pair, o *pair) {
	p.num ^= o.num
	for x := range o.ids {
		s._toggle(p.ids, x)
	}
}

func (s *LinearBaseRestore) _toggle(set map[int]struct{}, x int) {
	if _, ok := set[x]; ok {
		delete(set, x)
	} else {
		set[x] = struct{}{}
	}
}

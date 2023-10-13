package main

import (
	"fmt"
	"math/rand"
)

type LinkedList01 struct {
	m       int32
	prev    []int32
	next    []int32
	visited []bool
}

func NewLinkedList01(n int) *LinkedList01 {
	m := int32(n)
	prev := make([]int32, m+2)
	next := make([]int32, m+2)
	visited := make([]bool, m)
	for i := int32(0); i < m+2; i++ {
		prev[i] = i - 1
		next[i] = i + 1
	}
	return &LinkedList01{m: m, prev: prev, next: next, visited: visited}
}

// 将位置i处的值置为1，返回i的左右两侧第一个未被访问过的位置(不包含i).如果不存在, 返回-1.
// 如果重复 SetOne 则会报错.
func (f *LinkedList01) SetOne(i int) (prev, next int) {
	if f.visited[i] {
		panic("can not add same element twice")
	}
	f.visited[i] = true
	i++
	f.next[f.prev[i]] = f.next[i]
	f.prev[f.next[i]] = f.prev[i]
	prev, next = int(f.prev[i])-1, int(f.next[i])-1
	if prev < 0 {
		prev = -1
	}
	if next >= int(f.m) {
		next = -1
	}
	return
}

func main() {
	list := NewLinkedList01(10)
	fmt.Println(list.SetOne(8))
	fmt.Println(list.SetOne(6))
	fmt.Println(list.SetOne(7))
	fmt.Println(list.SetOne(9))

	n := int(1e4)
	real := NewLinkedList01(n)
	mocker := NewMocker(make([]int, n))
	perm := rand.Perm(n)

	for _, i := range perm {
		prev, next := real.SetOne(i)
		prev2, next2 := mocker.SetOne(i)
		if prev != prev2 || next != next2 {
			panic(fmt.Sprintf("prev: %d, next: %d, prev2: %d, next2: %d", prev, next, prev2, next2))
		}
	}

	fmt.Println("ok")
}

type Mocker struct {
	nums []int
}

func NewMocker(nums []int) *Mocker {
	return &Mocker{nums: append([]int{}, nums...)}
}

func (m *Mocker) SetOne(i int) (prev, next int) {
	if m.nums[i] == 1 {
		panic("can not add same element twice")
	}
	m.nums[i] = 1
	prev, next = -1, -1
	for j := i - 1; j >= 0; j-- {
		if m.nums[j] == 0 {
			prev = j
			break
		}
	}
	for j := i + 1; j < len(m.nums); j++ {
		if m.nums[j] == 0 {
			next = j
			break
		}
	}
	return
}

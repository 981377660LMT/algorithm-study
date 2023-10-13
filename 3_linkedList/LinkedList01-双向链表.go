package main

import "fmt"

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
	fmt.Println(list.SetOne(9))
}

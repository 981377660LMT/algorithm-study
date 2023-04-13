// SkewBinaryList(斜二项堆)
// 纯函数式堆（纯函数式优先级队列）
// https://scrapbox.io/data-structures/Skew_Binary_List
// https://noshi91.github.io/Library/data_structure/persistent_skew_binary_random_access_list.cpp

package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	list := NewPersistentSkewBinaryRandomAccessList()
	fmt.Println(list.Empty())
	list = list.AppendLeft(1)
	fmt.Println(list.Empty())
	fmt.Println(list.Get(0))
	fmt.Println(list.Front())
	list = list.AppendLeft(2)
	fmt.Println(list.Empty())
	fmt.Println(list.Get(1))
	list = list.Set(1, 10)
	fmt.Println(list)
	list = list.PopLeft()
	fmt.Println(list)

	time1 := time.Now()
	for i := 0; i < 1e5; i++ {
		list = list.AppendLeft(i)
		list = list.Set(i, i+1)
	}
	for i := 0; i < 1e5; i++ {
		list = list.PopLeft()
	}
	fmt.Println(time.Since(time1))
}

// API:
//  get(index): 访问下标.
//  front(): 访问第一个元素.
//  set(index, x): 更新下标的元素, 返回新的List.
//  appendLeft(x): 在左边添加元素, 返回新的List.
//  popLeft(): 弹出第一个元素, 返回新的List.
//  empty(): 判断是否为空.

type PH = int

// 可持久化斜二项堆.
type PersistentSkewBinaryRandomAccessList struct {
	root *pDigit
}

func NewPersistentSkewBinaryRandomAccessList() *PersistentSkewBinaryRandomAccessList {
	return &PersistentSkewBinaryRandomAccessList{}
}

func (list *PersistentSkewBinaryRandomAccessList) Empty() bool {
	return list.root == nil
}

func (list *PersistentSkewBinaryRandomAccessList) Front() PH {
	return list.root.tree.value
}

func (list *PersistentSkewBinaryRandomAccessList) Get(i int) PH {
	return list.root.LookUp(i)
}

func (list *PersistentSkewBinaryRandomAccessList) Set(i int, x PH) *PersistentSkewBinaryRandomAccessList {
	if list.root == nil {
		panic("root is nil")
	}
	return &PersistentSkewBinaryRandomAccessList{root: list.root.Update(i, x)}
}

func (list *PersistentSkewBinaryRandomAccessList) AppendLeft(x PH) *PersistentSkewBinaryRandomAccessList {
	if list.root != nil && list.root.next != nil && list.root.size == list.root.next.size {
		return &PersistentSkewBinaryRandomAccessList{
			root: &pDigit{
				size: 1 + list.root.size + list.root.next.size,
				tree: &pTree{
					value: x,
					left:  list.root.tree,
					right: list.root.next.tree,
				},
				next: list.root.next.next,
			},
		}
	}
	return &PersistentSkewBinaryRandomAccessList{
		root: &pDigit{
			size: 1,
			tree: &pTree{
				value: x,
			},
			next: list.root,
		},
	}
}

func (list *PersistentSkewBinaryRandomAccessList) PopLeft() *PersistentSkewBinaryRandomAccessList {
	if list.root == nil {
		panic("root is nil")
	}
	if list.root.size == 1 {
		return &PersistentSkewBinaryRandomAccessList{root: list.root.next}
	}
	chSize := list.root.size >> 1
	return &PersistentSkewBinaryRandomAccessList{
		root: &pDigit{size: chSize, tree: list.root.tree, next: &pDigit{size: chSize, tree: list.root.tree, next: list.root.next}}}
}

func (list *PersistentSkewBinaryRandomAccessList) String() (res string) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	sb := make([]string, 0)
	i := 0
	for {
		sb = append(sb, fmt.Sprintf("%v", list.Get(i)))
		res = "List{" + strings.Join(sb, ", ") + "}"
		i++
	}
}

type pTree struct {
	value PH
	left  *pTree
	right *pTree
}

func (t *pTree) LookUp(size, index int) PH {
	if index == 0 {
		return t.value
	}
	remIndex := index - 1
	chSize := size >> 1
	if remIndex < chSize {
		return t.left.LookUp(chSize, remIndex)
	}
	return t.right.LookUp(chSize, remIndex-chSize)
}

func (t *pTree) Update(size, index int, x PH) *pTree {
	if index == 0 {
		return &pTree{x, t.left, t.right}
	}
	remIndex := index - 1
	chSize := size >> 1
	if remIndex < chSize {
		return &pTree{t.value, t.left.Update(chSize, remIndex, x), t.right}
	}
	return &pTree{t.value, t.left, t.right.Update(chSize, remIndex-chSize, x)}
}

type pDigit struct {
	size int
	tree *pTree
	next *pDigit
}

func (d *pDigit) LookUp(index int) PH {
	if index < d.size {
		return d.tree.LookUp(d.size, index)
	}
	return d.next.LookUp(index - d.size)
}

func (d *pDigit) Update(index int, x PH) *pDigit {
	if index < d.size {
		return &pDigit{d.size, d.tree.Update(d.size, index, x), d.next}
	}
	return &pDigit{d.size, d.tree, d.next.Update(index-d.size, x)}
}

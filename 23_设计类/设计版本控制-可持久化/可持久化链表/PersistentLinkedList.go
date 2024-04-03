// 可持久化链表

package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

func main() {
	list1 := NewPersistentLinkedList()
	list2 := list1.Insert(0, 1)
	list3 := list2.Insert(1, 2)
	list4 := list3.Insert(2, 3)
	list5 := list4.Set(0, 99)
	fmt.Println(list1.GetAll())
	fmt.Println(list2.GetAll())
	fmt.Println(list3.GetAll())
	fmt.Println(list4.GetAll())
	fmt.Println(list4.Slice(2, 3).GetAll())
	fmt.Println(list5.GetAll())
	a, b := list4.Split(1)
	fmt.Println(a.GetAll(), b.GetAll())
}

func init() {
	debug.SetGCPercent(-1)
}

type V = int32
type PersistentLinkedList struct {
	data *PersistentTreap
}

func NewPersistentLinkedList() *PersistentLinkedList {
	return &PersistentLinkedList{data: NIL}
}

func (p *PersistentLinkedList) Get(i int32) V {
	i++
	return GetValueByRank(p.data, i)
}

func (p *PersistentLinkedList) Set(i int32, value V) *PersistentLinkedList {
	i++
	split0, split1 := SplitByRank(p.data, i-1)
	split2, split3 := SplitByRank(split1, 1)
	split2 = split2.Clone()
	split2.value = value
	split1 = Merge(split2, split3)
	return &PersistentLinkedList{data: Merge(split0, split1)}
}

// 将 value 插入到第 i 个元素之前.
func (p *PersistentLinkedList) Insert(i int32, value V) *PersistentLinkedList {
	i++
	split0, split1 := SplitByRank(p.data, i-1)
	newNode := NewPersistentTreap(value)
	newNode.PushUp()
	split0 = Merge(split0, newNode)
	return &PersistentLinkedList{data: Merge(split0, split1)}
}

func (p *PersistentLinkedList) Slice(start, end int32) *PersistentLinkedList {
	start++
	size := p.Len()
	if start == 1 && end == size {
		return p
	}
	if start > end {
		return NewPersistentLinkedList()
	}
	if start == 1 {
		split0, _ := SplitByRank(p.data, end)
		return &PersistentLinkedList{data: split0}
	}
	if end == size {
		_, split1 := SplitByRank(p.data, start-1)
		return &PersistentLinkedList{data: split1}
	}
	split0, _ := SplitByRank(p.data, end)
	_, split1 := SplitByRank(split0, start-1)
	return &PersistentLinkedList{data: split1}
}

// s[:mid], s[mid:]
func (p *PersistentLinkedList) Split(mid int32) (*PersistentLinkedList, *PersistentLinkedList) {
	split0, split1 := SplitByRank(p.data, mid)
	return &PersistentLinkedList{data: split0}, &PersistentLinkedList{data: split1}
}

func (p *PersistentLinkedList) Merge(other *PersistentLinkedList) *PersistentLinkedList {
	return &PersistentLinkedList{data: Merge(p.data, other.data)}
}

func (p *PersistentLinkedList) Len() int32 {
	return p.data.size
}

func (p *PersistentLinkedList) GetAll() []V {
	n := p.Len()
	res := make([]V, n)
	for i := int32(0); i < n; i++ {
		res[i] = p.Get(i)
	}
	return res
}

func createNIL() *PersistentTreap {
	res := &PersistentTreap{}
	res.left, res.right = res, res
	return res
}

var random = NewRandom()
var NIL = createNIL()

type PersistentTreap struct {
	left, right *PersistentTreap
	size        int32
	value       V
}

func NewPersistentTreap(value V) *PersistentTreap {
	return &PersistentTreap{left: NIL, right: NIL, size: 1, value: value}
}

func (p *PersistentTreap) Clone() *PersistentTreap {
	if p == NIL {
		return p
	}
	return &PersistentTreap{left: p.left, right: p.right, size: p.size, value: p.value}
}

func (p *PersistentTreap) PushUp() {
	if p == NIL {
		return
	}
	p.size = p.left.size + p.right.size + 1
}

func SplitByRank(root *PersistentTreap, rank int32) (*PersistentTreap, *PersistentTreap) {
	if root == NIL {
		return NIL, NIL
	}
	root = root.Clone()
	var res0, res1 *PersistentTreap
	if root.left.size >= rank {
		res0, res1 = SplitByRank(root.left, rank)
		root.left = res1
		res1 = root
	} else {
		res0, res1 = SplitByRank(root.right, rank-(root.size-root.right.size))
		root.right = res0
		res0 = root
	}
	root.PushUp()
	return res0, res1
}

func Merge(a, b *PersistentTreap) *PersistentTreap {
	if a == NIL {
		return b
	}
	if b == NIL {
		return a
	}
	if int(random.Rng()*(uint64(a.size)+uint64(b.size))>>32) < int(a.size) {
		a = a.Clone()
		a.right = Merge(a.right, b)
		a.PushUp()
		return a
	} else {
		b = b.Clone()
		b.left = Merge(a, b.left)
		b.PushUp()
		return b
	}
}

func Clone(root *PersistentTreap) *PersistentTreap {
	if root == NIL {
		return NIL
	}
	clone := root.Clone()
	clone.left = Clone(root.left)
	clone.right = Clone(root.right)
	return clone
}

func GetValueByRank(root *PersistentTreap, k int32) V {
	for root.size > 1 {
		if root.left.size >= k {
			root = root.left
		} else {
			k -= root.left.size
			if k == 1 {
				break
			}
			k--
			root = root.right
		}
	}
	return root.value
}

type Random struct {
	seed     uint64
	hashBase uint64
}

func NewRandom() *Random                 { return &Random{seed: uint64(time.Now().UnixNano()/2 + 1)} }
func NewRandomWithSeed(seed int) *Random { return &Random{seed: uint64(seed)} }

func (r *Random) Rng() uint64 {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed & 0xFFFFFFFF
}

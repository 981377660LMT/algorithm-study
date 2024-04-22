// api:
//  1.Insert(index int32, v int8)
//  2.Pop(index int32) int8
//  3.不支持Set操作
//  4.Get(index int32) int8
//  5.Count0(end int32) int32
//  6.Count1(end int32) int32
//  7.Count(end int32, v int8) int32
//  8.Kth0(k int32) int32
//  9.Kth1(k int32) int32
// 10.Kth(k int32, v int8) int32
// 11.Len() int32
// 12.ToList() []int8
// 13.Debug()

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
)

// TODO
// panic: runtime error: negative shift amount
func main() {
	// demo()
	test()
}

func demo() {
	nums := []int8{0, 0}
	wm := NewSplayTreeBitVector(int32(len(nums)), func(i int32) int8 {
		return nums[i]
	})
	wm.Insert(0, 1)
	wm.Pop(2)
	fmt.Println(wm.ToList())
}

const W int8 = 63

type SplayTreeBitVector struct {
	root   int32
	end    int32
	bitLen []int8
	key    []uint64
	size   []int32
	total  []int32
	child  []int32
}

func NewSplayTreeBitVector(n int32, f func(i int32) int8) *SplayTreeBitVector {
	res := &SplayTreeBitVector{
		root:   0,
		end:    1,
		bitLen: []int8{0},
		key:    []uint64{0},
		size:   []int32{0},
		total:  []int32{0},
		child:  []int32{0, 0},
	}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (t *SplayTreeBitVector) Reserve(n int32) {
	n = n/int32(W) + 1
	t.bitLen = append(t.bitLen, make([]int8, n)...)
	t.key = append(t.key, make([]uint64, n)...)
	t.size = append(t.size, make([]int32, n)...)
	t.total = append(t.total, make([]int32, n)...)
	t.child = append(t.child, make([]int32, 2*n)...)
}

func (t *SplayTreeBitVector) Insert(index int32, v int8) {
	if t.root == 0 {
		t.root = t._makeNode(uint64(v), 1)
		return
	}
	n := t.Len()
	if index < 0 {
		index += n
	}
	if index < 0 {
		index = 0
	}
	if index > n {
		index = n
	}
	if index == t.size[t.root] {
		node := t._rightSplay(t.root)
		if t.bitLen[node] == W {
			tmp := t.key[node]<<1 | uint64(v)
			newNode := t._makeNode(tmp&1, 1)
			t.key[node] = tmp >> 1
			t.child[newNode<<1] = node
			t._update(node)
			t.size[newNode] += t.size[node]
			t.total[newNode] += t.total[node]
			t.root = newNode
		} else {
			tmp := t.key[node]
			bl := index - int32(t.bitLen[node]) - t.size[t.child[node<<1]]
			t.key[node] = (((tmp>>bl)<<1 | uint64(v)) << bl) | (tmp & ((1 << bl) - 1))
			t.bitLen[node]++
			t.size[node]++
			t.total[node] += int32(v)
			t.root = node
		}
	} else {
		node := t._kthElmSplay(t.root, index)
		if t.bitLen[node] == W {
			index -= t.size[t.child[node<<1]]
			tmp := t.key[node]
			bl := int32(t.bitLen[node]) - index
			tmp = (((tmp>>bl)<<1 | uint64(v)) << bl) | (tmp & ((1 << bl) - 1))
			newNode := t._makeNode(tmp>>W, 1)
			t.key[node] = tmp & ((1 << W) - 1)
			t._update(node)
			if t.child[node<<1] != 0 {
				t.child[newNode<<1] = t.child[node<<1]
				t.child[node<<1] = 0
				t._update(node)
			}
			t.child[newNode<<1|1] = node
			t._update(newNode)
			t.root = newNode
		} else {
			tmp := t.key[node]
			bl := int32(t.bitLen[node]) - index + t.size[t.child[node<<1]]
			t.key[node] = (((tmp>>bl)<<1 | uint64(v)) << bl) | (tmp & ((1 << bl) - 1))
			t.bitLen[node]++
			t.size[node]++
			t.total[node] += int32(v)
			t.root = node
		}
	}
}

func (t *SplayTreeBitVector) Pop(index int32) int8 {
	n := t.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	root := t._kthElmSplay(t.root, index)
	size, child, key, bitLen, total := t.size, t.child, t.key, t.bitLen, t.total
	b32 := int32(bitLen[root])
	index -= size[child[root<<1]]
	v := key[root]
	res := int8(v >> (b32 - index - 1) & 1)
	if b32 == 1 {
		if child[root<<1] == 0 {
			t.root = child[root<<1|1]
		} else if child[root<<1|1] == 0 {
			t.root = child[root<<1]
		} else {
			node := t._rightSplay(child[root<<1])
			child[node<<1|1] = child[root<<1|1]
			t._update(node)
			t.root = node
		}
	} else {
		key[root] = ((v >> (b32 - index)) << (b32 - index - 1)) | (v & ((1 << (b32 - index - 1)) - 1))
		bitLen[root]--
		size[root]--
		total[root] -= int32(res)
		t.root = root
	}
	return res
}

func (t *SplayTreeBitVector) Get(index int32) int8 {
	n := t.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	t.root = t._kthElmSplay(t.root, index)
	index -= t.size[t.child[t.root<<1]]
	return int8(t.key[t.root] >> (int32(t.bitLen[t.root]) - index - 1) & 1)
}

func (t *SplayTreeBitVector) Count0(end int32) int32 {
	if end < 0 {
		return 0
	}
	if n := t.Len(); end > n {
		end = n
	}
	return end - t._pref(end)
}
func (t *SplayTreeBitVector) Count1(end int32) int32 {
	if end < 0 {
		return 0
	}
	if n := t.Len(); end > n {
		end = n
	}
	return t._pref(end)
}
func (t *SplayTreeBitVector) Count(end int32, key int8) int32 {
	if key == 0 {
		return t.Count0(end)
	}
	return t.Count1(end)
}
func (t *SplayTreeBitVector) Kth0(k int32) int32 {
	n := t.Len()
	if k < 0 || t.Count0(n) <= k {
		return -1
	}
	l, r := int32(0), n
	for r-l > 1 {
		m := (l + r) >> 1
		if m-t._pref(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}
func (t *SplayTreeBitVector) Kth1(k int32) int32 {
	n := t.Len()
	if k < 0 || t.Count1(n) <= k {
		return -1
	}
	l, r := int32(0), n
	for r-l > 1 {
		m := (l + r) >> 1
		if t._pref(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}
func (t *SplayTreeBitVector) Kth(k int32, key int8) int32 {
	if key == 0 {
		return t.Kth0(k)
	}
	return t.Kth1(k)
}
func (t *SplayTreeBitVector) Len() int32 { return t.size[t.root] }
func (t *SplayTreeBitVector) ToList() []int8 {
	if t.root == 0 {
		return nil
	}
	child, key, bitLen := t.child, t.key, t.bitLen
	res := make([]int8, 0, t.Len())
	var rec func(node int32)
	rec = func(node int32) {
		if child[node<<1] != 0 {
			rec(child[node<<1])
		}
		for i := int32(bitLen[node] - 1); i >= 0; i-- {
			res = append(res, int8(key[node]>>i)&1)
		}
		if child[node<<1|1] != 0 {
			rec(child[node<<1|1])
		}
	}
	rec(t.root)
	return res
}
func (t *SplayTreeBitVector) Debug() {
	child, key := t.child, t.key
	var rec func(node int32) int32
	rec = func(node int32) int32 {
		acc := t._popcount(key[node])
		if child[node<<1] != 0 {
			acc += rec(child[node<<1])
		}
		if child[node<<1|1] != 0 {
			acc += rec(child[node<<1|1])
		}
		if acc != t.total[node] {
			fmt.Println(acc, t.total[node])
			panic("acc Error")
		}
		return acc
	}
	rec(t.root)

}
func (t *SplayTreeBitVector) _build(n int32, f func(i int32) int8) {
	end := t.end
	t.Reserve(n)
	index := end
	W32 := int32(W)
	for i := int32(0); i < n; i += W32 {
		j, v := int32(0), uint64(0)
		for j < W32 && i+j < n {
			v <<= 1
			v |= uint64(f(i + j))
			j++
		}
		t.key[index] = v
		t.bitLen[index] = int8(j)
		t.size[index] = j
		t.total[index] = t._popcount(v)
		index++
	}
	t.end = index

	var rec func(l, r int32) int32
	rec = func(l, r int32) int32 {
		mid := (l + r) >> 1
		if l != mid {
			t.child[mid<<1] = rec(l, mid)
			t.size[mid] += t.size[t.child[mid<<1]]
			t.total[mid] += t.total[t.child[mid<<1]]
		}
		if mid+1 != r {
			t.child[mid<<1|1] = rec(mid+1, r)
			t.size[mid] += t.size[t.child[mid<<1|1]]
			t.total[mid] += t.total[t.child[mid<<1|1]]
		}
		return mid
	}
	t.root = rec(end, t.end)
}

func (t *SplayTreeBitVector) _pref(r int32) int32 {
	if r == 0 {
		return 0
	}
	if r == t.Len() {
		return t.total[t.root]
	}
	t.root = t._kthElmSplay(t.root, r-1)
	r -= t.size[t.child[t.root<<1]]
	return t.total[t.root] - t._popcount(t.key[t.root]&((1<<(int32(t.bitLen[t.root])-r))-1)) - t.total[t.child[t.root<<1|1]]
}

func (t *SplayTreeBitVector) _makeNode(key uint64, bitLen int8) int32 {
	end := t.end
	if end >= int32(len(t.key)) {
		t.key = append(t.key, key)
		t.bitLen = append(t.bitLen, bitLen)
		t.size = append(t.size, int32(bitLen))
		t.total = append(t.total, t._popcount(key))
		t.child = append(t.child, 0, 0)
	} else {
		t.key[end] = key
		t.bitLen[end] = bitLen
		t.size[end] = int32(bitLen)
		t.total[end] = t._popcount(key)
	}
	t.end++
	return end
}

func (t *SplayTreeBitVector) _updateTriple(x, y, z int32) {
	child, bitLen, size, total := t.child, t.bitLen, t.size, t.total
	lx, rx := child[x<<1], child[x<<1|1]
	ly, ry := child[y<<1], child[y<<1|1]
	size[z] = size[x]
	size[x] = int32(bitLen[x]) + size[lx] + size[rx]
	size[y] = int32(bitLen[y]) + size[ly] + size[ry]
	total[z] = total[x]
	total[x] = total[lx] + t._popcount(t.key[x]) + total[rx]
	total[y] = total[ly] + t._popcount(t.key[y]) + total[ry]
}

func (t *SplayTreeBitVector) _updateDouble(x, y int32) {
	child, bitLen, size, total := t.child, t.bitLen, t.size, t.total
	lx, rx := child[x<<1], child[x<<1|1]
	size[y] = size[x]
	size[x] = int32(bitLen[x]) + size[lx] + size[rx]
	total[y] = total[x]
	total[x] = total[lx] + t._popcount(t.key[x]) + total[rx]
}

func (t *SplayTreeBitVector) _update(node int32) {
	lnode, rnode := t.child[node<<1], t.child[node<<1|1]
	t.size[node] = int32(t.bitLen[node]) + t.size[lnode] + t.size[rnode]
	t.total[node] = t._popcount(t.key[node]) + t.total[lnode] + t.total[rnode]
}

// 这里的path可以不传引用.
// TODO: d 是否需要是uint64
func (t *SplayTreeBitVector) _splay(path []int32, d int32) {
	child := t.child
	g := d & 1
	for len(path) > 1 {
		pnode := path[len(path)-1]
		path = path[:len(path)-1]
		gnode := path[len(path)-1]
		path = path[:len(path)-1]
		f := d >> 1 & 1
		node := child[pnode<<1|g^1] // TODO
		var nnode int32
		if g == f {
			nnode = pnode<<1 | f
		} else {
			nnode = node<<1 | f
		}
		child[pnode<<1|g^1] = child[node<<1|g]
		child[node<<1|g] = pnode
		child[gnode<<1|f^1] = child[nnode]
		child[nnode] = gnode
		t._updateTriple(gnode, pnode, node)
		if len(path) == 0 {
			return
		}
		d >>= 2
		g = d & 1
		child[path[len(path)-1]<<1|g^1] = node
	}
	pnode := path[len(path)-1]
	path = path[:len(path)-1]
	node := child[pnode<<1|g^1]
	child[pnode<<1|g^1] = child[node<<1|g]
	child[node<<1|g] = pnode
	t._updateDouble(pnode, node)
}

func (t *SplayTreeBitVector) _kthElmSplay(node, k int32) int32 {
	child, bitLen, size := t.child, t.bitLen, t.size
	d := int32(0)
	path := []int32{}
	for {
		b32 := int32(bitLen[node])
		tmp := size[child[node<<1]] + b32
		if tmp-b32 <= k && k < tmp {
			if len(path) > 0 {
				t._splay(path, d)
			}
			return node
		}
		if tmp > k {
			d = d<<1 | 1
		} else {
			d = d << 1
		}
		path = append(path, node)
		if tmp <= k {
			node = child[node<<1|1]
			k -= tmp
		} else {
			node = child[node<<1]
		}
	}
}

func (t *SplayTreeBitVector) _leftSplay(node int32) int32 {
	if node == 0 {
		return 0
	}
	child := t.child
	if child[node<<1] == 0 {
		return node
	}
	path := []int32{}
	for child[node<<1] != 0 {
		path = append(path, node)
		node = child[node<<1]
	}
	t._splay(path, 1<<len(path)-1)
	return node
}

func (t *SplayTreeBitVector) _rightSplay(node int32) int32 {
	if node == 0 {
		return 0
	}
	child := t.child
	if child[node<<1|1] == 0 {
		return node
	}
	path := []int32{}
	for child[node<<1|1] != 0 {
		path = append(path, node)
		node = child[node<<1|1]
	}
	t._splay(path, 0)
	return node
}

func (t *SplayTreeBitVector) _popcount(v uint64) int32 {
	return int32(bits.OnesCount64(v))
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func max64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func test() {
	for i := 0; i < 2; i++ {
		n := rand.Intn(1e6) + 50
		nums := make([]int8, n)
		for i := 0; i < n; i++ {
			nums[i] = int8(rand.Intn(2))
		}
		bv := NewSplayTreeBitVector(int32(n), func(i int32) int8 { return nums[i] })

		count := func(end int32, v int8) int32 {
			res := int32(0)
			for i := int32(0); i < end; i++ {
				if nums[i] == v {
					res++
				}
			}
			return res
		}

		kth := func(k int32, v int8) int32 {
			res := int32(0)
			for i := int32(0); i < int32(len(nums)); i++ {
				if nums[i] == v {
					if res == k {
						return i
					}
					res++
				}
			}
			return -1
		}

		insert := func(index int32, v int8) {
			nums = append(nums, 0)
			copy(nums[index+1:], nums[index:])
			nums[index] = v
		}
		_ = insert

		pop := func(index int32) int8 {
			res := nums[index]
			nums = append(nums[:index], nums[index+1:]...)
			return res
		}
		_ = pop

		for j := 0; j < 1000; j++ {
			// count
			countIndex := int32(rand.Intn(n + 1))
			if bv.Count0(countIndex) != count(countIndex, 0) {
				panic("error1")
			}
			if bv.Count1(countIndex) != count(countIndex, 1) {
				panic("error2")
			}

			// kth
			kthIndex := int32(rand.Intn(n + 1))
			if bv.Kth0(kthIndex) != kth(kthIndex, 0) {
				panic("error3")
			}
			if bv.Kth1(kthIndex) != kth(kthIndex, 1) {
				panic("error4")
			}

			// insert
			insertIndex := rand.Intn(n + 1)
			insertValue := int8(rand.Intn(2))
			insert(int32(insertIndex), insertValue)
			bv.Insert(int32(insertIndex), insertValue)

			// fmt.Println(wm.ToList(), nums, "after insert", insertIndex, insertValue)

			// pop
			popIndex := rand.Intn(len(nums))
			if bv.Pop(int32(popIndex)) != pop(int32(popIndex)) {
				panic("error")
			}

			// fmt.Println(wm.ToList(), nums, "after pop", popIndex)

			// len
			if bv.Len() != int32(len(nums)) {
				panic("error")
			}

			// get
			// for i := 0; i < len(nums); i++ {
			// 	if bv.Get(int32(i)) != nums[i] {
			// 		fmt.Println(bv.ToList(), nums, i, n, len(nums))
			// 		panic("error get")
			// 	}
			// }

			// toList
			list := bv.ToList()
			for i := 0; i < len(nums); i++ {
				if list[i] != nums[i] {
					fmt.Println(list, nums, i, list[i], nums[i])
					panic("error toList")
				}
			}

			bv.Debug()
		}
	}

	fmt.Println("ok")
}

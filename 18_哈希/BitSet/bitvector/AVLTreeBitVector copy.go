// api:
//  1.Insert(index int32, v int8)
//  2.Pop(index int32) int8
//  3.Set(index int32, v int8)
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

func main() {
	nums := []int8{0, 0}
	wm := NewAVLTreeBitVector(int32(len(nums)), func(i int32) int8 {
		return nums[i]
	})
	wm.Insert(0, 1)
	wm.Pop(2)
	fmt.Println(wm.ToList())
	wm.Set(0, 1)
	fmt.Println(wm.ToList())
	test()
}

type AVLTreeBitVector struct {
	root    int32
	end     int32 // 使用的结点数
	bitLen  []int32
	key     []uint64 // 结点mask
	total   []int32  // 子树onesCount之和
	size    []int32  // 子树大小
	left    []int32
	right   []int32
	balance []int8 // 左子树高度-右子树高度
}

const W int32 = 63

func NewAVLTreeBitVector(n int32, f func(i int32) int8) *AVLTreeBitVector {
	res := &AVLTreeBitVector{
		root:    0,
		end:     1,
		bitLen:  []int32{0},
		key:     []uint64{0},
		total:   []int32{0},
		size:    []int32{0},
		left:    []int32{0},
		right:   []int32{0},
		balance: []int8{0},
	}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (t *AVLTreeBitVector) Reserve(n int32) {
	n = n/W + 1
	t.bitLen = append(t.bitLen, make([]int32, n)...)
	t.key = append(t.key, make([]uint64, n)...)
	t.size = append(t.size, make([]int32, n)...)
	t.total = append(t.total, make([]int32, n)...)
	t.left = append(t.left, make([]int32, n)...)
	t.right = append(t.right, make([]int32, n)...)
	t.balance = append(t.balance, make([]int8, n)...)
}

func (t *AVLTreeBitVector) Insert(index int32, v int8) {
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

	if t.root == 0 {
		t.root = t._makeNode(uint64(v), 1)
		return
	}

	v32 := int32(v)

	node := t.root
	path := []int32{}
	d := int32(0)
	for node != 0 {
		tmp := t.size[t.left[node]] + t.bitLen[node]
		if tmp-t.bitLen[node] <= index && index <= tmp {
			break
		}
		d <<= 1
		t.size[node]++
		t.total[node] += v32
		path = append(path, node)
		if tmp > index {
			node = t.left[node]
			d |= 1
		} else {
			node = t.right[node]
			index -= tmp
		}
	}
	index -= t.size[t.left[node]]
	if t.bitLen[node] < W {
		k := t.key[node]
		bl := index
		t.key[node] = (((k>>bl)<<1 | uint64(v)) << bl) | (k & ((1 << bl) - 1))
		t.bitLen[node]++
		t.size[node]++
		t.total[node] += v32
		return
	}
	path = append(path, node)
	t.size[node]++
	t.total[node] += v32
	k := t.key[node]
	bl := uint64(W - index)
	k = (((k>>bl)<<1 | uint64(v)) << bl) | (k & ((1 << bl) - 1))
	leftKey := k >> W
	leftKeyPopcount := int32(leftKey & 1)
	t.key[node] = k & ((1 << W) - 1)
	node = t.left[node]
	d <<= 1
	d |= 1

	if node == 0 {
		last := path[len(path)-1]
		if t.bitLen[last] < W {
			t.bitLen[last]++
			t.key[last] = (t.key[last] << 1) | leftKey
			return
		} else {
			t.left[last] = t._makeNode(leftKey, 1)
		}
	} else {
		path = append(path, node)
		t.size[node]++
		t.total[node] += leftKeyPopcount
		d <<= 1
		for t.right[node] != 0 {
			node = t.right[node]
			path = append(path, node)
			t.size[node]++
			t.total[node] += leftKeyPopcount
			d <<= 1
		}
		if t.bitLen[node] < W {
			t.bitLen[node]++
			t.key[node] = (t.key[node] << 1) | leftKey
			return
		} else {
			t.right[node] = t._makeNode(leftKey, 1)
		}
	}
	newNode := int32(0)
	for len(path) > 0 {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if d&1 == 1 {
			t.balance[node]++
		} else {
			t.balance[node]--
		}
		d >>= 1
		if t.balance[node] == 0 {
			break
		}
		if t.balance[node] == 2 {
			if t.balance[t.left[node]] == -1 {
				newNode = t._rotateLR(node)
			} else {
				newNode = t._rotateL(node)
			}
			break
		} else if t.balance[node] == -2 {
			if t.balance[t.right[node]] == 1 {
				newNode = t._rotateRL(node)
			} else {
				newNode = t._rotateR(node)
			}
			break
		}
	}
	if newNode != 0 {
		if len(path) > 0 {
			if d&1 == 1 {
				t.left[path[len(path)-1]] = newNode
			} else {
				t.right[path[len(path)-1]] = newNode
			}
		} else {
			t.root = newNode
		}
	}
}

func (t *AVLTreeBitVector) Pop(index int32) int8 {
	if index < 0 {
		index += t.Len()
	}
	left, right, size := t.left, t.right, t.size
	bitLen, keys, total := t.bitLen, t.key, t.total
	node := t.root
	d := int32(0)
	path := []int32{}
	for node != 0 {
		t := size[left[node]] + bitLen[node]
		if t-bitLen[node] <= index && index < t {
			break
		}
		path = append(path, node)
		if t > index {
			node = left[node]
			d <<= 1
			d |= 1
		} else {
			node = right[node]
			index -= t
		}
	}
	index -= size[left[node]]
	v := keys[node]
	res := int32(v >> (bitLen[node] - index - 1) & 1)
	if bitLen[node] == 1 {
		t._popUnder(path, d, node, res)
		return int8(res)
	}
	keys[node] = ((v >> (bitLen[node] - index)) << (bitLen[node] - index - 1)) | (v & ((1 << (bitLen[node] - index - 1)) - 1))
	bitLen[node]--
	size[node]--
	total[node] -= res
	for _, p := range path {
		size[p]--
		total[p] -= res
	}
	return int8(res)
}

func (t *AVLTreeBitVector) Set(index int32, v int8) {
	n := t.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("index out of range: %d", index))
	}

	left, right, bitLen, size, key, total := t.left, t.right, t.bitLen, t.size, t.key, t.total
	node := t.root
	path := []int32{}
	for true {
		tmp := size[left[node]] + bitLen[node]
		path = append(path, node)
		if tmp-bitLen[node] <= index && index < tmp {
			index -= size[left[node]]
			if v == 1 {
				key[node] |= 1 << index
			} else {
				key[node] &= ^(1 << index)
			}
			break
		} else if tmp > index {
			node = left[node]
		} else {
			node = right[node]
			index -= tmp
		}
	}
	for len(path) > 0 {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		total[node] = t._popcount(key[node]) + total[left[node]] + total[right[node]]
	}

}

func (t *AVLTreeBitVector) Get(index int32) int8 {
	if index < 0 {
		index += t.Len()
	}
	left, right, bitLen, size, key := t.left, t.right, t.bitLen, t.size, t.key
	node := t.root
	for true {
		tmp := size[left[node]] + bitLen[node]
		if tmp-bitLen[node] <= index && index < tmp {
			index -= size[left[node]]
			return int8(key[node] >> index & 1)
		}
		if tmp > index {
			node = left[node]
		} else {
			node = right[node]
			index -= tmp
		}
	}
	panic("unreachable")
}

func (t *AVLTreeBitVector) Count0(end int32) int32 {
	if end < 0 {
		return 0
	}
	if n := t.Len(); end > n {
		end = n
	}
	return end - t._pref(end)
}

func (t *AVLTreeBitVector) Count1(end int32) int32 {
	if end < 0 {
		return 0
	}
	if n := t.Len(); end > n {
		end = n
	}
	return t._pref(end)
}
func (t *AVLTreeBitVector) Count(end int32, v int8) int32 {
	if v == 1 {
		return t.Count1(end)
	}
	return t.Count0(end)
}
func (t *AVLTreeBitVector) Kth0(k int32) int32 {
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
func (t *AVLTreeBitVector) Kth1(k int32) int32 {
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
func (t *AVLTreeBitVector) Kth(k int32, v bool) int32 {
	if v {
		return t.Kth1(k)
	}
	return t.Kth0(k)
}
func (t *AVLTreeBitVector) Len() int32 { return t.size[t.root] }

func (t *AVLTreeBitVector) ToList() []int8 {
	if t.root == 0 {
		return nil
	}
	left, right, key, bitLen := t.left, t.right, t.key, t.bitLen
	res := make([]int8, 0, t.Len())
	var rec func(node int32)
	rec = func(node int32) {
		if left[node] != 0 {
			rec(left[node])
		}
		for i := int32(0); i < bitLen[node]; i++ {
			res = append(res, int8(key[node]>>i&1))
		}
		if right[node] != 0 {
			rec(right[node])
		}
	}
	rec(t.root)
	return res
}

func (t *AVLTreeBitVector) Debug() {
	left, right, key := t.left, t.right, t.key
	var rec func(node int32) int32
	rec = func(node int32) int32 {
		acc := t._popcount(key[node])
		if left[node] != 0 {
			acc += rec(left[node])
		}
		if right[node] != 0 {
			acc += rec(right[node])
		}
		if acc != t.total[node] {
			fmt.Println(node, acc, t.total[node])
			panic("error")
		}
		return acc
	}
	rec(t.root)
}

func (t *AVLTreeBitVector) _build(n int32, f func(i int32) int8) {
	bit := bits.Len32(uint32(n)) + 2
	mask := int32(1<<bit - 1)
	end := t.end
	t.Reserve(n)
	index := end
	for i := int32(0); i < n; i += W {
		j, v := int32(0), uint64(0)
		for j < W && i+j < n {
			v <<= 1
			v |= uint64(f(i + j))
			j++
		}
		t.key[index] = v
		t.bitLen[index] = j
		t.size[index] = j
		t.total[index] = t._popcount(v)
		index++
	}
	t.end = index

	var rec func(lr int32) int32
	rec = func(lr int32) int32 {
		l, r := lr>>bit, lr&mask
		mid := (l + r) >> 1
		hl, hr := int32(0), int32(0)
		if l != mid {
			le := rec(l<<bit | mid)
			t.left[mid], hl = le>>bit, le&mask
			t.size[mid] += t.size[t.left[mid]]
			t.total[mid] += t.total[t.left[mid]]
		}
		if mid+1 != r {
			ri := rec((mid+1)<<bit | r)
			t.right[mid], hr = ri>>bit, ri&mask
			t.size[mid] += t.size[t.right[mid]]
			t.total[mid] += t.total[t.right[mid]]
		}
		t.balance[mid] = int8(hl - hr)
		return mid<<bit | (max32(hl, hr) + 1)
	}
	t.root = rec(end<<bit|t.end) >> bit
}

func (t *AVLTreeBitVector) _rotateL(node int32) int32 {
	left, right, size, balance, total := t.left, t.right, t.size, t.balance, t.total
	u := left[node]
	size[u] = size[node]
	total[u] = total[node]
	size[node] -= size[left[u]] + t.bitLen[u]
	total[node] -= total[left[u]] + t._popcount(t.key[u])
	left[node] = right[u]
	right[u] = node
	if balance[u] == 1 {
		balance[u] = 0
		balance[node] = 0
	} else {
		balance[u] = -1
		balance[node] = 1
	}
	return u
}

func (t *AVLTreeBitVector) _rotateR(node int32) int32 {
	left, right, size, balance, total := t.left, t.right, t.size, t.balance, t.total
	u := right[node]
	size[u] = size[node]
	total[u] = total[node]
	size[node] -= size[right[u]] + t.bitLen[u]
	total[node] -= total[right[u]] + t._popcount(t.key[u])
	right[node] = left[u]
	left[u] = node
	if balance[u] == -1 {
		balance[u] = 0
		balance[node] = 0
	} else {
		balance[u] = 1
		balance[node] = -1
	}
	return u
}

func (t *AVLTreeBitVector) _rotateLR(node int32) int32 {
	left, right, size, total := t.left, t.right, t.size, t.total
	B := left[node]
	E := right[B]
	size[E] = size[node]
	size[node] -= size[B] - size[right[E]]
	size[B] -= size[right[E]] + t.bitLen[E]
	total[E] = total[node]
	total[node] -= total[B] - total[right[E]]
	total[B] -= total[right[E]] + t._popcount(t.key[E])
	right[B] = left[E]
	left[E] = B
	left[node] = right[E]
	right[E] = node
	t._updateBalance(E)
	return E
}

func (t *AVLTreeBitVector) _rotateRL(node int32) int32 {
	left, right, size, total := t.left, t.right, t.size, t.total
	C := right[node]
	D := left[C]
	size[D] = size[node]
	size[node] -= size[C] - size[left[D]]
	size[C] -= size[left[D]] + t.bitLen[D]
	total[D] = total[node]
	total[node] -= total[C] - total[left[D]]
	total[C] -= total[left[D]] + t._popcount(t.key[D])
	left[C] = right[D]
	right[D] = C
	right[node] = left[D]
	left[D] = node
	t._updateBalance(D)
	return D
}

func (t *AVLTreeBitVector) _updateBalance(node int32) {
	balance := t.balance
	if b := balance[node]; b == 1 {
		balance[t.right[node]] = -1
		balance[t.left[node]] = 0
	} else if b == -1 {
		balance[t.right[node]] = 0
		balance[t.left[node]] = 1
	} else {
		balance[t.right[node]] = 0
		balance[t.left[node]] = 0
	}
	balance[node] = 0
}

func (t *AVLTreeBitVector) _pref(r int32) int32 {
	left, right, bitLen, size, key, total := t.left, t.right, t.bitLen, t.size, t.key, t.total
	node := t.root
	s := int32(0)
	for r > 0 {
		tmp := size[left[node]] + bitLen[node]
		if tmp-bitLen[node] < r && r <= tmp {
			r -= size[left[node]]
			s += total[left[node]] + t._popcount(key[node]>>(bitLen[node]-r))
			break
		}
		if tmp > r {
			node = left[node]
		} else {
			s += total[left[node]] + t._popcount(key[node])
			node = right[node]
			r -= tmp
		}
	}
	return s
}

func (t *AVLTreeBitVector) _makeNode(v uint64, bitLen int32) int32 {
	end := t.end
	if end >= int32(len(t.key)) {
		t.key = append(t.key, v)
		t.bitLen = append(t.bitLen, bitLen)
		t.size = append(t.size, bitLen)
		t.total = append(t.total, t._popcount(v))
		t.left = append(t.left, 0)
		t.right = append(t.right, 0)
		t.balance = append(t.balance, 0)
	} else {
		t.key[end] = v
		t.bitLen[end] = bitLen
		t.size[end] = bitLen
		t.total[end] = t._popcount(v)
	}
	t.end++
	return end
}

// 这里的path可以不用*[]int32
func (t *AVLTreeBitVector) _popUnder(path []int32, d int32, node int32, res int32) {
	left, right, size, bitLen, balance, keys, total := t.left, t.right, t.size, t.bitLen, t.balance, t.key, t.total
	fd, lmaxTotal, lmaxBitLen := int32(0), int32(0), int32(0)

	if left[node] != 0 && right[node] != 0 {
		path = append(path, node)
		d <<= 1
		d |= 1
		lmax := left[node]
		for right[lmax] != 0 {
			path = append(path, lmax)
			d <<= 1
			fd <<= 1
			fd |= 1
			lmax = right[lmax]
		}
		lmaxTotal = t._popcount(keys[lmax])
		lmaxBitLen = bitLen[lmax]
		keys[node] = keys[lmax]
		bitLen[node] = lmaxBitLen
		node = lmax
	}
	var cNode int32
	if left[node] == 0 {
		cNode = right[node]
	} else {
		cNode = left[node]
	}
	if len(path) > 0 {
		if d&1 == 1 {
			left[path[len(path)-1]] = cNode
		} else {
			right[path[len(path)-1]] = cNode
		}
	} else {
		t.root = cNode
		return
	}

	for len(path) > 0 {
		newNode := int32(0)
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if d&1 == 1 {
			balance[node]--
		} else {
			balance[node]++
		}
		if fd&1 == 1 {
			size[node] -= lmaxBitLen
			total[node] -= lmaxTotal
		} else {
			size[node]--
			total[node] -= res
		}
		d >>= 1
		fd >>= 1
		if balance[node] == 2 {
			if balance[left[node]] < 0 {
				newNode = t._rotateLR(node)
			} else {
				newNode = t._rotateL(node)
			}
		} else if balance[node] == -2 {
			if balance[right[node]] > 0 {
				newNode = t._rotateRL(node)
			} else {
				newNode = t._rotateR(node)
			}
		} else if balance[node] != 0 {
			break
		}
		if newNode != 0 {
			if len(path) == 0 {
				t.root = newNode
				return
			}
			if d&1 == 1 {
				left[path[len(path)-1]] = newNode
			} else {
				right[path[len(path)-1]] = newNode
			}
			if balance[newNode] != 0 {
				break
			}
		}
	}

	for len(path) > 0 {
		node := path[len(path)-1]
		path = path[:len(path)-1]
		if fd&1 == 1 {
			size[node] -= lmaxBitLen
			total[node] -= lmaxTotal
		} else {
			size[node]--
			total[node] -= res
		}
		fd >>= 1
	}
}

func (t *AVLTreeBitVector) _popcount(v uint64) int32 {
	return int32(bits.OnesCount64(v))
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func test() {
	for i := 0; i < 5000; i++ {
		n := rand.Intn(70) + 1
		nums := make([]int8, n)
		for i := 0; i < n; i++ {
			nums[i] = int8(rand.Intn(2))
		}
		wm := NewAVLTreeBitVector(int32(n), func(i int32) int8 { return nums[i] })

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

		// count
		for i := 0; i < len(nums); i++ {
			if wm.Count0(int32(i)) != count(int32(i), 0) {
				panic("error")
			}
			if wm.Count1(int32(i)) != count(int32(i), 1) {
				panic("error")
			}
		}

		// kth
		for i := 0; i < len(nums); i++ {
			if wm.Kth0(int32(i)) != kth(int32(i), 0) {
				panic("error")
			}
			if wm.Kth1(int32(i)) != kth(int32(i), 1) {
				panic("error")
			}
		}

		// insert
		// insertIndex := rand.Intn(n + 1)
		// insertValue := int8(rand.Intn(2))
		// insert(int32(insertIndex), insertValue)
		// wm.Insert(int32(insertIndex), insertValue)

		// fmt.Println(wm.ToList(), nums, "after insert", insertIndex, insertValue)

		// pop
		// popIndex := rand.Intn(len(nums))
		// if wm.Pop(int32(popIndex)) != pop(int32(popIndex)) {
		// 	panic("error")
		// }

		// fmt.Println(wm.ToList(), nums, "after pop", popIndex)

		// set
		// setIndex := rand.Intn(len(nums))
		// setValue := int8(rand.Intn(2))
		// nums[setIndex] = setValue
		// wm.Set(int32(setIndex), setValue)
		// fmt.Println(wm.ToList(), nums, "after set", setIndex, setValue)

		// len
		if wm.Len() != int32(len(nums)) {
			panic("error")
		}

		// get
		for i := 0; i < len(nums); i++ {
			if wm.Get(int32(i)) != nums[i] {
				fmt.Println(wm.ToList(), nums, i, len(nums))
				panic("error get")
			}
		}

		// toList
		list := wm.ToList()
		for i := 0; i < len(nums); i++ {
			if list[i] != nums[i] {
				fmt.Println(list, nums, i, list[i], nums[i])
				panic("error toList")
			}
		}

		wm.Debug()
	}

	fmt.Println("ok")
}

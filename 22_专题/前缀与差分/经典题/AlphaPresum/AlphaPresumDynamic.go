// Api:
//  1. IsAscending(start, end int) bool
//  2. IsDescending(start, end int) bool
//  3. Min(start, end int) int32
//  4. Max(start, end int) int32
//  5. Set(index int32, value int32)
//  6. Get(index int32) int32
//  7. Count(start, end int, c int32) int32
//  8. CountAll(start, end int) []int32

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"
)

func main() {
	test()
}

const SIGMA int32 = 26
const OFFSET int32 = 97

type AlphaPresumDynamic struct {
	n    int32
	arr  []int32
	tree [SIGMA]*bITArray01

	revTree      [SIGMA]*bITArray01
	revTreeBuilt bool
}

func NewAlphaPresumDynamic(n int32, f func(i int32) int32) *AlphaPresumDynamic {
	arr := make([]int32, n)
	for i := int32(0); i < n; i++ {
		arr[i] = f(i) - OFFSET
	}
	tree := [SIGMA]*bITArray01{}
	for i := int32(0); i < SIGMA; i++ {
		tree[i] = newBITArray01(n)
	}
	for i := int32(0); i < n; i++ {
		tree[arr[i]].Add(i)
	}
	res := &AlphaPresumDynamic{n: n, arr: arr, tree: tree}
	return res
}

// 如果一个范围内的元素是单调递增的，那么对于这个范围内的任何字符，
// 它们出现的位置应该是连续的，且随着字符值的增加，出现的位置也应该是递增的.
func (a *AlphaPresumDynamic) IsAscending(start, end int32) bool {
	if start < 0 {
		start = 0
	}
	if end > a.n {
		end = a.n
	}
	if start >= end {
		return true
	}

	ptr := start
	for i := int32(0); i < SIGMA; i++ {
		c := a.tree[i].QueryRange(start, end)
		ptr += c
		if ptr > end || a.tree[i].QueryRange(start, ptr) != c {
			return false
		}
	}
	return true
}

func (a *AlphaPresumDynamic) IsDescending(start, end int32) bool {
	if start < 0 {
		start = 0
	}
	if end > a.n {
		end = a.n
	}
	if start >= end {
		return true
	}

	a.buildRevTree()
	ptr := start
	for i := int32(0); i < SIGMA; i++ {
		c := a.revTree[i].QueryRange(start, end)
		ptr += c
		if ptr > end || a.revTree[i].QueryRange(start, ptr) != c {
			return false
		}
	}
	return true
}

func (a *AlphaPresumDynamic) Min(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > a.n {
		end = a.n
	}
	if start >= end {
		return -1
	}
	for i := int32(0); i < SIGMA; i++ {
		if a.tree[i].QueryRange(start, end) > 0 {
			return i + OFFSET
		}
	}
	return -1
}

func (a *AlphaPresumDynamic) Max(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > a.n {
		end = a.n
	}
	if start >= end {
		return -1
	}
	for i := SIGMA - 1; i >= 0; i-- {
		if a.tree[i].QueryRange(start, end) > 0 {
			return i + OFFSET
		}
	}
	return -1
}

func (a *AlphaPresumDynamic) Set(index, value int32) {
	oldValue := a.arr[index]
	if oldValue == value {
		return
	}
	value -= OFFSET
	a.tree[oldValue].Remove(index)
	a.tree[value].Add(index)
	a.arr[index] = value
	if a.revTreeBuilt {
		a.revTree[SIGMA-1-oldValue].Remove(index)
		a.revTree[SIGMA-1-value].Add(index)
	}
}

func (a *AlphaPresumDynamic) Get(index int32) int32 {
	return a.arr[index] + OFFSET
}

func (a *AlphaPresumDynamic) Count(start, end, c int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > a.n {
		end = a.n
	}
	if start >= end {
		return 0
	}
	return a.tree[c-OFFSET].QueryRange(start, end)
}

func (a *AlphaPresumDynamic) CountAll(start, end int32) [SIGMA]int32 {
	if start < 0 {
		start = 0
	}
	if end > a.n {
		end = a.n
	}
	if start >= end {
		return [SIGMA]int32{}
	}
	res := [SIGMA]int32{}
	for i := int32(0); i < SIGMA; i++ {
		res[i] = a.tree[i].QueryRange(start, end)
	}
	return res
}

func (a *AlphaPresumDynamic) buildRevTree() {
	if a.revTreeBuilt {
		return
	}
	for i := int32(0); i < SIGMA; i++ {
		a.revTree[i] = newBITArray01(a.n)
	}
	for i := int32(0); i < a.n; i++ {
		a.revTree[SIGMA-1-a.arr[i]].Add(i)
	}
	a.revTreeBuilt = true
}

// !Point Add Range Sum, 0-based.
type bITArray32 struct {
	n     int32
	total int32
	data  []int32
}

func newBitArray(n int32) *bITArray32 {
	res := &bITArray32{n: n, data: make([]int32, n)}
	return res
}

func newBitArrayFrom(n int32, f func(i int32) int32) *bITArray32 {
	total := int32(0)
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &bITArray32{n: n, total: total, data: data}
}

func (b *bITArray32) Add(index int32, v int32) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *bITArray32) QueryPrefix(end int32) int32 {
	if end > b.n {
		end = b.n
	}
	res := int32(0)
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *bITArray32) QueryRange(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := int32(0), int32(0)
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *bITArray32) QueryAll() int32 {
	return b.total
}

// 01树状数组.
type bITArray01 struct {
	n    int32
	size int32 // data、bit的长度
	data []uint64
	bit  *bITArray32
}

func newBITArray01(n int32) *bITArray01 {
	return newBITArray01From(n, func(index int32) bool { return false })
}

func newBITArray01From(n int32, f func(index int32) bool) *bITArray01 {
	size := n>>6 + 1
	data := make([]uint64, size)
	for i := int32(0); i < n; i++ {
		if f(i) {
			data[i>>6] |= 1 << (i & 63)
		}
	}
	bit := newBitArrayFrom(size, func(i int32) int32 { return int32(bits.OnesCount64(data[i])) })
	return &bITArray01{n: n, size: size, data: data, bit: bit}
}

func (bit01 *bITArray01) QueryAll() int32 {
	return bit01.bit.QueryAll()
}

func (bit01 *bITArray01) QueryPrefix(end int32) int32 {
	i, j := end>>6, end&63
	res := bit01.bit.QueryPrefix(i)
	res += int32(bits.OnesCount64(bit01.data[i] & ((1 << j) - 1)))
	return res
}

func (bit01 *bITArray01) QueryRange(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > bit01.n {
		end = bit01.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return bit01.QueryPrefix(end)
	}
	res := int32(0)
	res -= int32(bits.OnesCount64(bit01.data[start>>6] & ((1 << (start & 63)) - 1)))
	res += int32(bits.OnesCount64(bit01.data[end>>6] & ((1 << (end & 63)) - 1)))
	res += bit01.bit.QueryRange(start>>6, end>>6)
	return res
}

func (bit01 *bITArray01) Add(index int32) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 1 {
		return false
	}
	bit01.data[i] |= 1 << j
	bit01.bit.Add(i, 1)
	return true
}

func (bit01 *bITArray01) Remove(index int32) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 0 {
		return false
	}
	bit01.data[i] ^= 1 << j
	bit01.bit.Add(i, -1)
	return true
}

func (bit01 *bITArray01) Has(index int32) bool {
	i, j := index>>6, index&63
	return (bit01.data[i]>>j)&1 == 1
}

func (bit01 *bITArray01) String() string {
	res := []string{}
	for i := int32(0); i < bit01.n; i++ {
		if bit01.QueryRange(i, i+1) == 1 {
			res = append(res, fmt.Sprintf("%d", i))
		}
	}
	return fmt.Sprintf("BITArray01: [%v]", strings.Join(res, ", "))
}

func test() {
	for i := 0; i < 10; i++ {
		n := int32(rand.Intn(500) + 1)
		arr := make([]int32, n)
		for i := int32(0); i < n; i++ {
			arr[i] = int32(rand.Intn(26) + 97)
		}
		S := NewAlphaPresumDynamic(n, func(i int32) int32 { return arr[i] })

		isAscendingBruteForce := func(start, end int32) bool {
			for i := start; i < end-1; i++ {
				if arr[i] > arr[i+1] {
					return false
				}
			}
			return true
		}

		isDescendingBruteForce := func(start, end int32) bool {
			for i := start; i < end-1; i++ {
				if arr[i] < arr[i+1] {
					return false
				}
			}
			return true
		}
		_ = isDescendingBruteForce

		minBruteForce := func(start, end int32) int32 {
			if start >= end {
				return -1
			}
			res := arr[start]
			for i := start + 1; i < end; i++ {
				if res > arr[i] {
					res = arr[i]
				}
			}
			return res
		}

		maxBruteForce := func(start, end int32) int32 {
			if start >= end {
				return -1
			}

			res := arr[start]
			for i := start + 1; i < end; i++ {
				if res < arr[i] {
					res = arr[i]
				}
			}
			return res
		}
		_ = maxBruteForce

		countBruteForce := func(start, end, c int32) int32 {
			res := int32(0)
			for i := start; i < end; i++ {
				if arr[i] == c {
					res++
				}
			}
			return res
		}
		_ = countBruteForce

		countAllBruteForce := func(start, end int32) [SIGMA]int32 {
			res := [SIGMA]int32{}
			for i := start; i < end; i++ {
				res[arr[i]-97]++
			}
			return res
		}
		_ = countAllBruteForce

		setBruteForce := func(index, value int32) {
			arr[index] = value
		}
		_ = setBruteForce

		getBruteForce := func(index int32) int32 {
			return arr[index]
		}
		_ = getBruteForce

		for j := 0; j < 10000; j++ {
			start, end := int32(rand.Intn(int(n+1))), int32(rand.Intn(int(n+1)))
			if start > end {
				start, end = end, start
			}
			if S.IsAscending(start, end) != isAscendingBruteForce(start, end) {
				panic("IsAscending Error")
			}
			if S.IsDescending(start, end) != isDescendingBruteForce(start, end) {
				fmt.Println(start, end, S.IsDescending(start, end), isDescendingBruteForce(start, end))
				panic("IsDescending Error")
			}

			if S.Min(start, end) != minBruteForce(start, end) {
				fmt.Println(start, end)
				fmt.Println(S.Min(start, end), minBruteForce(start, end))
				panic("Min Error")
			}

			if S.Max(start, end) != maxBruteForce(start, end) {
				fmt.Println(start, end)
				fmt.Println(S.Max(start, end), maxBruteForce(start, end))
				panic("Max Error")
			}

			c := int32(rand.Intn(26) + 97)
			if S.Count(start, end, c) != countBruteForce(start, end, c) {
				fmt.Println(start, end, c)
				fmt.Println(S.Count(start, end, c), countBruteForce(start, end, c))
				panic("Count Error")
			}

			if res := S.CountAll(start, end); res != countAllBruteForce(start, end) {
				fmt.Println(start, end)
				fmt.Println(res, countAllBruteForce(start, end))
				panic("CountAll Error")
			}

			index := int32(rand.Intn(int(n)))
			value := int32(rand.Intn(26) + 97)
			setBruteForce(index, value)
			S.Set(index, value)

			if S.Get(index) != getBruteForce(index) {
				fmt.Println(index)
				fmt.Println(S.Get(index), getBruteForce(index))
				panic("Get Error")
			}
		}

	}
	fmt.Println("Pass")
}

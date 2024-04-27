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
	"math/rand"
)

func main() {
	test()
}

const SIGMA int32 = 26
const OFFSET int32 = 97

type AlphaPresumDynamic struct {
	n    int32
	arr  []int32
	tree [SIGMA]*BITArray32
}

func NewAlphaPresumDynamic(n int32, f func(i int32) int32) *AlphaPresumDynamic {
	arr := make([]int32, n)
	for i := int32(0); i < n; i++ {
		arr[i] = f(i) - OFFSET
	}
	tree := [SIGMA]*BITArray32{}
	for i := int32(0); i < SIGMA; i++ {
		tree[i] = NewBitArray(n)
	}
	for i := int32(0); i < n; i++ {
		tree[arr[i]].Add(i, 1)
	}
	res := &AlphaPresumDynamic{n: n, arr: arr, tree: tree}
	return res
}

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
	ptr := end - 1
	for i := SIGMA - 1; i >= 0; i-- {
		c := a.tree[i].QueryRange(start, end)
		ptr -= c
		if ptr < start || a.tree[i].QueryRange(ptr, end) != c {
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
	a.tree[oldValue].Add(index, -1)
	a.tree[value].Add(index, 1)
	a.arr[index] = value
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

// !Point Add Range Sum, 0-based.
type BITArray32 struct {
	n     int32
	total int32
	data  []int32
}

func NewBitArray(n int32) *BITArray32 {
	res := &BITArray32{n: n, data: make([]int32, n)}
	return res
}

func NewBitArrayFrom(n int32, f func(i int32) int32) *BITArray32 {
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
	return &BITArray32{n: n, total: total, data: data}
}

func (b *BITArray32) Add(index int32, v int32) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray32) QueryPrefix(end int32) int32 {
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
func (b *BITArray32) QueryRange(start, end int32) int32 {
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

func (b *BITArray32) QueryAll() int32 {
	return b.total
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
			// if S.IsDescending(start, end) != isDescendingBruteForce(start, end) {
			// 	panic("IsDescending Error")
			// }

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

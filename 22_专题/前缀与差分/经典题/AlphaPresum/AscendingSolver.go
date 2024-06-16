package main

import (
	"cmp"
	"fmt"
	"math/bits"
	"math/rand"
	"strings"
)

func main() {
	arr := []int{3, 2, 3}
	less := func(a, b int) bool { return a < b }
	solver := NewAscendingSolver(
		int32(len(arr)), func(i int32) int { return arr[i] },
		less,
	)
	// up，down 信息在后面那个数字上
	fmt.Println(solver.down) // [1] => arr[1]>arr[0]
	fmt.Println(solver.up)   // [2] => arr[2]>arr[1]

	test()
}

// 树状数组维护区间递增/区间递减
// 区间元素个数<=1，视为递增/递减.
type AscendingSolver[V cmp.Ordered] struct {
	n    int32
	arr  []V
	less func(a, b V) bool
	down *bITArray01 // down[i] = 1 表示 arr[i-1] > arr[i]
	up   *bITArray01 // up[i] = 1 表示 arr[i-1] < arr[i]
}

func NewAscendingSolver[V cmp.Ordered](
	n int32, f func(i int32) V, less func(a, b V) bool,
) *AscendingSolver[V] {
	arr := make([]V, n)
	for i := int32(0); i < n; i++ {
		arr[i] = f(i)
	}
	solver := &AscendingSolver[V]{n: n, arr: arr, less: less}
	down, up := newBITArray01(n), newBITArray01(n)
	for i := int32(1); i < n; i++ {
		if less(arr[i-1], arr[i]) {
			up.Add(i)
		}
		if less(arr[i], arr[i-1]) {
			down.Add(i)
		}
	}
	solver.down, solver.up = down, up
	return solver
}

func (solver *AscendingSolver[V]) Set(i int32, v V) {
	if solver.arr[i] == v {
		return
	}
	if i > 0 {
		if solver.less(v, solver.arr[i-1]) {
			solver.down.Add(i)
		} else {
			solver.down.Remove(i)
		}
		if solver.less(solver.arr[i-1], v) {
			solver.up.Add(i)
		} else {
			solver.up.Remove(i)
		}
	}
	if i+1 < solver.n {
		if solver.less(solver.arr[i+1], v) {
			solver.down.Add(i + 1)
		} else {
			solver.down.Remove(i + 1)
		}
		if solver.less(v, solver.arr[i+1]) {
			solver.up.Add(i + 1)
		} else {
			solver.up.Remove(i + 1)
		}
	}
	solver.arr[i] = v
}

func (solver *AscendingSolver[V]) Get(i int32) V { return solver.arr[i] }

// 区间元素个数<=1，视为递增.
func (solver *AscendingSolver[V]) IsAscending(start, end int32) bool {
	if start < 0 {
		start = 0
	}
	if end > solver.n {
		end = solver.n
	}
	if start >= end {
		return true
	}
	return solver.down.QueryRange(start+1, end) == 0
}

// 区间元素个数<=1，视为递减.
func (solver *AscendingSolver[V]) IsDescending(start, end int32) bool {
	if start < 0 {
		start = 0
	}
	if end > solver.n {
		end = solver.n
	}
	if start >= end {
		return true
	}
	return solver.up.QueryRange(start+1, end) == 0
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
		if end == bit01.n {
			return bit01.QueryAll()
		}
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
	for i := 0; i < 20; i++ {
		n := rand.Intn(500) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(100)
		}

		less := func(a, b int) bool { return a < b }
		solver := NewAscendingSolver(
			int32(n), func(i int32) int { return arr[i] },
			less,
		)

		isAscending := func(start, end int32) bool {
			for i := start + 1; i < end; i++ {
				if less(arr[i], arr[i-1]) {
					return false
				}
			}
			return true
		}
		_ = isAscending

		isDescending := func(start, end int32) bool {
			for i := start + 1; i < end; i++ {
				if less(arr[i-1], arr[i]) {
					return false
				}
			}
			return true
		}
		_ = isDescending

		for s := 0; s < 100; s++ {

			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					res1, res2 := solver.IsAscending(int32(i), int32(j)), isAscending(int32(i), int32(j))
					if res1 != res2 {
						fmt.Println("Error1", i, j, res1, res2)
						fmt.Println(arr)
						panic("error1")
					}
					res1, res2 = solver.IsDescending(int32(i), int32(j)), isDescending(int32(i), int32(j))
					if res1 != res2 {
						fmt.Println("Error2")
						panic("error2")
					}
				}
			}

			index := int32(rand.Intn(n))
			v := rand.Intn(10)
			solver.Set(index, v)
			arr[index] = v

			// get
			for i := 0; i < n; i++ {
				res1, res2 := solver.Get(int32(i)), arr[i]
				if res1 != res2 {
					fmt.Println("Error3", i, res1, res2)
					fmt.Println(arr)
					panic("error3")
				}
			}
		}
	}

	fmt.Println("pass")
}

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	v1, v2 := NewLinearBase(nil), NewLinearBase(nil)
	v1.Add(1)
	v1.Add(2)
	v1.Add(3)
	v2.Add(1)
	v2.Add(5)
	v2.Add(7)
	v2.Add(8)

	fmt.Println(v1.Or(v2))
	fmt.Println(v1)
	fmt.Println(v1.And(v2))
	fmt.Println(v1.And(NewLinearBase(nil)))
	fmt.Println(v1)
}

func demo() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	xor := 0
	V1 := NewLinearBase(nil)
	for _, v := range nums {
		xor ^= v
		V1.Add(v)
	}

	mask := ^xor
	V2 := NewLinearBase(nil)
	V1.ForEach(func(v int) {
		V2.Add(v & mask)
	})

	res := V2.Max(0) + (xor ^ V2.Max(0))
	fmt.Fprintln(out, res)

}

// VectorSpace，线性基空间.支持线性基合并.
type LinearBase struct {
	bases  []int
	maxBit int
}

func NewLinearBase(nums []int) *LinearBase {
	res := &LinearBase{}
	for _, num := range nums {
		res.Add(num)
	}
	return res
}

func (lb *LinearBase) Add(num int) bool {
	for _, base := range lb.bases {
		if base == 0 || num == 0 {
			break
		}
		num = min(num, num^base)
	}
	if num != 0 {
		lb.bases = append(lb.bases, num)
		lb.maxBit = max(lb.maxBit, num)
		return true
	}
	return false
}

func (lb *LinearBase) Max(xor int) int {
	res := xor
	for _, base := range lb.bases {
		res = max(res, res^base)
	}
	return res
}

func (lb *LinearBase) Min(xorVal int) int {
	res := xorVal
	for _, base := range lb.bases {
		res = min(res, res^base)
	}
	return res
}

func (lb *LinearBase) Copy() *LinearBase {
	res := &LinearBase{}
	res.bases = append(res.bases, lb.bases...)
	res.maxBit = lb.maxBit
	return res
}

func (lb *LinearBase) Len() int {
	return len(lb.bases)
}

func (lb *LinearBase) ForEach(f func(base int)) {
	for _, base := range lb.bases {
		f(base)
	}
}

func (lb *LinearBase) Has(v int) bool {
	for _, w := range lb.bases {
		if v == 0 {
			break
		}
		v = min(v, v^w)
	}
	return v == 0
}

func (lb *LinearBase) Or(other *LinearBase) *LinearBase {
	v1, v2 := lb, other
	if v1.Len() < v2.Len() {
		v1, v2 = v2, v1
	}
	res := v1.Copy()
	for _, base := range v2.bases {
		res.Add(base)
	}
	return res
}

func (lb *LinearBase) And(other *LinearBase) *LinearBase {
	maxDim := max(lb.maxBit, other.maxBit)
	x := lb.orthogonalSpace(maxDim)
	y := other.orthogonalSpace(maxDim)
	if x.Len() < y.Len() {
		x, y = y, x
	}
	for _, base := range y.bases {
		x.Add(base)
	}
	return x.orthogonalSpace(maxDim)
}

func (lb *LinearBase) String() string {
	return fmt.Sprintf("%v", lb.bases)
}

// 正交空间.
func (lb *LinearBase) orthogonalSpace(maxDim int) *LinearBase {
	lb.normalize(true)
	m := maxDim
	tmp := make([]int, m)
	for _, base := range lb.bases {
		tmp[bits.Len(uint(base))-1] = base
	}
	tmp = Transpose(m, m, tmp, true)
	res := &LinearBase{}
	for j := 0; j < m; j++ {
		if tmp[j]>>j&1 == 1 {
			continue
		}
		res.Add(tmp[j] | 1<<j)
	}
	return res
}

func (lb *LinearBase) normalize(reverse bool) {
	n := len(lb.bases)
	for j := 0; j < n; j++ {
		for i := 0; i < j; i++ {
			lb.bases[i] = min(lb.bases[i], lb.bases[i]^lb.bases[j])
		}
	}
	if !reverse {
		sort.Ints(lb.bases)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(lb.bases)))
	}
}

// 矩阵转置,O(n+m)log(n+m)
func Transpose(row, col int, grid []int, inPlace bool) []int {
	if len(grid) != row {
		panic("row not match")
	}
	if !inPlace {
		grid = append(grid[:0:0], grid...)
	}
	log := 0
	max_ := max(row, col)
	for 1<<log < max_ {
		log++
	}
	if len(grid) < 1<<log {
		*&grid = append(grid, make([]int, 1<<log-len(grid))...)
	}
	width := 1 << log
	mask := int(1)
	for i := 0; i < log; i++ {
		mask |= (mask << (1 << i))
	}
	for t := 0; t < log; t++ {
		width >>= 1
		mask ^= (mask >> width)
		for i := 0; i < 1<<t; i++ {
			for j := 0; j < width; j++ {
				x := &grid[width*(2*i)+j]
				y := &grid[width*(2*i+1)+j]
				*x = ((*y << width) & mask) ^ *x
				*y = ((*x & mask) >> width) ^ *y
				*x = ((*y << width) & mask) ^ *x
			}
		}
	}
	return grid[:col]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

// G - Partial Xor Enumeration
// https://atcoder.jp/contests/abc283/submissions/me
func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, left, right int
	fmt.Fscan(in, &n, &left, &right)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	lb := NewLinearBaseFromList(nums)
	for i := left; i <= right; i++ {
		fmt.Fprint(out, lb.KthXor(i), " ")
	}
}

type LinearBase struct {
	Bases []int       // 基底
	Rows  map[int]int // 高斯消元的行
	Bit   int         // 最大数的位数
}

func NewLinearBase(bit int) *LinearBase {
	return &LinearBase{
		Rows: make(map[int]int),
		Bit:  bit,
	}
}

func NewLinearBaseFromList(list []int) *LinearBase {
	maxNum := 0
	for _, num := range list {
		if num > maxNum {
			maxNum = num
		}
	}
	res := NewLinearBase(bits.Len(uint(maxNum)))
	for _, num := range list {
		res.Add(num)
	}
	res.Build()
	return res
}

// 插入一个向量,如果插入成功返回True,否则返回False.
func (lb *LinearBase) Add(x int) bool {
	x = lb.Normalize(x)
	if x == 0 {
		return false
	}
	i := bits.Len(uint(x)) - 1
	for j := 0; j < lb.Bit; j++ {
		if (lb.Rows[j]>>i)&1 == 1 {
			lb.Rows[j] ^= x
		}
	}
	lb.Rows[i] = x
	return true
}

func (lb *LinearBase) Build() {
	res := make([]int, 0)
	sortedKeys := make([]int, 0, len(lb.Rows))
	for k := range lb.Rows {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	for _, k := range sortedKeys {
		if v := lb.Rows[k]; v > 0 {
			res = append(res, v)
		}
	}
	lb.Bases = res
}

// 子序列(子集,包含空集)第k小的异或 1<=k<=2**len(self.bases).
func (lb *LinearBase) KthXor(k int) int {
	k -= 1
	res := 0
	for i := 0; i < bits.Len(uint(k)); i++ {
		if (k>>i)&1 == 1 {
			res ^= lb.Bases[i]
		}
	}
	return res
}

func (lb *LinearBase) MaxXor() int {
	return lb.KthXor(1 << len(lb.Bases))
}

func (lb *LinearBase) Normalize(x int) int {
	for i := bits.Len(uint(x)) - 1; i >= 0; i-- {
		if (x>>i)&1 == 1 {
			x ^= lb.Rows[i]
		}
	}
	return x
}

// x是否能由线性基表出.
func (lb *LinearBase) Has(x int) bool {
	return lb.Normalize(x) == 0
}

func (lb *LinearBase) Copy() *LinearBase {
	res := &LinearBase{
		Bases: append(lb.Bases[:0:0], lb.Bases...),
		Rows:  make(map[int]int, len(lb.Rows)),
		Bit:   lb.Bit,
	}
	for k, v := range lb.Rows {
		res.Rows[k] = v
	}
	return res
}

func (lb *LinearBase) Len() int {
	return len(lb.Bases)
}

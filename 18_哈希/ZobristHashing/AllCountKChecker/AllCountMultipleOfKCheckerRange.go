package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// https://www.luogu.com.cn/problem/CF1746F
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	checker := NewAllCountMultipleOfKChecker(nums)
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 1 {
			var index, value int
			fmt.Fscan(in, &index, &value)
			index--
			checker.Set(index, value)
		} else {
			var start, end, k int
			fmt.Fscan(in, &start, &end, &k)
			start--
			if checker.Query(start, end, k) {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

type Value = int

// 判断数组区间每个数出现的次数是否均k的`倍数`.
// !思路是：给每个数分配一个随机的正整数哈希值，判断区间内的哈希值之和是否是k的倍数即可
// 误判概率大，需要分配多个哈希值.
type AllCountMultipleOfKCheckerRange struct {
	arr    []Value
	pool   map[Value]E
	random *Random
	bit    *BITGroup
}

func NewAllCountMultipleOfKChecker(arr []Value) *AllCountMultipleOfKCheckerRange {
	arr = append(arr[:0:0], arr...)
	pool := make(map[Value]E)
	random := NewRandom()
	res := &AllCountMultipleOfKCheckerRange{arr: arr, pool: pool, random: random}
	bit := NewBITGroupFrom(len(arr), func(index int) E { return res.hash(arr[index]) })
	res.bit = bit
	return res
}

func (c *AllCountMultipleOfKCheckerRange) Set(index int, newValue Value) {
	if c.arr[index] == newValue {
		return
	}
	curH, preH := c.hash(newValue), c.hash(c.arr[index])
	diff := op(curH, inv(preH))
	c.bit.Update(index, diff)
	c.arr[index] = newValue
}

func (c *AllCountMultipleOfKCheckerRange) Query(start, end int, k int) bool {
	sum := c.bit.QueryRange(start, end)
	for i := range sum {
		if sum[i]%k != 0 {
			return false
		}
	}
	return true
}

func (c *AllCountMultipleOfKCheckerRange) hash(v Value) E {
	if v, has := c.pool[v]; has {
		return v
	}
	var h E
	for i := range h {
		h[i] = int(c.random.RandInt(1, 1<<32-1))
	}
	c.pool[v] = h
	return h
}

type E = [25]int // 哈希

func e() E {
	return E{}
}
func op(e1, e2 E) E {
	for i := range e1 {
		e1[i] += e2[i]
	}
	return e1
}
func inv(e E) E {
	for i := range e {
		e[i] = -e[i]
	}
	return e
}

type BITGroup struct {
	n     int
	data  []E
	total E
}

func NewBITGroup(n int) *BITGroup {
	data := make([]E, n)
	for i := range data {
		data[i] = e()
	}
	return &BITGroup{n: n, data: data, total: e()}
}

func NewBITGroupFrom(n int, f func(index int) E) *BITGroup {
	total := e()
	data := make([]E, n)
	for i := range data {
		data[i] = f(i)
		total = op(total, data[i])
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[j-1], data[i-1])
		}
	}
	return &BITGroup{n: n, data: data, total: total}
}

func (fw *BITGroup) Update(i int, x E) {
	fw.total = op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

func (fw *BITGroup) QueryAll() E { return fw.total }

// [0, end)
func (fw *BITGroup) QueryPrefix(end int) E {
	if end > fw.n {
		end = fw.n
	}
	res := e()
	for end > 0 {
		res = op(res, fw.data[end-1])
		end &= end - 1
	}
	return res
}

// [start, end)
func (fw *BITGroup) QueryRange(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	if start > end {
		return e()
	}
	pos, neg := e(), e()
	for end > start {
		pos = op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = op(neg, fw.data[start-1])
		start &= start - 1
	}
	return op(pos, inv(neg))
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
	return r.seed
}
func (r *Random) Rng61() uint64 { return r.Rng() & ((1 << 61) - 1) }

// [left, right]
func (r *Random) RandInt(min, max int) uint64 { return uint64(min) + r.Rng()%(uint64(max-min+1)) }

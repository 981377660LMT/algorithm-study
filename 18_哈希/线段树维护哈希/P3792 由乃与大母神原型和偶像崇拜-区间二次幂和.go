// P3792 由乃与大母神原型和偶像崇拜-线段树维护区间二次幂(pow2)的和
// 线段树维护区间 pow2Sum

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"os"
	"strconv"
)

const INF int = 1e18
const MOD int = 1e9 + 7

var fp *FastPow

func init() {
	fp = NewFastPow(2, 1e9+10)
}

// 1 start end : 判断区间[start, end)是否可以重排为值域上连续的一段(区间值域连续)
// 2 pos val : 将 nums[pos] 修改为 val
//
// 通过求区间的 max,min 求出这个区间的左右端点，然后用公式验证这个区间的和是否符合连续区间的性质
func 由乃与大母神原型和偶像崇拜(nums []int, operations [][3]int) []bool {
	leaves := make([]E, len(nums))
	for i := range nums {
		leaves[i] = FromElement(nums[i])
	}
	tree := _NewSegTree(leaves)

	res := []bool{}
	for _, op := range operations {
		kind := op[0]
		if kind == 1 {
			pos, val := op[1], op[2]
			tree.Set(pos, FromElement(val))
		} else {
			start, end := op[1], op[2]
			cur := tree.Query(start, end)
			min, max := cur.min, cur.max
			pow2 := cur.pow2
			ok1 := max-min+1 == end-start
			ok2 := fp.RangePow2Sum(min, max+1) == pow2
			ok := ok1 && ok2
			res = append(res, ok)
		}
	}

	return res
}

// 光速幂.
type FastPow struct {
	max     int
	divData []int
	modData []int
}

// O(sqrt(maxN))预处理,O(1)查询.
//
//	base: 幂运算的基.
//	maxN: 最大的幂.
func NewFastPow(base int, maxN int) *FastPow {
	max := int(math.Ceil(math.Sqrt(float64(maxN))))
	res := &FastPow{max: max, divData: make([]int, max+1), modData: make([]int, max+1)}
	cur := 1
	for i := 0; i <= max; i++ {
		res.modData[i] = cur
		cur = cur * base % MOD
	}
	cur = 1
	last := res.modData[max]
	for i := 0; i <= max; i++ {
		res.divData[i] = cur
		cur = cur * last % MOD
	}
	return res
}

// n<=maxN.
func (fp *FastPow) Pow(n int) int {
	return (fp.divData[n/fp.max] * fp.modData[n%fp.max] % MOD)
}

// 区间以2为底的幂和 (2^start + 2^(start+1) + ... + 2^(end-1)) % MOD.
func (fp *FastPow) RangePow2Sum(start, end int) int {
	if start >= end {
		return 0
	}
	res := (fp.Pow(end) - fp.Pow(start)) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}

type E = struct {
	min, max int
	pow2     int
}

func FromElement(v int) E {
	return E{
		min: v, max: v,
		pow2: fp.Pow(v),
	}
}

func (*_SegTree) e() E {
	return E{
		min: INF, max: -INF,
		pow2: 0,
	}
}

func (*_SegTree) op(a, b E) E {
	return E{
		min: min(a.min, b.min), max: max(a.max, b.max),
		pow2: (a.pow2 + b.pow2) % MOD,
	}
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

type _SegTree struct {
	n, size int
	data    []E
}

func _NewSegTree(leaves []E) *_SegTree {
	res := &_SegTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}

	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.data = seg
	return res
}
func (st *_SegTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.data[index+st.size]
}
func (st *_SegTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.data[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.data[index] = st.op(st.data[index<<1], st.data[index<<1|1])
	}
}
func (st *_SegTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.data[index] = st.op(st.data[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.data[index] = st.op(st.data[index<<1], st.data[index<<1|1])
	}
}

// [start, end)
func (st *_SegTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.data[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.data[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}

func (st *_SegTree) QueryAll() E { return st.data[1] }

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, q := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}
	operations := make([][3]int, q)
	for i := 0; i < q; i++ {
		op := io.NextInt()
		if op == 1 {
			pos, val := io.NextInt(), io.NextInt()
			pos--
			operations[i] = [3]int{op, pos, val}
		} else {
			start, end := io.NextInt(), io.NextInt()
			start--
			operations[i] = [3]int{op, start, end}
		}
	}

	res := 由乃与大母神原型和偶像崇拜(nums, operations)
	for _, ok := range res {
		if ok {
			io.Println("damushen")
		} else {
			io.Println("yuanxing")
		}
	}
}

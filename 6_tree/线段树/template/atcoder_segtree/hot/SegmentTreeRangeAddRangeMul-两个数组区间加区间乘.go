package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

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

const MOD int = 998244353

// F - Two Sequence Queries
// https://atcoder.jp/contests/abc357/tasks/abc357_f
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, q := io.NextInt(), io.NextInt()
	A := make([]int, n)
	B := make([]int, n)
	for i := 0; i < n; i++ {
		A[i] = io.NextInt()
	}
	for i := 0; i < n; i++ {
		B[i] = io.NextInt()
	}

	seg := NewSegmentTreeRangeAddRangeMul(n, func(i int) E {
		return E{
			aSum: A[i],
			bSum: B[i],
			mul:  A[i] * B[i] % MOD,
			size: 1,
		}
	})

	for i := 0; i < q; i++ {
		op := io.NextInt()
		if op == 1 {
			l, r, x := io.NextInt(), io.NextInt(), io.NextInt()
			l--
			seg.Update(l, r, Id{x: x})
		} else if op == 2 {
			l, r, y := io.NextInt(), io.NextInt(), io.NextInt()
			l--
			seg.Update(l, r, Id{y: y})
		} else {
			l, r := io.NextInt(), io.NextInt()
			l--
			res := seg.Query(l, r)
			io.Println(res.mul % MOD)
		}
	}
}

// RangeAddRangeMul(线段树维护两个数组的区间加区间乘积)
// !(a+x)*(b+y) = ab + ay + xb + xy

const INF = 1e18

type E = struct{ aSum, bSum, mul, size int }
type Id = struct{ x, y int }

func (*SegmentTreeRangeAddRangeMul) e() E   { return E{} }
func (*SegmentTreeRangeAddRangeMul) id() Id { return Id{} }
func (*SegmentTreeRangeAddRangeMul) op(left, right E) E {
	left.aSum = (left.aSum + right.aSum) % MOD
	left.bSum = (left.bSum + right.bSum) % MOD
	left.mul = (left.mul + right.mul) % MOD
	left.size += right.size
	return left
}
func (*SegmentTreeRangeAddRangeMul) mapping(f Id, g E) E {
	return E{
		aSum: (g.aSum + f.x*g.size) % MOD,
		bSum: (g.bSum + f.y*g.size) % MOD,
		// !(a+x)*(b+y) = ab + ay + xb + xy
		mul:  (g.mul + f.x*g.bSum + f.y*g.aSum + f.x*f.y%MOD*g.size) % MOD,
		size: g.size,
	}
}
func (*SegmentTreeRangeAddRangeMul) composition(f, g Id) Id {
	f.x = (f.x + g.x) % MOD
	f.y = (f.y + g.y) % MOD
	return f
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// !template
type SegmentTreeRangeAddRangeMul struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewSegmentTreeRangeAddRangeMul(n int, f func(int) E) *SegmentTreeRangeAddRangeMul {
	tree := &SegmentTreeRangeAddRangeMul{}
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

func NewSegmentTreeRangeAddRangeMulFrom(leaves []E) *SegmentTreeRangeAddRangeMul {
	tree := &SegmentTreeRangeAddRangeMul{}
	n := len(leaves)
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = leaves[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *SegmentTreeRangeAddRangeMul) Query(left, right int) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *SegmentTreeRangeAddRangeMul) QueryAll() E {
	return tree.data[1]
}
func (tree *SegmentTreeRangeAddRangeMul) GetAll() []E {
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *SegmentTreeRangeAddRangeMul) Update(left, right int, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := 1; i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *SegmentTreeRangeAddRangeMul) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *SegmentTreeRangeAddRangeMul) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *SegmentTreeRangeAddRangeMul) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *SegmentTreeRangeAddRangeMul) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *SegmentTreeRangeAddRangeMul) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *SegmentTreeRangeAddRangeMul) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := 0; i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}

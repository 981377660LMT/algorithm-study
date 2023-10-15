// !区间仿射变换，区间查询

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"os"
	"strconv"
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

// https://atcoder.jp/contests/practice2/tasks/practice2_k
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

	sqrt := NewSqrtDecomposition(nums, 1+int(math.Sqrt(float64(n))))
	for i := 0; i < q; i++ {
		op := io.NextInt()
		if op == 0 {
			start, end, mul, add := io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt()
			sqrt.Update(start, end, Id{mul: mul, add: add})
		} else {
			start, end := io.NextInt(), io.NextInt()
			res := 0
			sqrt.Query(start, end, func(cur E) {
				res += cur
				res %= MOD
			})
			if res < 0 {
				res += MOD
			}
			io.Println(res)
		}
	}
}

type E = int                     // 需要查询的值的类型
type Id = struct{ mul, add int } // 懒标记的类型

type Block struct {
	id, start, end int // 0 <= id < bs, 0 <= start <= end <= n
	nums           []E // block内的原序列

	sum  E
	lazy Id
}

func (b *Block) Init() {
	b.lazy = Id{mul: 1, add: 0}
	for i := range b.nums {
		b.sum += b.nums[i]
	}
	b.sum %= MOD
}

func (b *Block) QueryPart(start, end int) E {
	res := 0
	add, mul := b.lazy.add, b.lazy.mul
	for i := start; i < end; i++ {
		res += b.nums[i]*mul + add
		res %= MOD
	}
	return res
}

func (b *Block) UpdatePart(start, end int, lazy Id) {
	// !更新[start,end)前将区间标记下放
	b.sum = 0
	bMul, bAdd := b.lazy.mul, b.lazy.add
	for i := 0; i < start; i++ {
		b.nums[i] = b.nums[i]*bMul + bAdd
		b.nums[i] %= MOD
		b.sum += b.nums[i]
	}
	mul, add := lazy.mul, lazy.add
	for i := start; i < end; i++ {
		b.nums[i] = b.nums[i]*bMul + bAdd
		b.nums[i] %= MOD
		b.nums[i] = b.nums[i]*mul + add
		b.nums[i] %= MOD
		b.sum += b.nums[i]
	}
	for i := end; i < len(b.nums); i++ {
		b.nums[i] = b.nums[i]*bMul + bAdd
		b.nums[i] %= MOD
		b.sum += b.nums[i]
	}
	b.sum %= MOD
	b.lazy = Id{mul: 1, add: 0}
}

func (b *Block) QueryAll() E {
	mul, add := b.lazy.mul, b.lazy.add
	return (b.sum*mul + add*(b.end-b.start)) % MOD
}

func (b *Block) UpdateAll(lazy Id) {
	b.lazy.mul = b.lazy.mul * lazy.mul % MOD
	b.lazy.add = (b.lazy.add*lazy.mul + lazy.add) % MOD
}

//
//
//
// dont modify the template below
//
//
//

type SqrtDecomposition struct {
	n      int
	bs     int
	bls    []Block
	belong []int
}

// 指定维护的序列和分块大小初始化.
//
//	blockSize:分块大小,一般取根号n(300)
func NewSqrtDecomposition(nums []E, blockSize int) *SqrtDecomposition {
	nums = append(nums[:0:0], nums...)
	res := &SqrtDecomposition{
		n:      len(nums),
		bs:     blockSize,
		bls:    make([]Block, len(nums)/blockSize+1),
		belong: make([]int, len(nums)),
	}
	for i := range res.belong {
		res.belong[i] = i / blockSize
	}
	for i := range res.bls {
		res.bls[i].id = i
		res.bls[i].start = i * blockSize
		res.bls[i].end = min((i+1)*blockSize, len(nums))
		res.bls[i].nums = nums[res.bls[i].start:res.bls[i].end]
		res.bls[i].Init()
	}
	return res
}

// 更新左闭右开区间[start,end)的值.
//
//	0<=start<=end<=n
func (s *SqrtDecomposition) Update(start, end int, lazy Id) {
	if start >= end {
		return
	}
	id1, id2 := s.belong[start], s.belong[end-1]
	pos1, pos2 := start-s.bs*id1, end-s.bs*id2
	if id1 == id2 {
		s.bls[id1].UpdatePart(pos1, pos2, lazy)
	} else {
		s.bls[id1].UpdatePart(pos1, s.bs, lazy)
		for i := id1 + 1; i < id2; i++ {
			s.bls[i].UpdateAll(lazy)
		}
		s.bls[id2].UpdatePart(0, pos2, lazy)
	}
}

// 查询左闭右开区间[start,end)的值.
//
//	0<=start<=end<=n
func (s *SqrtDecomposition) Query(start, end int, cb func(blockRes E)) {
	if start >= end {
		return
	}
	id1, id2 := s.belong[start], s.belong[end-1]
	pos1, pos2 := start-s.bs*id1, end-s.bs*id2
	if id1 == id2 {
		cb(s.bls[id1].QueryPart(pos1, pos2))
		return
	}
	cb(s.bls[id1].QueryPart(pos1, s.bs))
	for i := id1 + 1; i < id2; i++ {
		cb(s.bls[i].QueryAll())
	}
	cb(s.bls[id2].QueryPart(0, pos2))
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

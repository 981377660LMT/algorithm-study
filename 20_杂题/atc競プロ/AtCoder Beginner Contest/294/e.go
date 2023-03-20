package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

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

//
//
//
// !Range Add Range Sum, 0-based.
type BITArray2 struct {
	n     int
	tree1 map[int]int
	tree2 map[int]int
}

func NewBITArray2(n int) *BITArray2 {
	return &BITArray2{
		n:     n,
		tree1: make(map[int]int),
		tree2: make(map[int]int),
	}
}

// 切片内[start, end)的每个元素加上delta.
//  0<=start<=end<=n
func (b *BITArray2) Add(start, end, delta int) {
	end--
	b.add(start, delta)
	b.add(end+1, -delta)
}

// 求切片内[start, end)的和.
//  0<=start<=end<=n
func (b *BITArray2) Query(start, end int) int {
	end--
	return b.query(end) - b.query(start-1)
}

func (b *BITArray2) add(index, delta int) {
	index++
	rawIndex := index
	for index <= b.n {
		b.tree1[index] += delta
		b.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (b *BITArray2) query(index int) (res int) {
	index++
	if index > b.n {
		index = b.n
	}
	rawIndex := index
	for index > 0 {
		res += rawIndex*b.tree1[index] - b.tree2[index]
		index -= index & -index
	}
	return
}

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	// L, N1, N2 = map(int, input().split())
	// # !快速查询区间里某种数有多少个
	// bits = defaultdict(lambda: BIT2(int(1e12) + 10))  # 数 => BIT
	// start = 1
	// for _ in range(N1):
	// 		v, l = map(int, input().split())
	// 		bits[v].add(start, start + l - 1, 1)
	// 		start += l
	// res = 0
	// start = 1
	// for _ in range(N2):
	// 		v, l = map(int, input().split())
	// 		res += bits[v].query(start, start + l - 1)
	// 		start += l
	// print(res)

	// map 可以用指针吗 /
	// !动态开点线段树谁快
	L, N1, N2 := io.NextInt(), io.NextInt(), io.NextInt()
	bits := make(map[int]BITArray2)
	start := 0
	for i := 0; i < N1; i++ {
		v, l := io.NextInt(), io.NextInt()
		if _, ok := bits[v]; !ok {
			bits[v] = *NewBITArray2(L + 10)
		}
		bit := bits[v]
		bit.Add(start, start+l, 1)
		start += l
	}
	res := 0
	start = 0
	for i := 0; i < N2; i++ {
		v, l := io.NextInt(), io.NextInt()
		if _, ok := bits[v]; !ok {
			bits[v] = *NewBITArray2(L + 10)
		}
		bit := bits[v]
		res += bit.Query(start, start+l)
		start += l
	}

	io.Println(res)

}

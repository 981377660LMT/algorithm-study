package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()
	// 	A∩B=∅ を満たす
	// 2 つの整数の集合
	// A,B に対して、
	// f(A,B) を以下のように定義します。

	// A∪B に含まれる要素を昇順に並べた数列を
	// C=(C
	// 1
	// ​
	//  ,C
	// 2
	// ​
	//  ,…,C
	// ∣A∣+∣B∣
	// ​
	//  ) とする。
	// A={C
	// k
	// 1
	// ​

	// ​
	//  ,C
	// k
	// 2
	// ​

	// ​
	//  ,…,C
	// k
	// ∣A∣
	// ​

	// ​
	//  } となるような
	// k
	// 1
	// ​
	//  ,k
	// 2
	// ​
	//  ,…,k
	// ∣A∣
	// ​
	//   をとる。 このとき、
	// f(A,B)=
	// i=1
	// ∑
	// ∣A∣
	// ​
	//  k
	// i
	// ​
	//   とする。
	// 例えば、
	// A={1,3},B={2,8} のとき、
	// C=(1,2,3,8) より
	// A={C
	// 1
	// ​
	//  ,C
	// 3
	// ​
	//  } なので、
	// f(A,B)=1+3=4 です。

	// それぞれが
	// M 個の要素からなる
	// N 個の整数の集合
	// S
	// 1
	// ​
	//  ,S
	// 2
	// ​
	//  …,S
	// N
	// ​
	//   があり、各
	// i (1≤i≤N) について
	// S
	// i
	// ​
	//  ={A
	// i,1
	// ​
	//  ,A
	// i,2
	// ​
	//  ,…,A
	// i,M
	// ​
	//  } です。 ただし、
	// S
	// i
	// ​
	//  ∩S
	// j
	// ​
	//  =∅ (i
	// 
	// =j) が保証されます。

	// 1≤i<j≤N
	// ∑
	// ​
	//  f(S
	// i
	// ​
	//  ,S
	// j
	// ​
	//  ) を求めてください。

	n, m := io.NextInt(), io.NextInt()
	allNums := make([][]uint32, n)
	nums := make([]uint32, 0, n*m)
	for i := 0; i < n; i++ {
		row := make([]uint32, m)
		for j := 0; j < m; j++ {
			row[j] = uint32(io.NextInt())
		}
		sort.Slice(row, func(i, j int) bool { return row[i] < row[j] })
		allNums[i] = row
		nums = append(nums, row...)
	}
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	rank := make(map[uint32]int, n*m)
	for i, v := range nums {
		rank[v] = i
	}

	res := 0
	leaves := make([]int, n*m)
	for i := range leaves {
		leaves[i] = 1
	}
	bit := NewBitArrayFrom(leaves)
	for _, row := range allNums {
		for i := 0; i < len(row); i++ {
			cur := row[i]
			curRank := rank[cur]
			bit.Add(curRank, -1)
		}
		for i := len(row) - 1; i >= 0; i-- {
			cur := row[i]
			curRank := rank[cur]
			res += bit.Query(curRank + 1)
		}
	}

	io.Println(res + m*(m+1)/2*n*(n-1)/2)
}

type BitArray struct {
	n    int
	log  int
	data []int
}

func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

func NewBitArrayFrom(arr []int) *BitArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

func (b *BitArray) Add(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r)
func (b *BitArray) Query(r int) int {
	res := 0
	for ; r > 0; r &= r - 1 {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BitArray) QueryRange(l, r int) int {
	return b.Query(r) - b.Query(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *BitArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

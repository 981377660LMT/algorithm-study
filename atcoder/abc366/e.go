package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, d := io.NextInt(), io.NextInt()
	xs, ys := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		xs[i], ys[i] = io.NextInt(), io.NextInt()
	}

	sort.Ints(xs)
	sort.Ints(ys)

	Dx := DistSum(xs)
	Dy := DistSum(ys)
	minX, maxX := xs[0], xs[n-1]
	midY := ys[n/2]
	distToMidY := Dy(midY)
	res := 0
	lowerX, upperX := minX-d, maxX+d
	for x := lowerX; x <= upperX; x++ {
		distToAllX := Dx(x)
		remainY := d - distToAllX
		if remainY < distToMidY {
			continue
		}

		// 二分查找 y 的范围
		lower1, upper1 := midY, midY+remainY
		for lower1 <= upper1 {
			mid := (lower1 + upper1) / 2
			if Dy(mid) <= remainY {
				lower1 = mid + 1
			} else {
				upper1 = mid - 1
			}
		}
		res += max(0, upper1-midY)
		lower2, upper2 := midY-remainY, midY-1
		for lower2 <= upper2 {
			mid := (lower2 + upper2) / 2
			if Dy(mid) <= remainY {
				upper2 = mid - 1
			} else {
				lower2 = mid + 1
			}
		}
		res += max(0, midY-lower2+1)
	}

	io.Println(res)
}

// 有序数组中所有点到`x=k`的距离之和.
func DistSum(sortedNums []int) func(k int) int {
	n := len(sortedNums)
	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + sortedNums[i]
	}

	return func(k int) int {
		pos := sort.SearchInts(sortedNums, k+1)
		leftSum := k*pos - preSum[pos]
		rightSum := preSum[n] - preSum[pos] - k*(n-pos)
		return leftSum + rightSum
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"sort"
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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	points := make([][2]int, n)
	S := NewRectangleSum()
	for i := 0; i < n; i++ {
		points[i][0], points[i][1] = io.NextInt(), io.NextInt()
		if points[i][0] > points[i][1] {
			points[i][0], points[i][1] = points[i][1], points[i][0]
		}
		S.AddPoint(points[i][0], points[i][1], 1)
	}
	S.Build()

	for i := 0; i < n; i++ {
		curX, curY := points[i][0], points[i][1]
		count := S.Query(curX, curY, curY+1, n+10) + S.Query(curX, curY, 0, curX)
		if count > 0 {
			io.Println("Yes")
			return
		}
	}
	io.Println("No")

}

type StaticRectangleSum struct {
	points [][3]int
	xs     []int
	ys     []int
	wm     *waveletMatrix
}

func NewRectangleSum() *StaticRectangleSum {
	return &StaticRectangleSum{
		points: [][3]int{},
	}
}

func (s *StaticRectangleSum) AddPoint(x, y, w int) {
	s.points = append(s.points, [3]int{x, y, w})
}

func (s *StaticRectangleSum) Build() {
	sort.Slice(s.points, func(i, j int) bool {
		return s.points[i][0] < s.points[j][0]
	})
	n := len(s.points)
	xs, ys, ws := make([]int, n), make([]int, n), make([]int, n)
	for i, p := range s.points {
		xs[i], ys[i], ws[i] = p[0], p[1], p[2]
	}
	s.xs = xs

	set := make(map[int]struct{}, len(ys))
	for _, y := range ys {
		set[y] = struct{}{}
	}
	sortedSet := make([]int, 0, len(set))
	for y := range set {
		sortedSet = append(sortedSet, y)
	}
	sort.Ints(sortedSet)
	s.ys = sortedSet

	comp := make(map[int]int, len(sortedSet))
	for i, y := range sortedSet {
		comp[y] = i
	}

	newYs := make([]int, len(ys))
	for i, y := range ys {
		newYs[i] = comp[y]
	}

	maxLog := bits.Len(uint(len(sortedSet)))
	s.wm = newWaveletMatrix(newYs, ws, maxLog)
}

// 求矩形x1<=x<x2,y1<=y<y2的权值和 注意是左闭右开
func (s *StaticRectangleSum) Query(x1, x2, y1, y2 int) int {
	return s.rectSum(x1, x2, y2) - s.rectSum(x1, x2, y1)
}

func (s *StaticRectangleSum) rectSum(left, right, upper int) int {
	left = sort.SearchInts(s.xs, left)
	right = sort.SearchInts(s.xs, right)
	upper = sort.SearchInts(s.ys, upper)
	return s.wm.RectSum(left, right, upper)
}

func sortedSet(nums []int) []int {
	set := make(map[int]struct{}, len(nums))
	for _, num := range nums {
		set[num] = struct{}{}
	}
	res := make([]int, 0, len(set))
	for num := range set {
		res = append(res, num)
	}
	sort.Ints(res)
	return res
}

func newWaveletMatrix(ys, ws []int, maxLog int) *waveletMatrix {
	n := len(ys)
	mat := make([]*bitVector, 0, maxLog)
	zs := make([]int, 0, maxLog)
	data := make([][]int, maxLog)
	for i := range data {
		data[i] = make([]int, n+1)
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	for d := maxLog - 1; d >= 0; d-- {
		vec := newBitVector(n + 1)
		ls, rs := make([]int, 0, n), make([]int, 0, n)
		for i, val := range order {
			if (ys[val]>>uint(d))&1 == 1 {
				rs = append(rs, val)
				vec.Set(i)
			} else {
				ls = append(ls, val)
			}
		}
		vec.Build()
		mat = append(mat, vec)
		zs = append(zs, len(ls))
		order = append(ls, rs...)
		for i, val := range order {
			data[maxLog-d-1][i+1] = data[maxLog-d-1][i] + ws[val]
		}
	}

	return &waveletMatrix{
		n:      n,
		maxLog: maxLog,
		mat:    mat,
		zs:     zs,
		data:   data,
	}
}

type waveletMatrix struct {
	n      int
	maxLog int
	mat    []*bitVector
	zs     []int
	data   [][]int
}

func (w *waveletMatrix) RectSum(left, right, upper int) int {
	res := 0
	for d := 0; d < w.maxLog; d++ {
		if (upper>>(w.maxLog-d-1))&1 == 1 {
			res += w.data[d][w.mat[d].Count(0, right)]
			res -= w.data[d][w.mat[d].Count(0, left)]
			left = w.mat[d].Count(1, left) + w.zs[d]
			right = w.mat[d].Count(1, right) + w.zs[d]
		} else {
			left = w.mat[d].Count(0, left)
			right = w.mat[d].Count(0, right)
		}
	}
	return res
}

type bitVector struct {
	n     int
	block []int
	sum   []int
}

func newBitVector(n int) *bitVector {
	blockCount := (n + 63) >> 6
	return &bitVector{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int, blockCount+1),
	}
}

func (f *bitVector) Set(i int) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *bitVector) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + bits.OnesCount(uint(f.block[i]))
	}
}

// 统计 [0,end) 中 value 的个数
func (f *bitVector) Count(value, end int) int {
	if value == 1 {
		return f.count1(end)
	}
	return end - f.count1(end)
}

func (f *bitVector) count1(k int) int {
	mask := (1 << uint(k&63)) - 1
	return f.sum[k>>6] + bits.OnesCount(uint(f.block[k>>6]&mask))
}

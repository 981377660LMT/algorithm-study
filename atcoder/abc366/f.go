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

const INF int = 4e18

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, k := io.NextInt(), io.NextInt()
	pairs := make([][2]int, n)
	for i := 0; i < n; i++ {
		a, b := io.NextInt(), io.NextInt()
		pairs[i] = [2]int{a, b}
	}

	sort.Slice(pairs, func(i, j int) bool {
		a1, b1 := pairs[i][0], pairs[i][1]
		a2, b2 := pairs[j][0], pairs[j][1]
		return a1*b2 > a2*b1
	})

	// jump or not
	memo := make([][]int, (n+1)*(k+1))
	for i := range memo {
		memo[i] = make([]int, (k + 1))
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var dfs (func(int, int) int)
	dfs = func(index, remain int) int {
		if index == n {
			if remain == 0 {
				return 1
			}
			return -INF
		}
		if remain == 0 {
			return 1
		}
		if memo[index][remain] != -1 {
			return memo[index][remain]
		}
		res := dfs(index+1, remain)
		a, b := pairs[index][0], pairs[index][1]
		nextRes := dfs(index+1, remain-1)
		if nextRes != -INF {
			res = max(res, a*nextRes+b)
		}

		memo[index][remain] = res
		return res
	}

	res := dfs(0, k)
	io.Println(res)
}

type E = struct{ mul, add int }

func (*SlidingWindowAggregation) e() E { return E{1, 0} }
func (*SlidingWindowAggregation) op(a, b E) E {
	return E{a.mul * b.mul, (a.add*b.mul + b.add)}
}

type SlidingWindowAggregation struct {
	cumL []E
	cumR E
	dat  []E
	sz   int
}

func NewSlidingWindowAggregation() *SlidingWindowAggregation {
	res := &SlidingWindowAggregation{}
	res.cumL = []E{res.e()}
	res.cumR = res.e()
	return res
}

func (s *SlidingWindowAggregation) Len() int {
	return s.sz
}

func (s *SlidingWindowAggregation) Append(x E) {
	s.sz++
	s.cumR = s.op(s.cumR, x)
	s.dat = append(s.dat, x)
}

func (s *SlidingWindowAggregation) PopLeft() {
	s.sz--
	s.cumL = s.cumL[:len(s.cumL)-1]
	if len(s.cumL) == 0 {
		s.cumL = []E{s.e()}
		s.cumR = s.e()
		for len(s.dat) > 1 {
			s.cumL = append(s.cumL, s.op(s.dat[len(s.dat)-1], s.cumL[len(s.cumL)-1]))
			s.dat = s.dat[:len(s.dat)-1]
		}
		s.dat = s.dat[:0]
	}
}

func (s *SlidingWindowAggregation) Query() E {
	return s.op(s.cumL[len(s.cumL)-1], s.cumR)
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
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

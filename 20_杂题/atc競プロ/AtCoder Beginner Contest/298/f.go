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

// from collections import defaultdict
// from heapq import heappop, heappush
// import sys

// sys.setrecursionlimit(int(1e9))
// input = lambda: sys.stdin.readline().rstrip("\r\n")
// MOD = 998244353
// INF = int(4e18)

// if __name__ == "__main__":
//     # 站在交叉点上

//     n = int(input())
//     points = [tuple(map(int, input().split())) for _ in range(n)]
//     S = set((x, y) for x, y, _ in points)
//     rowSum, colSum = defaultdict(int), defaultdict(int)
//     for x, y, w in points:
//         rowSum[x] += w
//         colSum[y] += w

//     res = 0
//     for r, c, w in points:
//         res = max(res, rowSum[r] + colSum[c] - w)

//     # 不站在交叉点上
//     # !两个数组选数,最大化和
//     row = sorted(rowSum.items(), key=lambda x: x[1], reverse=True)
//     col = sorted(colSum.items(), key=lambda x: x[1], reverse=True)
//     pq = [(-(row[0][1] + col[0][1]), row[0][0], col[0][0], 0, 0)]
//     while pq:
//         s, x, y, ptr1, ptr2 = heappop(pq)
//         s = -s
//         if (x, y) not in S:
//             res = max(res, s)
//             print(res)
//             exit(0)
//         else:
//             if ptr1 + 1 < len(row):
//                 nextS = s - row[ptr1][1] + row[ptr1 + 1][1]
//                 heappush(pq, (-nextS, row[ptr1 + 1][0], y, ptr1 + 1, ptr2))
//             if ptr2 + 1 < len(col):
//                 nextS = s - col[ptr2][1] + col[ptr2 + 1][1]
//                 heappush(pq, (-nextS, x, col[ptr2 + 1][0], ptr1, ptr2 + 1))
//     print(res)

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	points := make([][3]int, n)
	for i := 0; i < n; i++ {
		x, y, w := io.NextInt(), io.NextInt(), io.NextInt()
		points[i] = [3]int{x, y, w}
	}

	S := make(map[[2]int]struct{})
	rowSum, colSum := make(map[int]int), make(map[int]int)
	for _, p := range points {
		x, y, w := p[0], p[1], p[2]
		S[[2]int{x, y}] = struct{}{}
		rowSum[x] += w
		colSum[y] += w
	}

	res := 0
	for _, p := range points {
		x, y, w := p[0], p[1], p[2]
		res = max(res, rowSum[x]+colSum[y]-w)
	}

	rows := make([][2]int, 0, len(rowSum))
	cols := make([][2]int, 0, len(colSum))
	for k, v := range rowSum {
		rows = append(rows, [2]int{k, v})
	}
	for k, v := range colSum {
		cols = append(cols, [2]int{k, v})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i][1] > rows[j][1] })
	sort.Slice(cols, func(i, j int) bool { return cols[i][1] > cols[j][1] })

	pq := NewHeap(func(a, b [3]int) bool { return a[0] > b[0] }, nil)
	pq.Push([3]int{(rows[0][1] + cols[0][1]), 0, 0})
	for pq.Len() > 0 {
		item := pq.Pop()
		s, ptr1, ptr2 := item[0], item[1], item[2]
		if _, ok := S[[2]int{rows[ptr1][0], cols[ptr2][0]}]; !ok {
			res = max(res, s)
			io.Println(res)
			return
		} else {
			if ptr1+1 < len(rows) {
				nextS := s - rows[ptr1][1] + rows[ptr1+1][1]
				pq.Push([3]int{nextS, ptr1 + 1, ptr2})
			}
			if ptr2+1 < len(cols) {
				nextS := s - cols[ptr2][1] + cols[ptr2+1][1]
				pq.Push([3]int{nextS, ptr1, ptr2 + 1})
			}
		}
	}

	io.Println(res)

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

type H = [3]int

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
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

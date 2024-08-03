package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/rand"
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
	// fmt.Println(N, M)
	// fmt.Println(T)
	// fmt.Println(P)
	// for i := 0; i < 10000; i++ {
	// 	check()
	// }
	// fmt.Println("pass")

	// 	2 6 [2 4 5 6 8 10] [0 1 1 1 1 0] 8 [2 2 2 2 2 2 2 2] [1 1 1 1 1 1 1 1]
	// [8 8 8 8 8 8 8 8] [3 3 3 3 3 3 3 3]
	N, M := 2, 6
	T := []int{2, 4, 5, 6, 8, 10}
	P := []int{0, 1, 1, 1, 1, 0}
	Q := 8
	fmt.Println(solve1(N, M, T, P, Q, []int{2, 2, 2, 2, 2, 2, 2, 2}, []int{1, 1, 1, 1, 1, 1, 1, 1}))
	// fmt.Println(solve2(N, M, T, P, Q, []int{2, 2, 2, 2, 2, 2, 2, 2}, []int{1, 1, 1, 1, 1, 1, 1, 1}))
}

func check() {
	N := 2
	M := 6
	T := make([]int, M)
	for i := 0; i < M; i++ {
		T[i] = rand.Intn(5) + 1
	}
	sort.Ints(T)
	for i := 1; i < M; i++ {
		if T[i] == T[i-1] {
			T[i]++
		}
	}
	remain := make([][2]int, M)
	for i := 0; i < M; i++ {
		remain[i] = [2]int{i, T[i]}
	}

	P := make([]int, M)
	take2 := func(p int) {
		l, r := rand.Intn(len(remain)), rand.Intn(len(remain))
		if l == r {
			r = (r + 1) % len(remain)
		}
		if l > r {
			l, r = r, l
		}
		t1, t2 := remain[l][0], remain[r][0]
		P[t1] = p
		P[t2] = p
		tmp := append(remain[:l], remain[l+1:r]...)
		tmp = append(tmp, remain[r+1:]...)
		remain = tmp
	}

	for i := 0; i < M/2; i++ {
		p := rand.Intn(N)
		take2(p)
	}

	Q := rand.Intn(10) + 1
	A := make([]int, Q)
	B := make([]int, Q)
	for i := 0; i < Q; i++ {
		A[i] = rand.Intn(N) + 1
		B[i] = rand.Intn(N) + 1
		if A[i] == B[i] {
			B[i] = (B[i] + 1) % N
			if B[i] == 0 {
				B[i] = N
			}
		}
	}

	res1 := solve1(N, M, T, P, Q, A, B)
	fmt.Println(res1)
}

func solve1(N, M int, T, P []int, Q int, A, B []int) []int {
	times := make([][]int, N)
	for i := 0; i < M; i++ {
		times[P[i]] = append(times[P[i]], T[i])
	}

	intervals := make([][][2]int, N)
	for i := 0; i < N; i++ {
		curTimes := times[i]
		curIntervals := make([][2]int, 0, len(curTimes)/2)
		for j := 0; j < len(curTimes); j += 2 {
			curIntervals = append(curIntervals, [2]int{curTimes[j], curTimes[j+1]})
		}
		sort.Slice(curIntervals, func(i, j int) bool { return curIntervals[i][0] < curIntervals[j][0] })
		intervals[i] = curIntervals
	}

	const threshold int = 2 // sqrt(nm/q)
	isBig := func(id int) bool { return len(intervals[id]) > threshold }
	var big []int
	for i := 0; i < N; i++ {
		if isBig(i) {
			big = append(big, i)
		}
	}
	bigRes := make(map[int][]int, len(big))
	for i := 0; i < len(big); i++ {
		bigRes[big[i]] = make([]int, N)
	}

	// !扫描线+前缀和处理每个大区间与其他区间的交集范围
	type event = struct {
		id   int // in: id, out: -id-1
		time int
	}
	events := make([]event, 0, M)
	for i := 0; i < N; i++ {
		for j := 0; j < len(intervals[i]); j++ {
			events = append(events, event{id: i, time: intervals[i][j][0]})
			events = append(events, event{id: ^i, time: intervals[i][j][1]})
		}
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].time != events[j].time {
			return events[i].time < events[j].time
		}
		return events[i].id > events[j].id // in first
	})

	dp := make([]int, N)
	for bigId := 0; bigId < N; bigId++ {
		if !isBig(bigId) {
			continue
		}
		for i := range dp {
			dp[i] = 0
		}
		pre, curSum, inside := 0, 0, 0
		for _, e := range events {
			id, time := e.id, e.time
			if id < 0 {
				id = ^id
			}
			if id == bigId {
				curSum += inside * (time - pre)
				pre = time
				inside ^= 1
			}
			dp[id] = curSum + inside*(time-pre) - dp[id]
		}
		for i := 0; i < N; i++ {
			bigRes[bigId] = append(bigRes[bigId], dp[i])
		}
	}

	res := make([]int, Q)
	for i := 0; i < Q; i++ {
		a, b := A[i], B[i]
		a--
		b--
		if !isBig(a) && !isBig(b) {
			res[i] = IntervalsIntersection(intervals[a], intervals[b])
		} else {
			if !isBig(a) {
				a, b = b, a
			}
			res[i] = bigRes[a][b]
		}
	}

	return res
}

// 区间交集长度.
func IntervalsIntersection(intervals1, intervals2 [][2]int) int {
	n1, n2 := len(intervals1), len(intervals2)
	res := 0
	left, right := 0, 0
	for left < n1 && right < n2 {
		s1, e1, s2, e2 := intervals1[left][0], intervals1[left][1], intervals2[right][0], intervals2[right][1]
		if (s1 <= e2 && e2 <= e1) || (s2 <= e1 && e1 <= e2) {
			res += min(e1, e2) - max(s1, s2)
		}
		if e1 < e2 {
			left++
		} else {
			right++
		}
	}
	return res
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

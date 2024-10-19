package main

import (
	"bufio"
	"fmt"
	stdio "io"
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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, Q := int32(io.NextInt()), int32(io.NextInt())
	H, T := make([]byte, Q), make([]int32, Q)
	for i := int32(0); i < Q; i++ {
		s, t := io.Text(), io.NextInt()
		H[i] = s[0]
		T[i] = int32(t - 1)
	}

	distLeft := func(from_, to int32) int32 {
		if from_ >= to {
			return from_ - to
		}
		return from_ + N - to
	}

	distRight := func(from_, to int32) int32 {
		if to >= from_ {
			return to - from_
		}
		return to + N - from_
	}

	onPathLeft := func(from_, to, x int32) bool {
		if from_ == to {
			return false
		}
		if x < to {
			x += N
		}
		if from_ < to {
			from_ += N
		}
		return to <= x && x <= from_
	}

	onPathRight := func(from_, to, x int32) bool {
		if from_ == to {
			return false
		}
		if from_ > to {
			to += N
		}
		if from_ > x {
			x += N
		}
		return from_ <= x && x <= to
	}

	moveLeft := func(cur, to, other int32) (int32, int32) {
		if !onPathLeft(cur, to, other) {
			return distLeft(cur, to), other
		}
		otherTo := to - 1
		if to == 0 {
			otherTo = N - 1
		}
		return distLeft(cur, to) + distLeft(other, otherTo), otherTo
	}

	moveRight := func(cur, to, other int32) (int32, int32) {
		if !onPathRight(cur, to, other) {
			return distRight(cur, to), other
		}
		otherTo := to + 1
		if to == N-1 {
			otherTo = 0
		}
		return distRight(cur, to) + distRight(other, otherTo), otherTo
	}

	memo := make(map[int]int32)

	const INF32 int32 = 1e9 + 10

	var dfs func(int32, int32, int32) int32
	dfs = func(index, posL, posR int32) int32 {
		if index == Q {
			return 0
		}
		hash_ := int(index)<<24 | int(posL)<<12 | int(posR)
		if res, ok := memo[hash_]; ok {
			return res
		}
		to := T[index]
		res := INF32
		if H[index] == 'L' {
			d1, r1 := moveLeft(posL, to, posR)
			res = min32(res, d1+dfs(index+1, to, r1))
			d2, r2 := moveRight(posL, to, posR)
			res = min32(res, d2+dfs(index+1, to, r2))
		} else {
			d1, l1 := moveLeft(posR, to, posL)
			res = min32(res, d1+dfs(index+1, l1, to))
			d2, l2 := moveRight(posR, to, posL)
			res = min32(res, d2+dfs(index+1, l2, to))
		}
		memo[hash_] = res
		return res
	}

	res := dfs(0, 0, 1)

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

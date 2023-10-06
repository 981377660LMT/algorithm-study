// https://www.luogu.com.cn/problem/P6578
// P6578 [Ynoi2019]魔法少女网站 (第十分块)
//

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
	"strconv"
)

// 1 pos x: 将 pos 位置的值修改为 x
// 2 start end x: 查询区间 [start, end) 中有多少子数组最大值不超过 x.
//
// 不带修怎么做? -> 全部元素小于给定数的极长连续段, 大于等于 x 的数为1, 小于 x 的数为0.
func 魔法少女网站(nums []int, operations [][4]int) []int {

}

// 给定一个数组nums和一些查询(start,end,x)，对每个查询回答区间[start,end)内有多少个子数组最大值不超过x.
// nums.length<=1e5,nums[i]<=1e5,查询个数<=1e5
// 线段树或者分块
func 魔法少女网站无修改版本(nums []int, queries [][3]int) []int {
	type queryWithId struct{ start, end, x, id int }
	qs := make([]queryWithId, len(queries))
	for i, q := range queries {
		qs[i] = queryWithId{q[0], q[1], q[2], i}
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })

}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(nums []int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	n := len(nums)
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}

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

	n, q := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := range nums {
		nums[i] = io.NextInt()
	}

	operations := make([][4]int, q)
	for i := range operations {
		op := io.NextInt()
		if op == 1 {
			pos, x := io.NextInt(), io.NextInt()
			pos--
			operations[i] = [4]int{1, pos, x, 0}
		} else {
			start, end, x := io.NextInt(), io.NextInt(), io.NextInt()
			start--
			operations[i] = [4]int{2, start, end, x}
		}
	}

	res := 魔法少女网站(nums, operations)
	for _, v := range res {
		io.Println(v)
	}
}

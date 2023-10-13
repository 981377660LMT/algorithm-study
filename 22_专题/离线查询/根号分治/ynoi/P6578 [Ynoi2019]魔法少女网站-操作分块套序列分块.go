// https://www.luogu.com.cn/problem/P6578
// P6578 [Ynoi2019]魔法少女网站 (第十分块)
//

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

// 1 pos x: 将 pos 位置的值修改为 x
// 2 start end x: 查询区间 [start, end) 中有多少子数组最大值不超过 x.
//
// 不带修怎么做? ->
// 将所有查询按照x从小到大排序,然后从小到大依次处理，即维护一个01数组.
// 长为len的极长连续段的贡献为len*(len+1)/2.
// 用线段树维护.
// 带修 ->
// 线段树不好维护,用分块维护。
// 看一下分块不带修改怎么做。
// !对于一个询问，我们只需要从左到右合并当前阈值下 0,1 对应的信息就好了，但复杂度过大，考虑序列分块，从左到右逐块处理，维护每个询问合并到当前块的信息。
// 如果遇到散块，就直接暴力合并信息。否则考虑对于一个大小O(sqrt(n))的块也只有O(sqrt(n))种本质不用的x。
// 带修改需要操作分块。
// !https://www.luogu.com.cn/blog/ryoku/solution-p6578
type Node struct {
	size      int // 区间长度
	preOne    int // 	前缀连续1的个数
	sufOne    int // 后缀连续1的个数
	pairCount int // !区间贡献
}

func 魔法少女网站(nums []int, operations [][4]int) []int {

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

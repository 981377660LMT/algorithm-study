package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"os"
	"sort"
	"strconv"
)

// https://www.luogu.com.cn/problem/solution/P3604
// P3604 美好的每一天-求区间可重排为回文串的个数(区间回文)
// 字符串只包含小写字母.
// n,q<=6e4
// !等价于每次查询区间中有几个二元组的的异或值在二进制下1的个数<=1.
// 异或前缀和.
func Solve(s string, ranges [][2]int) []int {
	n, q := len(s), len(ranges)
	M := NewMoAlgo(n, q)
	for _, query := range ranges {
		start, end := query[0], query[1]
		M.AddQuery(start, end)
	}

	preXor := make([]int, n+1)
	for i := 0; i < n; i++ {
		preXor[i+1] = preXor[i] ^ (1 << (s[i] - 'a'))
	}

	counter := [1 << 26]int{}
	counter[0] = 1

	res := make([]int, q)
	count := 0
	M.Run(
		// add
		func(index, delta int) {
			// 窗口向右扩张
			if delta == 1 {
				index++
			}
			cur := preXor[index]
			count += counter[cur]
			counter[cur]++
			for i := 0; i < 26; i++ {
				count += counter[cur^(1<<i)]
			}
		},
		// remove
		func(index, delta int) {
			// 窗口向左收缩
			if delta == -1 {
				index++
			}
			cur := preXor[index]
			counter[cur]--
			count -= counter[cur]
			for i := 0; i < 26; i++ {
				count -= counter[cur^(1<<i)]
			}
		},
		// query
		func(qid int) {
			res[qid] = count
		},
	)
	return res
}

type MoAlgo struct {
	queryOrder int
	chunkSize  int
	buckets    [][]query
}

type query struct{ qi, left, right int }

func NewMoAlgo(n, q int) *MoAlgo {
	chunkSize := max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	buckets := make([][]query, n/chunkSize+1)
	return &MoAlgo{chunkSize: chunkSize, buckets: buckets}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//
//	0 <= left <= right <= n
func (mo *MoAlgo) AddQuery(left, right int) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right})
	mo.queryOrder++
}

// 返回每个查询的结果.
//
//	add: 将数据添加到窗口. delta: 1 表示向右移动，-1 表示向左移动.
//	remove: 将数据从窗口移除. delta: 1 表示向右移动，-1 表示向左移动.
//	query: 查询窗口内的数据.
func (mo *MoAlgo) Run(
	add func(index, delta int),
	remove func(index, delta int),
	query func(qid int),
) {
	left, right := 0, 0

	for i, bucket := range mo.buckets {
		if i&1 == 1 {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right < bucket[j].right })
		} else {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right > bucket[j].right })
		}

		for _, q := range bucket {
			// !窗口扩张
			for left > q.left {
				left--
				add(left, -1)
			}
			for right < q.right {
				add(right, 1)
				right++
			}

			// !窗口收缩
			for left < q.left {
				remove(left, 1)
				left++
			}
			for right > q.right {
				right--
				remove(right, -1)
			}

			query(q.qi)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

	_, q := io.NextInt(), io.NextInt()
	s := io.Text()
	ranges := make([][2]int, q)
	for i := 0; i < q; i++ {
		ranges[i] = [2]int{io.NextInt() - 1, io.NextInt()}
	}

	res := Solve(s, ranges)
	for _, v := range res {
		io.Println(v)
	}
}

// 给定一个数组和一些查询区间
// 求区间内的三元组数目(i,j,k)使得i<j<k且a[i]=a[j]=a[k]

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
	mo := NewMoAlgo(n, q)
	for i := 0; i < q; i++ {
		l, r := io.NextInt()-1, io.NextInt()
		mo.AddQuery(l, r)
	}

	//
	//
	comb3 := func(n int) int {
		return n * (n - 1) * (n - 2) / 6
	}

	res := make([]int, q)
	counter := [2e5 + 10]int{}
	count := 0 // 相等的三元组的个数
	mo.Run(
		func(index, _ int) {
			counter[nums[index]] += 1
			count += comb3(counter[nums[index]]) - comb3(counter[nums[index]]-1)
		},
		func(index, _ int) {
			counter[nums[index]] -= 1
			count += comb3(counter[nums[index]]) - comb3(counter[nums[index]]+1)
		},
		func(qid int) {
			res[qid] = count
		},
	)
	for _, v := range res {
		io.Println(v)
	}
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
//  0 <= left <= right <= n
func (mo *MoAlgo) AddQuery(left, right int) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right})
	mo.queryOrder++
}

// 返回每个查询的结果.
//  add: 将数据添加到窗口. delta: 1 表示向右移动，-1 表示向左移动.
//  remove: 将数据从窗口移除. delta: 1 表示向右移动，-1 表示向左移动.
//  query: 查询窗口内的数据.
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

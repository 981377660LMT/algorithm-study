package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
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

type PointSetModSum struct {
	nums          []int
	stepThreshold int
	// dp[step][start] 表示步长为step,起点为start的所有元素的和.
	// `dp[step][start] = dp[step][start+step] + nums[start]`.
	dp [][]int
}

// stepThreshold: 步长阈值,当步长小于等于该值时,使用dp数组预处理答案,否则直接遍历.
// 预处理时间空间复杂度均为`O(n*stepThreshold)`.
// 单次遍历时间复杂度为`O(n/stepThreshold)`.
// 取50较为合适.
func NewPointSetModSum(nums []int, stepThreshold int) *PointSetModSum {
	n := len(nums)
	dp := make([][]int, stepThreshold)
	for step := 1; step <= stepThreshold; step++ {
		curSum := make([]int, n+1)
		for start := n - 1; start >= 0; start-- {
			curSum[start] = curSum[min(n, start+step)] + nums[start]
		}
		dp[step-1] = curSum
	}
	return &PointSetModSum{nums: nums, stepThreshold: stepThreshold, dp: dp}
}

func (pss *PointSetModSum) Set(index, value int) {
	if index < 0 || index >= len(pss.nums) {
		return
	}
	pre := pss.nums[index]
	if pre == value {
		return
	}
	pss.nums[index] = value
	delta := value - pre
	for step := 1; step <= pss.stepThreshold; step++ {
		pss.dp[step-1][index%step] += delta
	}
}

// 查询 sum(nums[start::step]).
func (pss *PointSetModSum) Query(start, step int) int {
	if start < 0 {
		start = 0
	}
	if step <= pss.stepThreshold {
		return pss.dp[step-1][start]
	}
	sum := 0
	for i := start; i < len(pss.nums); i += step {
		sum += pss.nums[i]
	}
	return sum
}

func (pss *PointSetModSum) String() string {
	return fmt.Sprintf("PointSetModSum{nums: %v}", pss.nums)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 一列に並んだ
// N 個のマス
// 1,2,…,N と長さ
// N の数列
// A=(A
// 1
// ​
//  ,A
// 2
// ​
//  ,…,A
// N
// ​
//  ) があります。
// 最初、マス
// 1 は黒く、他の
// N−1 個のマスは白く塗られており、
// 1 つのコマがマス
// 1 に置かれています。

// 以下の操作を
// 0 回以上好きな回数繰り返します。

// コマがマス
// i にあるとき、ある正整数
// x を決めてコマをマス
// i+A
// i
// ​
//
//	×x に移動させる。
//
// 但し、
// i+A
// i
// ​
//
//	×x>N となるような移動はできません。
//
// その後、マス
// i+A
// i
// ​
//
//	×x を黒く塗る。
//
// 操作を終えた時点で黒く塗られたマスの集合として考えられるものの数を
// 998244353 で割った余りを求めてください。

const MOD int = 998244353

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}

	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = 1
	}
	S := NewPointSetModSum(values, int(math.Sqrt(float64(n))))
	res := 0
	for i := n - 1; i >= 0; i-- {
		res += S.Query(i, nums[i])
		res %= MOD
		S.Set(i, res)
	}
	io.Println(res)
}

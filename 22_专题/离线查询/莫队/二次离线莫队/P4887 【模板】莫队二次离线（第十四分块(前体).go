// P4887 【模板】莫队二次离线（第十四分块(前体)
package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"math/bits"
	"os"
	"sort"
	"strconv"
)

// !每次查询一个区间中有几个二元组的的异或值在二进制下有k个1
// 1<=n,q<=1e5 0<=ai,k<=16384 (2^14)
// 莫队二次离线要求O(1)查询,O(sqrt(n))修改=>注意到C(14,k)最多C(14,7)=3432种情况.转移时可以枚举.
// !预处理出所有有 k 个1的数组ok，利用异或性质即可得到答案(a^b=c => a^c=b, b^c=a)
// 每一次插入一个数，就把它与ok数组中所有的数异或起来的值放进一个桶里，然后查询就直接查当前数在桶中出现了几次
func Solve(nums []int, ranges [][2]int, k int) []int {
	n, q := len(nums), len(ranges)
	M := NewMoOfflineAgain(n, q, -1)
	for _, query := range ranges {
		start, end := query[0], query[1]
		M.AddQuery(start, end)
	}

	ok := []int{} // 有k个1的数
	for i := 0; i <= 1<<14; i++ {
		if bits.OnesCount(uint(i)) == k {
			ok = append(ok, i)
		}
	}
	xorCounter := make([]int, 1<<15)
	res := M.Run(
		func(index int) {
			cur := nums[index]
			for _, v := range ok {
				xorCounter[v^cur]++
			}
		},
		func(index int) AbelianGroup {
			cur := nums[index]
			return xorCounter[cur]
		},
		func(index int) AbelianGroup {
			cur := nums[index]
			return xorCounter[cur]
		},
	)
	return res
}

// 可交换群(commutative group).
type AbelianGroup = int

func e() AbelianGroup                   { return 0 }
func op(a, b AbelianGroup) AbelianGroup { return a + b }
func inv(a AbelianGroup) AbelianGroup   { return -a }

type MoOfflineAgain struct {
	n           int
	q           int
	blockSize   int
	queryBlocks [][]query
	queryOrder  int
}

type query struct{ qi, left, right int }

// n: 数组长度, q: 查询个数, blockSize: 块大小,-1 表示使用默认值.
func NewMoOfflineAgain(n, q int, blockSize int) *MoOfflineAgain {
	if blockSize == -1 {
		blockSize = max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	}
	queryBlocks := make([][]query, n/blockSize+1)
	return &MoOfflineAgain{n: n, q: q, blockSize: blockSize, queryBlocks: queryBlocks}
}

// 添加一个查询，查询范围为`左闭右开区间` [start, end).
//
//	0 <= start <= end <= n
func (mo *MoOfflineAgain) AddQuery(start, end int) {
	bid := start / mo.blockSize
	mo.queryBlocks[bid] = append(mo.queryBlocks[bid], query{qi: mo.queryOrder, left: start, right: end})
	mo.queryOrder++
}

// add: 将`A[index]`加入窗口中.
// queryLeft: 窗口最左侧的`A[index]`对答案的贡献.
// queryRight: 窗口最右侧的`A[index]`对答案的贡献.
// `add` 操作次数为`O(n)`，`query` 操作次数为`O(nsqrt(n))`.
func (mo *MoOfflineAgain) Run(
	add func(index int),
	queryLeft func(index int) AbelianGroup,
	queryRight func(index int) AbelianGroup,
) []AbelianGroup {
	n, q, blocks := mo.n, mo.q, mo.queryBlocks
	type event struct{ qi, start, end, kind int }
	eventGroups := make([][]event, n+1)

	left, right := 0, 0
	for i, block := range blocks {
		if i&1 == 1 {
			sort.Slice(block, func(i, j int) bool { return block[i].right < block[j].right })
		} else {
			sort.Slice(block, func(i, j int) bool { return block[i].right > block[j].right })
		}

		for _, q := range block {
			qi, ql, qr := q.qi, q.left, q.right
			if ql < left {
				eventGroups[right] = append(eventGroups[right], event{qi: qi, start: ql, end: left, kind: 2})
				left = ql
			}
			if right < qr {
				eventGroups[left] = append(eventGroups[left], event{qi: qi, start: right, end: qr, kind: 1})
				right = qr
			}
			if left < ql {
				eventGroups[right] = append(eventGroups[right], event{qi: qi, start: left, end: ql, kind: 0})
				left = ql
			}
			if qr < right {
				eventGroups[left] = append(eventGroups[left], event{qi: qi, start: qr, end: right, kind: 3})
				right = qr
			}
		}

	}

	rightSum, leftSum := make([]AbelianGroup, n+1), make([]AbelianGroup, n+1)
	rightSum[0], leftSum[0] = e(), e()
	res := make([]AbelianGroup, q)
	for i := range res {
		res[i] = e()
	}
	for i := 0; i <= n; i++ {
		events := eventGroups[i]
		for _, event := range events {
			qi, start, end, kind := event.qi, event.start, event.end, event.kind
			sum := e()
			if kind&1 != 0 {
				for j := start; j < end; j++ {
					sum = op(sum, queryRight(j))
				}
			} else {
				for j := start; j < end; j++ {
					sum = op(sum, queryLeft(j))
				}
			}

			if kind&2 != 0 {
				res[qi] = op(res[qi], inv(sum))
			} else {
				res[qi] = op(res[qi], sum)
			}

		}

		if i < n {
			rightSum[i+1] = op(rightSum[i], queryRight(i))
			add(i)
			leftSum[i+1] = op(leftSum[i], queryLeft(i))
		}

	}

	curSum := e()
	for _, block := range blocks {
		for j := range block {
			qi, ql, qr := block[j].qi, block[j].left, block[j].right
			curSum = op(curSum, res[qi])
			res[qi] = op(op(leftSum[ql], rightSum[qr]), inv(curSum))
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

	n, q, k := io.NextInt(), io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}
	ranges := make([][2]int, q)
	for i := 0; i < q; i++ {
		ranges[i] = [2]int{io.NextInt() - 1, io.NextInt()}
	}

	res := Solve(nums, ranges, k)
	for _, v := range res {
		io.Println(v)
	}
}

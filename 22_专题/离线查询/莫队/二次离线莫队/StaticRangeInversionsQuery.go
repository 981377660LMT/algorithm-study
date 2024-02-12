// 静态区间逆序对
// 二次离线莫队

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 静态区间逆序对-离线.
// 时间复杂度O(nsqrt(n)),空间复杂度O(n).
// https://judge.yosupo.jp/problem/static_range_inversions_query
func StaticRangeInversionsQuery(nums []int, ranges [][2]int) []int {
	n, q := len(nums), len(ranges)
	rank, newNums := discretize(nums)
	bit := _newBITRangeBlockFastQuery(len(rank))
	M := NewMoOfflineAgain(n, q, -1)
	for _, query := range ranges {
		start, end := query[0], query[1]
		M.AddQuery(start, end)
	}
	res := M.Run(
		func(index int) {
			bit.Add(newNums[index], 1)
		},
		func(index int) AbelianGroup {
			return bit.QueryRange(0, newNums[index])
		},
		func(index int) AbelianGroup {
			return bit.QueryRange(newNums[index]+1, len(rank))
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

// add: 将`A[index]`加入数据结构中.
// queryLeft: 查询`A[index]`左侧(包含index本身)的贡献之和.
// queryRight: 查询`A[index]`右侧(包含index本身)的贡献之和.
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

// 基于分块实现的`树状数组`.
// `O(sqrt(n))`单点加，`O(1)`查询区间和.
// 一般配合莫队算法使用.
type _BITRangeBlockFastQuery struct {
	_n           int
	_belong      []int
	_blockStart  []int
	_blockEnd    []int
	_blockCount  int
	_partPreSum  []int
	_blockPreSum []int
}

func _newBITRangeBlockFastQuery(lengthOrArray interface{}) *_BITRangeBlockFastQuery {
	var n int
	var isArray bool
	if length, ok := lengthOrArray.(int); ok {
		n = length
	} else {
		n = len(lengthOrArray.([]int))
		isArray = true
	}
	blockSize := int(math.Sqrt(float64(n)) + 1)
	blockCount := 1 + (n / blockSize)
	belong := make([]int, n)
	for i := range belong {
		belong[i] = i / blockSize
	}
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	for i := range blockStart {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	partPreSum := make([]int, n)
	blockPreSum := make([]int, blockCount)
	res := &_BITRangeBlockFastQuery{
		_n:           n,
		_belong:      belong,
		_blockStart:  blockStart,
		_blockEnd:    blockEnd,
		_blockCount:  blockCount,
		_partPreSum:  partPreSum,
		_blockPreSum: blockPreSum,
	}
	if isArray {
		res.Build(lengthOrArray.([]int))
	}
	return res
}

func (b *_BITRangeBlockFastQuery) Add(index int, delta int) {
	if index < 0 || index >= b._n {
		panic("index out of range")
	}
	bid := b._belong[index]
	for i := index; i < b._blockEnd[bid]; i++ {
		b._partPreSum[i] += delta
	}
	for id := bid + 1; id < b._blockCount; id++ {
		b._blockPreSum[id] += delta
	}
}

func (b *_BITRangeBlockFastQuery) QueryRange(start int, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return 0
	}
	return b._query(end) - b._query(start)
}

func (b *_BITRangeBlockFastQuery) Build(arr []int) {
	if len(arr) != b._n {
		panic("array length mismatch n")
	}
	curBlockSum := 0
	for bid := 0; bid < b._blockCount; bid++ {
		curPartSum := 0
		for i := b._blockStart[bid]; i < b._blockEnd[bid]; i++ {
			curPartSum += arr[i]
			b._partPreSum[i] = curPartSum
		}
		b._blockPreSum[bid] = curBlockSum
		curBlockSum += curPartSum
	}
}

func (b *_BITRangeBlockFastQuery) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := range b._partPreSum {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	sb.WriteString("}")
	return sb.String()
}

func (b *_BITRangeBlockFastQuery) _query(end int) int {
	if end <= 0 {
		return 0
	}
	return b._partPreSum[end-1] + b._blockPreSum[b._belong[end-1]]
}

// (紧)离散化.
//
//	rank: 给定一个在 nums 中的值,返回它的排名(0~len(rank)-1).
//	newNums: 离散化后的数组.
func discretize(nums []int) (rank map[int]int, newNums []int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	rank = make(map[int]int, len(allNums))
	for i, v := range allNums {
		rank[v] = i
	}
	newNums = make([]int, len(nums))
	for i, v := range nums {
		newNums[i] = rank[v]
	}
	return
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

	n, q := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}
	ranges := make([][2]int, q)
	for i := 0; i < q; i++ {
		ranges[i] = [2]int{io.NextInt(), io.NextInt()}
	}

	res := StaticRangeInversionsQuery(nums, ranges)
	for _, v := range res {
		io.Println(v)
	}
}

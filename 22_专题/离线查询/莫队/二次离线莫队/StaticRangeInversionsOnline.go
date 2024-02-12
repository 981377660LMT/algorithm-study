// 静态区间逆序对，分块在线算法.

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

// https://leetcode.cn/problems/shu-zu-zhong-de-ni-xu-dui-lcof/
func reversePairs(record []int) int {
	R := NewStaticRangeInversionsOnline(record, -1)
	return R.Query(0, len(record))
}

// 在线求解静态区间逆序对问题。
// 时间复杂度O((n+q)sqrt(n)),空间复杂度O(nsqrt(n)).
// https://stackoverflow.com/questions/21763392/counting-inversions-in-ranges
type StaticRangeInversionsOnline struct {
	n            int32
	blockSize    int32
	blockCount   int32
	values       []int32
	sortedValues [][2]int32
	presuf       [][]int32
	sufG, preH   []int32
	preSum2d     [][]int
}

// blockSize: 分块大小.-1表示使用默认值`2*sqrt(n)`.
func NewStaticRangeInversionsOnline(nums []int, blockSize int) *StaticRangeInversionsOnline {
	res := &StaticRangeInversionsOnline{}
	n := int32(len(nums))
	blockSize32 := int32(blockSize)
	if blockSize == -1 {
		blockSize32 = 2 * (int32(math.Ceil(math.Sqrt(float64(n)))) + 1)
	}
	blockCount := (n + blockSize32 - 1) / blockSize32

	dict := append(nums[:0:0], nums...)
	sort.Ints(dict)
	uniqueSortedInplace(&dict)
	d := int32(len(dict))
	values := make([]int32, 0, n)
	sortedValues := make([][2]int32, 0, n)
	for _, v := range nums {
		pos := int32(sort.SearchInts(dict, v))
		values = append(values, pos)
		sortedValues = append(sortedValues, [2]int32{pos, int32(len(values)) - 1})
	}

	presuf := make([][]int32, blockCount)
	for i := range presuf {
		presuf[i] = make([]int32, n)
	}
	sufG := make([]int32, n)
	preH := make([]int32, n)

	for bid := int32(0); bid < blockCount; bid++ {
		start := bid * blockSize32
		end := min32(start+blockSize32, n)
		tmp := sortedValues[start:end]
		sort.Slice(tmp, func(i, j int) bool {
			if tmp[i][0] != tmp[j][0] {
				return tmp[i][0] < tmp[j][0]
			}
			return tmp[i][1] < tmp[j][1]
		})
		counter := make([]int32, d+1)
		for i := start; i < end; i++ {
			counter[values[i]+1]++
		}
		for i := int32(0); i < d; i++ {
			counter[i+1] += counter[i]
		}
		for b := int32(0); b < bid; b++ {
			for i := (b+1)*blockSize32 - 1; i >= b*blockSize32; i-- {
				presuf[bid][i] = presuf[bid][i+1] + counter[values[i]]
			}
		}
		for b := bid + 1; b < blockSize32; b++ {
			for i := b * blockSize32; i < min32((b+1)*blockSize32, n); i++ {
				presuf[bid][i] = counter[d] - counter[values[i]+1]
				if i != b*blockSize32 {
					presuf[bid][i] += presuf[bid][i-1]
				}
			}
		}
		for i := start + 1; i < end; i++ {
			inv := int32(0)
			for j := start; j < i; j++ {
				if values[j] > values[i] {
					inv++
				}
			}
			preH[i] = preH[i-1] + inv
		}
		for i := end - 2; i >= start; i-- {
			inv := int32(0)
			for j := i + 1; j < end; j++ {
				if values[j] < values[i] {
					inv++
				}
			}
			sufG[i] = sufG[i+1] + inv
		}
	}

	preSum2d := make([][]int, blockCount)
	for i := range preSum2d {
		preSum2d[i] = make([]int, blockCount)
	}
	for i := blockCount - 1; i >= 0; i-- {
		preSum2d[i][i] = int(sufG[i*blockSize32])
		for j := i + 1; j < blockCount; j++ {
			preSum2d[i][j] = preSum2d[i][j-1] + preSum2d[i+1][j] - preSum2d[i+1][j-1] + int(presuf[j][i*blockSize32])
		}
	}

	res.n = n
	res.blockSize = blockSize32
	res.blockCount = blockCount
	res.values = values
	res.sortedValues = sortedValues
	res.presuf = presuf
	res.sufG = sufG
	res.preH = preH
	res.preSum2d = preSum2d
	return res
}

func (sr *StaticRangeInversionsOnline) Query(start, end int) (res int) {
	l, r := int32(start), int32(end)
	n, bs, bc := sr.n, sr.blockSize, sr.blockCount
	sortedValues := sr.sortedValues
	if l < 0 {
		l = 0
	}
	if r > n {
		r = n
	}
	if l >= r {
		return 0
	}

	if b := l / bs; b == (r-1)/bs {
		res += int(sr.preH[r-1])
		if l%bs != 0 {
			res -= int(sr.preH[l-1])
		}
		lessCount := 0
		for p, q := b*bs, min32((b+1)*bs, n); p < q; p++ {
			tmp := sortedValues[p][1]
			if tmp >= l && tmp < r {
				lessCount++
			}
			if tmp < l {
				res -= lessCount
			}
		}
		return
	}

	lb := (l + bs - 1) / bs
	var rb int32
	if r == n {
		rb = bc - 1
	} else {
		rb = r/bs - 1
	}

	res += sr.preSum2d[lb][rb]
	if bs*lb > l {
		res += int(sr.sufG[l])
		for b := lb; b <= rb; b++ {
			res += int(sr.presuf[b][l])
		}
	}
	if bs*(rb+1) < r {
		res += int(sr.preH[r-1])
		for b := lb; b <= rb; b++ {
			res += int(sr.presuf[b][r-1])
		}
	}
	lessCount := 0
	j := (rb + 1) * bs
	for p, q := max32(0, (lb-1)*bs), lb*bs; p < q; p++ {
		if sortedValues[p][1] >= l {
			for j < min32(n, (rb+2)*bs) && (sortedValues[j][1] >= r || sortedValues[j][0] < sortedValues[p][0]) {
				if sortedValues[j][1] < r {
					lessCount++
				}
				j++
			}
			res += lessCount
		}
	}
	return
}

func uniqueSortedInplace(sorted *[]int) {
	if len(*sorted) == 0 {
		return
	}
	tmp := *sorted
	slow := 0
	for fast := 0; fast < len(tmp); fast++ {
		if tmp[fast] != tmp[slow] {
			slow++
			tmp[slow] = tmp[fast]
		}
	}
	*sorted = tmp[:slow+1]
}

// 静态区间逆序对-离线.
// 时间复杂度O(nsqrt(n)),空间复杂度O(n).
// https://judge.yosupo.jp/problem/static_range_inversions_query
func StaticRangeInversionsQuery(nums []int, ranges [][2]int) []int {
	S := NewStaticRangeInversionsOnline(nums, -1)
	res := make([]int, len(ranges))
	for i, r := range ranges {
		res[i] = S.Query(r[0], r[1])
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

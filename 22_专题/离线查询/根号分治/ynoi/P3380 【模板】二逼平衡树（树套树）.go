// https://www.luogu.com.cn/problem/P3380

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

const INF int = 2147483647

// 1 start end num: 查询num在[start, end)内的排名.
// 2 start end k: 查询[start, end)内排名为k(k>=0)的数.
// 3 pos num: 将pos位置的数修改为num.
// 4 start end num: 查询[start, end)内num的严格前驱.不存在输出-2147483647.
// 5 start end num: 查询[start, end)内num的严格后继.不存在输出2147483647.
func Solve(nums []int, operations [][4]int) []int {
	allNums := append(nums[:0:0], nums...)
	for i := range operations {
		kind := operations[i][0]
		if kind == 1 || kind == 4 || kind == 5 {
			allNums = append(allNums, operations[i][3])
		} else if kind == 3 {
			allNums = append(allNums, operations[i][2])
		}
	}
	set := make(map[int]struct{})
	for _, v := range allNums {
		set[v] = struct{}{}
	}
	sorted := make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	valueToId := make(map[int]int, len(sorted))
	idToValue := make([]int, len(sorted))
	for i, v := range sorted {
		valueToId[v] = i
		idToValue[i] = v
	}

	for i := range nums {
		nums[i] = valueToId[nums[i]]
	}

	blockSize1 := int(math.Sqrt(float64(len(nums)))+1) * 4 // !增加块的大小以减少块的数量，减少内存占用.
	block1 := UseBlock(len(nums), blockSize1)
	belong1, blockStart1, blockEnd1, blockCount1 := block1.belong, block1.blockStart, block1.blockEnd, block1.blockCount
	block2 := UseBlock(len(sorted), int(math.Sqrt(float64(len(sorted)))+1))
	belong2, blockStart2, blockEnd2, blockCount2 := block2.belong, block2.blockStart, block2.blockEnd, block2.blockCount
	sum1 := make([][]int, blockCount1+1)
	for i := range sum1 {
		sum1[i] = make([]int, blockCount2)
	}
	sum2 := make([][]int, blockCount1+1)
	for i := range sum2 {
		sum2[i] = make([]int, len(sorted))
	}
	fragmentCounter1 := make([]int, blockCount2) // 查询时临时保存散块内出现在第i个值域块中的数的个数.
	fragmentCounter2 := make([]int, len(sorted)) // 查询时临时保存散块内等于i的数的个数.

	// 初始化前缀和.
	init := func() {
		for bid := 0; bid < blockCount1; bid++ {
			copy(sum1[bid+1], sum1[bid])
			copy(sum2[bid+1], sum2[bid])
			s1, s2 := sum1[bid+1], sum2[bid+1]
			for i := blockStart1[bid]; i < blockEnd1[bid]; i++ {
				num := nums[i]
				vid := belong2[num]
				s1[vid]++
				s2[num]++
			}
		}
	}
	init()

	// 更新散块信息.
	updateFragment := func(start, end, delta int) {
		for i := start; i < end; i++ {
			num := nums[i]
			vid := belong2[num]
			fragmentCounter1[vid] += delta
			fragmentCounter2[num] += delta
		}
	}

	// 单点修改.
	pointSet := func(pos, target int) {
		if nums[pos] == target {
			return
		}
		preValue, curValue := nums[pos], target
		preVid, curVid := belong2[preValue], belong2[curValue]
		for bid := belong1[pos]; bid < blockCount1; bid++ {
			s1, s2 := sum1[bid+1], sum2[bid+1]
			s1[preVid]--
			s1[curVid]++
			s2[preValue]--
			s2[curValue]++
		}
		nums[pos] = target
	}

	// 排名，即严格小于target的元素个数(0-based).
	queryRank := func(start, end, target int) (res int) {
		bid1, bid2 := belong1[start], belong1[end-1]

		if bid1 == bid2 {
			for i := start; i < end; i++ {
				if nums[i] < target {
					res++
				}
			}
			return
		} else {
			for i := start; i < blockEnd1[bid1]; i++ {
				if nums[i] < target {
					res++
				}
			}
			for i := blockStart1[bid2]; i < end; i++ {
				if nums[i] < target {
					res++
				}
			}
			// 计算整块[bid1+1,bid2)内严格小于target的元素个数.
			// !这样的元素在值域的前缀的若干个整块+一个散块中.
			vid := belong2[target]
			for j := 0; j < vid; j++ {
				res += sum1[bid2][j] - sum1[bid1+1][j]
			}
			for j := blockStart2[vid]; j < target; j++ {
				res += sum2[bid2][j] - sum2[bid1+1][j]
			}
			return
		}
	}

	// 查询区间第k小(k>=0).
	queryKth := func(start, end, kth int) int {
		bid1, bid2 := belong1[start], belong1[end-1]

		if bid1 == bid2 {
			updateFragment(start, end, 1)
			todo := kth + 1
			for vid := 0; vid < blockCount2; vid++ {
				if tmp := todo - fragmentCounter1[vid]; tmp > 0 {
					todo = tmp
					continue
				}
				// 答案在这个值域中.
				for j := blockStart2[vid]; j < blockEnd2[vid]; j++ {
					todo -= fragmentCounter2[j]
					if todo <= 0 {
						updateFragment(start, end, -1)
						return j
					}
				}
			}
		} else {
			updateFragment(start, blockEnd1[bid1], 1)
			updateFragment(blockStart1[bid2], end, 1)
			todo := kth + 1
			for vid := 0; vid < blockCount2; vid++ {
				curCount := fragmentCounter1[vid] + sum1[bid2][vid] - sum1[bid1+1][vid]
				if tmp := todo - curCount; tmp > 0 {
					todo = tmp
					continue
				}
				for j := blockStart2[vid]; j < blockEnd2[vid]; j++ {
					curCount := fragmentCounter2[j] + sum2[bid2][j] - sum2[bid1+1][j]
					todo -= curCount
					if todo <= 0 {
						updateFragment(start, blockEnd1[bid1], -1)
						updateFragment(blockStart1[bid2], end, -1)
						return j
					}
				}
			}
		}

		panic("unreachable")
	}

	queryLower := func(start, end, target int) (res int, ok bool) {
		bid1, bid2 := belong1[start], belong1[end-1]
		vid := belong2[target]

		if bid1 == bid2 {
			updateFragment(start, end, 1)
			cand := target - 1
			vStart := blockStart2[vid]
			for cand >= vStart && fragmentCounter2[cand] == 0 {
				cand--
			}
			// 前驱在同一个值域块中.
			if cand >= vStart {
				updateFragment(start, end, -1)
				return cand, true
			}

			candVid := vid - 1
			for candVid >= 0 && fragmentCounter1[candVid] == 0 {
				candVid--
			}
			if candVid == -1 {
				updateFragment(start, end, -1)
				return
			}
			// 前驱在当前块candVid中.
			for j := blockEnd2[candVid] - 1; j >= blockStart2[candVid]; j-- {
				if fragmentCounter2[j] > 0 {
					updateFragment(start, end, -1)
					return j, true
				}
			}
			panic("unreachable")
		} else {
			updateFragment(start, blockEnd1[bid1], 1)
			updateFragment(blockStart1[bid2], end, 1)
			cand := target - 1
			vStart := blockStart2[vid]
			for cand >= vStart && (fragmentCounter2[cand]+sum2[bid2][cand]-sum2[bid1+1][cand]) == 0 {
				cand--
			}
			if cand >= vStart {
				updateFragment(start, blockEnd1[bid1], -1)
				updateFragment(blockStart1[bid2], end, -1)
				return cand, true
			}

			candVid := vid - 1
			for candVid >= 0 && (fragmentCounter1[candVid]+sum1[bid2][candVid]-sum1[bid1+1][candVid]) == 0 {
				candVid--
			}
			if candVid == -1 {
				updateFragment(start, blockEnd1[bid1], -1)
				updateFragment(blockStart1[bid2], end, -1)
				return
			}
			for j := blockEnd2[candVid] - 1; j >= blockStart2[candVid]; j-- {
				if fragmentCounter2[j]+sum2[bid2][j]-sum2[bid1+1][j] > 0 {
					updateFragment(start, blockEnd1[bid1], -1)
					updateFragment(blockStart1[bid2], end, -1)
					return j, true
				}
			}
			panic("unreachable")
		}
	}

	queryHigher := func(start, end, target int) (res int, ok bool) {
		bid1, bid2 := belong1[start], belong1[end-1]
		vid := belong2[target]

		if bid1 == bid2 {
			updateFragment(start, end, 1)
			cand := target + 1
			vEnd := blockEnd2[vid]
			for cand < vEnd && fragmentCounter2[cand] == 0 {
				cand++
			}
			if cand < vEnd {
				updateFragment(start, end, -1)
				return cand, true
			}

			candVid := vid + 1
			for candVid < blockCount2 && fragmentCounter1[candVid] == 0 {
				candVid++
			}
			if candVid == blockCount2 {
				updateFragment(start, end, -1)
				return
			}
			for j := blockStart2[candVid]; j < blockEnd2[candVid]; j++ {
				if fragmentCounter2[j] > 0 {
					updateFragment(start, end, -1)
					return j, true
				}
			}
			panic("unreachable")
		} else {
			updateFragment(start, blockEnd1[bid1], 1)
			updateFragment(blockStart1[bid2], end, 1)
			cand := target + 1
			vEnd := blockEnd2[vid]
			for cand < vEnd && (fragmentCounter2[cand]+sum2[bid2][cand]-sum2[bid1+1][cand]) == 0 {
				cand++
			}
			if cand < vEnd {
				updateFragment(start, blockEnd1[bid1], -1)
				updateFragment(blockStart1[bid2], end, -1)
				return cand, true
			}

			candVid := vid + 1
			for candVid < blockCount2 && (fragmentCounter1[candVid]+sum1[bid2][candVid]-sum1[bid1+1][candVid]) == 0 {
				candVid++
			}
			if candVid == blockCount2 {
				updateFragment(start, blockEnd1[bid1], -1)
				updateFragment(blockStart1[bid2], end, -1)
				return
			}
			for j := blockStart2[candVid]; j < blockEnd2[candVid]; j++ {
				if fragmentCounter2[j]+sum2[bid2][j]-sum2[bid1+1][j] > 0 {
					updateFragment(start, blockEnd1[bid1], -1)
					updateFragment(blockStart1[bid2], end, -1)
					return j, true
				}
			}
			panic("unreachable")
		}
	}

	res := make([]int, 0, len(operations))
	for _, op := range operations {
		kind := op[0]
		if kind == 1 {
			start, end, num := op[1], op[2], op[3]
			num = valueToId[num]
			rank := queryRank(start, end, num)
			res = append(res, rank+1)
		} else if kind == 2 {
			start, end, k := op[1], op[2], op[3]
			kth := queryKth(start, end, k)
			res = append(res, idToValue[kth])
		} else if kind == 3 {
			pos, num := op[1], op[2]
			num = valueToId[num]
			pointSet(pos, num)
		} else if kind == 4 {
			start, end, num := op[1], op[2], op[3]
			num = valueToId[num]
			lower, ok := queryLower(start, end, num)
			if !ok {
				res = append(res, -INF)
			} else {
				res = append(res, idToValue[lower])
			}
		} else if kind == 5 {
			start, end, num := op[1], op[2], op[3]
			num = valueToId[num]
			higher, ok := queryHigher(start, end, num)
			if !ok {
				res = append(res, INF)
			} else {
				res = append(res, idToValue[higher])
			}
		}
	}

	return res
}

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
			start, end, num := io.NextInt(), io.NextInt(), io.NextInt()
			start--
			operations[i] = [4]int{1, start, end, num}
		} else if op == 2 {
			start, end, k := io.NextInt(), io.NextInt(), io.NextInt()
			start--
			k--
			operations[i] = [4]int{2, start, end, k}
		} else if op == 3 {
			pos, num := io.NextInt(), io.NextInt()
			pos--
			operations[i] = [4]int{3, pos, num, 0}
		} else if op == 4 || op == 5 {
			start, end, num := io.NextInt(), io.NextInt(), io.NextInt()
			start--
			operations[i] = [4]int{op, start, end, num}
		}
	}

	res := Solve(nums, operations)
	for _, v := range res {
		io.Println(v)
	}
}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

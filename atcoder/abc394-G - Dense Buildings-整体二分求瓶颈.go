// abc394-G - Dense Buildings-整体二分
// https://atcoder.jp/contests/abc394/tasks/abc394_g
//
// 在一个H×W的城市网格中，每个格子有一座F_i,j层高的建筑。
// 高桥君有两种移动方式：
//
// 使用楼梯在同一建筑内上下移动（每次上下一层计为1次使用楼梯）
// 使用空中通道从当前建筑的X层横向移动到相邻建筑的X层（前提是目标建筑至少有X层）
// !计算从起点建筑的某一层到终点建筑的某一层所需的最少楼梯使用次数。
//
// !本质上是找到一条从起点到终点的路径，使得路径上所有建筑的最低高度尽可能高（因为这决定了我们可以在多高的楼层使用空中通道）
// 使用整体二分求两点之间所有路径中的点权最小值的最大值.

package main

import (
	"bufio"
	"fmt"
	stdio "io"
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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	h, w := io.NextInt(), io.NextInt()
	grid := make([][]int, h)
	for i := range grid {
		grid[i] = make([]int, w)
		for j := range grid[i] {
			grid[i][j] = io.NextInt()
		}
	}

	q := io.NextInt()
	type query struct{ x1, y1, h1, x2, y2, h2 int }
	queries := make([]query, 0, q)
	for i := 0; i < q; i++ {
		x1, y1, h1 := io.NextInt(), io.NextInt(), io.NextInt()
		x1, y1 = x1-1, y1-1
		x2, y2, h2 := io.NextInt(), io.NextInt(), io.NextInt()
		x2, y2 = x2-1, y2-1
		queries = append(queries, query{x1, y1, h1, x2, y2, h2})
	}

	weightToEdges := func() map[int][][2]int {
		res := make(map[int][][2]int)
		for i := 0; i < h; i++ {
			for j := 0; j < w; j++ {
				if i+1 < h {
					cost := min(grid[i][j], grid[i+1][j])
					res[cost] = append(res[cost], [2]int{i*w + j, (i+1)*w + j})
				}
				if j+1 < w {
					cost := min(grid[i][j], grid[i][j+1])
					res[cost] = append(res[cost], [2]int{i*w + j, i*w + j + 1})
				}
			}
		}
		return res
	}()
	weightsReverseSorted := func() []int {
		res := make([]int, 0, len(weightToEdges))
		for weight := range weightToEdges {
			res = append(res, weight)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(res)))
		return res
	}()

	// !整体二分求两点之间所有路径中的点权最小值的最大值
	uf := NewUnionFindArray(h * w)
	reset := func() { uf.Clear() }
	mutate := func(mutationId int) {
		for _, edge := range weightToEdges[weightsReverseSorted[mutationId]] {
			uf.Union(edge[0], edge[1])
		}
	}
	predicate := func(queryId int) bool {
		q := queries[queryId]
		return uf.Find(q.x1*w+q.y1) == uf.Find(q.x2*w+q.y2)
	}

	res := ParallelBinarySearch(len(weightsReverseSorted), q, reset, mutate, predicate)
	for i, q := range queries {
		h1, h2 := q.h1, q.h2
		if res[i] == -1 {
			io.Println(abs(h1 - h2))
			continue
		}
		minH := min(min(h1, h2), weightsReverseSorted[res[i]])
		io.Println(h1 + h2 - 2*minH)
	}
}

// 整体二分解决这样一类问题:
//   - 给定一个长度为n的操作序列, 按顺序执行这些操作;
//   - 给定q个查询,每个查询形如:"条件qi为真(满足条件)是在第几次操作之后?".
//     !要求对条件为真的判定具有单调性，即某个操作后qi为真,后续操作都会满足qi为真.
//
// 返回:
//   - -1 => 不需要操作就满足条件的查询.
//   - [0, n) => 满足条件的最早的操作的编号(0-based).
//   - n => 执行完所有操作后都不满足条件的查询.
//
// https://betrue12.hateblo.jp/entry/2019/08/14/152227
func ParallelBinarySearch(
	n, q int,
	reset func(), // 重置操作序列，一共调用 logn 次.
	mutate func(mutationId int), // 执行第 mutationId 次操作，一共调用 nlogn 次.
	predicate func(queryId int) bool, // 判断第 queryId 次查询是否满足条件，一共调用 qlogn 次.
) []int {
	left, right := make([]int, q), make([]int, q)
	for i := 0; i < q; i++ {
		left[i], right[i] = 0, n
	}

	// 不需要操作就满足条件的查询
	for i := 0; i < q; i++ {
		if predicate(i) {
			right[i] = -1
		}
	}

	for {
		mids := make([]int, q)
		for i := range mids {
			mids[i] = -1
		}
		for i := 0; i < q; i++ {
			if left[i] <= right[i] {
				mids[i] = (left[i] + right[i]) >> 1
			}
		}

		// csr 数组保存二元对 (qi,mid).
		indeg := make([]int, n+2)
		for i := 0; i < q; i++ {
			mid := mids[i]
			if mid != -1 {
				indeg[mid+1]++
			}
		}
		for i := 0; i < n+1; i++ {
			indeg[i+1] += indeg[i]
		}
		total := indeg[n+1]
		if total == 0 {
			break
		}
		counter := append(indeg[:0:0], indeg...)
		csr := make([]int, total)
		for i := 0; i < q; i++ {
			mid := mids[i]
			if mid != -1 {
				csr[counter[mid]] = i
				counter[mid]++
			}
		}

		reset()
		times := 0
		for _, pos := range csr {
			for times < mids[pos] {
				mutate(times)
				times++
			}
			if predicate(pos) {
				right[pos] = times - 1
			} else {
				left[pos] = times + 1
			}
		}
	}

	return right
}

type UnionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{data: data}
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) Clear() {
	for i := range ufa.data {
		ufa.data[i] = -1
	}
}

func (ufa *UnionFindArray) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

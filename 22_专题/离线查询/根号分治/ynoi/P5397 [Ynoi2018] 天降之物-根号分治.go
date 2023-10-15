// P5397 [Ynoi2018] 天降之物-根号分治 (第四分块)
// https://www.cnblogs.com/DCH233/p/17315715.html

// 给你一个长为n的序列，有两种操作：
// 1 x y: 将序列中所有值为x的数变为y
// 2 x y: 询问序列中x与y的最近距离.不存在则输出Ikaros.
// 本题强制在线，每次的 x,y 需要 xor 上次答案.如果输出 Ikaros，或是第一次查询, 则上次答案为 0.
// 数组中所有数的值在 1 到 1e5 之间.
//
//  1. 先考虑静态问题，根据频率根号分治，小数暴力归并，大数预处理答案.
//  2. 考虑合并操作，建立一个实际颜色对应代表颜色的映射；那么在合并时就可以强制钦定小的向大的合并
//     假设 x 的出现次数小于等于 y 的出现次数，暴力的想法是每次合并两个数就直接重构.
//     但是这样会导致每次合并的复杂度为 O(n).
//     注意到如果合并的两个数都是大数，这种情况最多出现O(sqrt(n))次，可以接受.
//     考虑两个数中至少有一个是小数.
//     !每个数维护一个邻接表(缓冲), 相当于加入 x的数的懒标记。
//     在合并时如果邻接表中的数超过 sqrt(n) 个，就直接重构(大数)，重构后的大数不存在邻接表.
//     查询时需要把每个数的邻接表中的数计算进来.
//     每次邻接表合并的复杂度为 O(sqrt(n))，`重构大数的次数为O(sqrt(n))`，因此总体复杂度为 O(nsqrt(n)).

package main

import (
	"bufio"
	"fmt"
	"math"

	stdio "io"
	"os"
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

	const N int = int(1e5 + 10)
	SQRT := 2*int(math.Sqrt(float64(n))) + 1
	facade := [N]int{}   // realValue -> innerValue
	adjMap := [N][]int{} // innerValue -> indices
	largeId := NewDictionary()
	largeRes := make([][]int, (n/SQRT)+10) // largeId -> res
	for i := range facade {
		facade[i] = -1
	}
	for i, v := range nums {
		facade[v] = v
		adjMap[v] = append(adjMap[v], i)
	}

	// !重构大数.
	rebuildLarge := func(x int) {
		largeId := largeId.Id(x)
		res := make([]int, N)
		for i := range res {
			res[i] = n
		}
		dist := n
		for i := 0; i < n; i++ {
			cur := nums[i]
			if cur == x {
				dist = 0
			} else {
				dist++
				res[cur] = min(res[cur], dist)
			}
		}
		dist = n
		for i := n - 1; i >= 0; i-- {
			cur := nums[i]
			if cur == x {
				dist = 0
			} else {
				dist++
				res[cur] = min(res[cur], dist)
			}
		}
		largeRes[largeId] = res
		adjMap[x] = nil // !清空邻接表.
	}

	for v := range adjMap {
		if len(adjMap[v]) > SQRT {
			rebuildLarge(v)
		}
	}

	// !将所有的x变为y.
	update := func(x, y int) {
		px, py := facade[x], facade[y]
		if px == py || px == -1 {
			return
		}
		if py == -1 {
			facade[y] = px
			facade[x] = -1
			return
		}

		if largeId.Has(px) {
			facade[x], facade[y] = facade[y], facade[x]
			px, py = py, px
		}
		for i := 0; i < largeId.Size(); i++ {
			res := largeRes[i]
			res[py] = min(res[py], res[px])
		}

		if largeId.Has(px) {
			for i := 0; i < n; i++ {
				if nums[i] == px {
					nums[i] = py
				}
			}
			rebuildLarge(py)
		} else {
			for _, i := range adjMap[px] {
				nums[i] = py
			}
			adjMap[py] = MergeTwoSortedArray(adjMap[py], adjMap[px])
			if len(adjMap[py]) > SQRT {
				rebuildLarge(py)
			}
		}

		adjMap[px] = nil
		facade[x] = -1
	}

	// !将所有的x变为y.指针写法.
	// update := func(x, y int) {
	// 	px, py := &facade[x], &facade[y]
	// 	if *px == *py || *px == -1 {
	// 		return
	// 	}
	// 	if *py == -1 {
	// 		*py = *px
	// 		*px = -1
	// 		return
	// 	}

	// 	if largeId.Has(*px) {
	// 		*px, *py = *py, *px
	// 	}
	// 	for i := 0; i < largeId.Size(); i++ {
	// 		res := largeRes[i]
	// 		res[*py] = min(res[*py], res[*px])
	// 	}

	// 	if largeId.Has(*px) {
	// 		for i := 0; i < n; i++ {
	// 			if nums[i] == *px {
	// 				nums[i] = *py
	// 			}
	// 		}
	// 		rebuildLarge(*py)
	// 	} else {
	// 		for _, i := range adjMap[*px] {
	// 			nums[i] = *py
	// 		}
	// 		adjMap[*py] = MergeTwoSortedArray(adjMap[*py], adjMap[*px])
	// 		if len(adjMap[*py]) > SQRT {
	// 			rebuildLarge(*py)
	// 		}
	// 	}

	// 	adjMap[*px] = nil
	// 	*px = -1
	// }

	// !查询x和y的最近距离.
	query := func(x, y int) (res int, ok bool) {
		x, y = facade[x], facade[y]
		if x == -1 || y == -1 {
			return
		}
		if x == y {
			return 0, true
		}

		res = n
		pos1, pos2 := adjMap[x], adjMap[y]
		i, j := 0, 0
		for i < len(pos1) && j < len(pos2) {
			res = min(res, abs(pos1[i]-pos2[j]))
			if pos1[i] < pos2[j] {
				i++
			} else {
				j++
			}
		}

		if largeId.Has(x) {
			id := largeId.Id(x)
			res = min(res, largeRes[id][y])
		}
		if largeId.Has(y) {
			id := largeId.Id(y)
			res = min(res, largeRes[id][x])
		}
		return res, true
	}

	preRes := 0
	for qi := 0; qi < q; qi++ {
		op, x, y := io.NextInt(), io.NextInt()^preRes, io.NextInt()^preRes
		if op == 1 {
			update(x, y)
		} else {
			res, ok := query(x, y)
			if !ok {
				io.Println("Ikaros")
				preRes = 0
			} else {
				io.Println(res)
				preRes = res
			}
		}
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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func MergeTwoSortedArray(a, b []int) []int {
	res := make([]int, len(a)+len(b))
	i, j, k := 0, 0, 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			res[k] = a[i]
			i++
		} else {
			res[k] = b[j]
			j++
		}
		k++
	}
	for i < len(a) {
		res[k] = a[i]
		i++
		k++
	}
	for j < len(b) {
		res[k] = b[j]
		j++
		k++
	}
	return res
}

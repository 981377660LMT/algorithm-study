// https://atcoder.jp/contests/abc266/submissions/34262996 (树状数组套树状数组)
// https://www.luogu.com.cn/problem/P3810 (三维偏序)

package main

import (
	"bufio"
	"fmt"
	stdio "io"
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

	n := io.NextInt()
	goods := make([][4]int, n)
	for i := 0; i < n; i++ {
		goods[i] = [4]int{io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt()}
	}
	io.Println(snukePanic(goods, 0, 0))
}

// 二维打地鼠,每次可以向左右移动或者`向上移动`
// !注意不能向下移动
// n个地鼠,每个地鼠有三个参数: t, x, y, s, t表示地鼠出现的时间, x, y表示地鼠出现的坐标, s表示地鼠的分数
// 初始时在(startX,startY),求最大分数(保证相同时间相同位置最多仅有一只地鼠)
//
// dp[t][x][y] 表示在时间t,坐标(x, y)时的最大分数
// dp[ti][xi][yi] = max(dp[tj][xj][yj]) + s,其中
// !xi-xj+yi-yj<=ti-tj,-(xi-xj)+yi-yj<=ti-tj
// 令A=ti-xi-yi, B=ti+xi-yi,则可以转化为三维偏序问题
// !dp[Ai][Bi][yi] = max(dp[Aj][Bj][yj]) + s,其中 Aj<=Ai, Bj<=Bi, yj<=yi
// 排序+二维树状数组解决
func snukePanic(goods [][4]int, startX, startY int) int {
	n := len(goods)
	newGoods := make([][4]int, 0, n) // y, A, B, score
	A, B := make([]int, 0, n), make([]int, 0, n)
	A, B = append(A, 0), append(B, 0) // 初始坐标
	for i := 0; i < n; i++ {
		t, x, y, s := goods[i][0], goods[i][1], goods[i][2], goods[i][3]
		x, y = x-startX, y-startY
		if t-x-y >= 0 && t+x-y >= 0 {
			newGoods = append(newGoods, [4]int{y, t - x - y, t + x - y, s})
			A = append(A, t-x-y)
			B = append(B, t+x-y)
		}
	}
	sort.Slice(newGoods, func(i, j int) bool { // 按照y维度扫描,维护另外两个维度的前缀最大值
		if newGoods[i][0] == newGoods[j][0] {
			if newGoods[i][1] == newGoods[j][1] {
				return newGoods[i][2] < newGoods[j][2]
			}
			return newGoods[i][1] < newGoods[j][1]
		}
		return newGoods[i][0] < newGoods[j][0]
	})

	sortedA, mp1 := sortedSet(A) // 离散化
	sortedB, mp2 := sortedSet(B)

	// 查询二维前缀最大值
	bit2d := NewBIT2D(len(sortedA), len(sortedB))
	res := 0
	for _, good := range newGoods {
		a, b, s := mp1[good[1]], mp2[good[2]], good[3]
		curMax := bit2d.Query(a+1, b+1) + s
		res = max(res, curMax)
		bit2d.Upfate(a, b, curMax)
	}
	return res
}

func sortedSet(nums []int) (sorted []int, rank map[int]int) {
	set := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	rank = make(map[int]int, len(sorted))
	for i, v := range sorted {
		rank[v] = i
	}
	return
}

// !幺元为0, 二元操作为op
var op = max

type BIT2D struct {
	h, w int
	data []*bit1D
}

// 二维树状数组,初始化height个一维树状数组,每个一维树状数组管理的纵坐标长度为width.
//  需要提前将所有点的坐标离散化.
func NewBIT2D(height, width int) *BIT2D {
	bits := make([]*bit1D, height)
	for i := range bits {
		bits[i] = newBIT(width)
	}
	return &BIT2D{height, width, bits}
}

// 单点更新(x,y)处的元素,x和y都是离散化后的坐标.
// 0 <= x < h, 0 <= y < w
func (f *BIT2D) Upfate(x, y, value int) {
	for x++; x <= f.h; x += x & -x {
		f.data[x-1].Update(y, value)
	}
}

// 查询前缀区间 [0,rightX) * [0,rightY) 的值, rightX和rightY都是离散化后的坐标.
// 0 <= rightX <= h, 0 <= rightY <= w
func (f *BIT2D) Query(rightX, rightY int) int {
	res := 0
	if rightX > f.h {
		rightX = f.h
	}
	if rightY > f.w {
		rightY = f.w
	}
	for ; rightX > 0; rightX -= rightX & -rightX {
		res = op(res, f.data[rightX-1].Query(rightY))
	}
	return res
}

type bit1D struct {
	n    int
	data map[int]int
}

func newBIT(n int) *bit1D {
	return &bit1D{n, make(map[int]int)}
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *bit1D) Update(index int, value int) {
	for index++; index <= f.n; index += index & -index {
		f.data[index-1] = op(f.data[index-1], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= right <= n
func (f *bit1D) Query(right int) int {
	res := 0
	if right > f.n {
		right = f.n
	}
	for ; right > 0; right -= right & -right {
		res = op(res, f.data[right-1])
	}
	return res
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

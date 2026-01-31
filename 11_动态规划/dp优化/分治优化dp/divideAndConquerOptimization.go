// 决策单调性优化dp.
// 単調最小値DP (aka. 分割統治DP) 优化 offlineDp
// https://ei1333.github.io/library/dp/divide-and-conquer-optimization.hpp
// !用于高速化 dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j<=n) => !将区间[0,n)分成k组的最小代价
//
//	如果f满足决策单调性 那么对转移的每一行，可以采用 monotoneminima 寻找最值点
//	O(kn^2)优化到O(knlogn)
//
// https://www.cnblogs.com/alex-wei/p/DP_optimization_method_II.html
// https://www.cnblogs.com/purplevine/p/16990286.html
// https://www.luogu.com/article/vx7a76on

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

const INF int = 1e18

func main() {
	// CF321E()
	// CF833B()
	// CF868F()
	P5574()
}

// Ciel and Gondolas
// https://www.luogu.com.cn/problem/CF321E
//
// 转移代价为二维前缀和.
func CF321E() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, k := io.NextInt(), io.NextInt()
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			grid[i][j] = io.NextInt()
		}
	}

	preSum2d := NewPreSum2DFrom(grid)

	f := func(i, j int, _ int) int {
		res := preSum2d.QueryRange(i, i, j-1, j-1)
		return res
	}
	dp := DivideAndConquerOptimization(k, n, f)
	res := dp[k][n]
	res /= 2
	io.Println(res)
}

// CF833B-The Bakery (决策单调性+莫队维护区间颜色个数)
// https://www.luogu.com.cn/problem/CF833B
// 将一个数组分为k段，使得总价值最大。
// 一段区间的价值表示为区间内不同数字的个数。
// n=3e4,k<=50
//
// dp[i][j]=max{dp[i-1][k]+cost(k+1,j) 1<=k<j
// dp[i][j]意为前j个数被分成i段时的最大总价值.
func CF833B() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, k := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}

	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	counter := make([]int, D.Size())
	left, right, kind := 0, 0, 0
	add := func(i int) {
		counter[nums[i]]++
		if counter[nums[i]] == 1 {
			kind++
		}
	}
	remove := func(i int) {
		counter[nums[i]]--
		if counter[nums[i]] == 0 {
			kind--
		}
	}
	f := func(l, r int, _ int) int {
		for left > l {
			left--
			add(left)
		}
		for right < r {
			add(right)
			right++
		}
		for left < l {
			remove(left)
			left++
		}
		for right > r {
			right--
			remove(right)
		}
		return -kind // 要求最大值，因此取负
	}

	dp := DivideAndConquerOptimization(k, n, f)
	io.Println(-dp[k][n])
}

// Yet Another Minimization Problem (决策单调性+莫队维护相同元素的对数)
// https://www.luogu.com.cn/problem/CF868F
// 有一个长度为 n 的序列，要求将其分成 k 个子段，每个子段的花费是子段内相同元素的对数，求最小花费。
// dp[k][i] 表示前 i 个元素分成 k 个子段的最小花费。
func CF868F() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, k := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}
	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	counter := make([]int, D.Size())
	left, right, cost := 0, 0, 0
	add := func(i int) {
		cost += counter[nums[i]]
		counter[nums[i]]++
	}
	remove := func(i int) {
		counter[nums[i]]--
		cost -= counter[nums[i]]
	}
	f := func(l, r int, _ int) int {
		for left > l {
			left--
			add(left)
		}
		for right < r {
			add(right)
			right++
		}
		for left < l {
			remove(left)
			left++
		}
		for right > r {
			right--
			remove(right)
		}
		return cost
	}

	dp := DivideAndConquerOptimization(k, n, f)
	io.Println(dp[k][n])
}

// P5574 [CmdOI2019] 任务分配问题 (决策单调性+莫队+树状数组维护逆序对)
// https://www.luogu.com.cn/problem/P5574
// 将数组分成 k 段，每段的代价是这段的逆序对数，求最小代价。
func P5574() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, k := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}

	for i := 0; i < n; i++ {
		nums[i] = -nums[i] // !这里的逆序对指的是前面的数小于后面，与一般的定义相反
	}

	getRank, size := DiscretizeSparse(nums, 0)
	for i, v := range nums {
		nums[i] = getRank(v)
	}

	bit := NewBitArray(size) // [0,size)
	inv := 0
	left, right := 0, 0
	addLeft := func(i int) {
		v := nums[i]
		inv += bit.QueryPrefix(v)
		bit.Add(v, 1)
	}
	addRight := func(i int) {
		v := nums[i]
		inv += bit.QueryAll() - bit.QueryPrefix(v+1)
		bit.Add(v, 1)
	}
	removeLeft := func(i int) {
		v := nums[i]
		inv -= bit.QueryPrefix(v)
		bit.Add(v, -1)
	}
	removeRight := func(i int) {
		v := nums[i]
		inv -= bit.QueryAll() - bit.QueryPrefix(v+1)
		bit.Add(v, -1)
	}

	f := func(l, r int, _ int) int {
		for left > l {
			left--
			addLeft(left)
		}
		for right < r {
			addRight(right)
			right++
		}
		for left < l {
			removeLeft(left)
			left++
		}
		for right > r {
			right--
			removeRight(right)
		}
		return inv
	}

	dp := DivideAndConquerOptimization(k, n, f)
	io.Println(dp[k][n])
}

// !f(i,j,step): 左闭右开区间[i,j)的代价(0<=i<j<=n)
//
//	可选:step表示当前在第几组(1<=step<=k)
func DivideAndConquerOptimization(k, n int, f func(i, j, step int) int) [][]int {
	dp := make([][]int, k+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0

	for k_ := 1; k_ <= k; k_++ {
		getCost := func(y, x int) int {
			if x >= y {
				return INF
			}
			return dp[k_-1][x] + f(x, y, k_)
		}
		res := monotoneminima(n+1, n+1, getCost)
		for j := 0; j <= n; j++ {
			dp[k_][j] = res[j][1]
		}
	}

	return dp
}

// 对每个 0<=i<H 求出 f(i,j) 取得最小值的 (j, f(i,j)) (0<=j<W)
func monotoneminima(H, W int, f func(i, j int) int) [][2]int {
	dp := make([][2]int, H) // dp[i] 表示第i行取到`最小值`的(索引,值)

	var dfs func(top, bottom, left, right int)
	dfs = func(top, bottom, left, right int) {
		if top > bottom {
			return
		}

		mid := (top + bottom) >> 1
		index := -1
		res := 0
		for i := left; i <= right; i++ {
			tmp := f(mid, i)
			if index == -1 || tmp < res { // !less if get min
				index = i
				res = tmp
			}
		}
		dp[mid] = [2]int{index, res}
		dfs(top, mid-1, left, index)
		dfs(mid+1, bottom, index, right)
	}

	dfs(0, H-1, 0, W-1)
	return dp
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

type PreSum2D struct {
	preSum [][]int
}

func NewPreSum2D(row, col int, f func(int, int) int) *PreSum2D {
	preSum := make([][]int, row+1)
	for i := range preSum {
		preSum[i] = make([]int, col+1)
	}
	for r := 0; r < row; r++ {
		for c := 0; c < col; c++ {
			preSum[r+1][c+1] = f(r, c) + preSum[r][c+1] + preSum[r+1][c] - preSum[r][c]
		}
	}
	return &PreSum2D{preSum}
}

func NewPreSum2DFrom(mat [][]int) *PreSum2D {
	return NewPreSum2D(len(mat), len(mat[0]), func(r, c int) int { return mat[r][c] })
}

// 查询sum(A[r1:r2+1, c1:c2+1])的值.
// 0 <= r1 <= r2 < row, 0 <= c1 <= c2 < col.
func (ps *PreSum2D) QueryRange(row1, col1, row2, col2 int) int {
	return ps.preSum[row2+1][col2+1] - ps.preSum[row2+1][col1] - ps.preSum[row1][col2+1] + ps.preSum[row1][col1]
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

// (松)离散化.
//
//	offset: 离散化的起始值偏移量.
//
//	getRank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//	count: 离散化(去重)后的元素个数.
func DiscretizeSparse(nums []int, offset int) (getRank func(int) int, count int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	count = len(set)
	allNums := make([]int, 0, count)
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	getRank = func(x int) int { return sort.SearchInts(allNums, x) + offset }
	return
}

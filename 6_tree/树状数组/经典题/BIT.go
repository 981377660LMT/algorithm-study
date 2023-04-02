package main

import (
	"fmt"
	"strings"
)

func demo() {
	bitArray := NewBITArrayWithIntSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(bitArray)
	bitArray.Add(1, 1)
	fmt.Println(bitArray)
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n    int
	tree []int
}

func NewBITArray(n int) *BITArray {
	return &BITArray{n: n, tree: make([]int, n+1)}
}

func NewBITArrayWithIntSlice(nums []int) *BITArray {
	bitArray := &BITArray{
		n:    len(nums),
		tree: make([]int, len(nums)+1),
	}
	bitArray.Build(nums)
	return bitArray
}

// 常数优化: dp O(n) 建树
// https://oi-wiki.org/ds/fenwick/#tricks
func (b *BITArray) Build(nums []int) {
	for i := 1; i < len(b.tree); i++ {
		b.tree[i] += nums[i-1]
		if j := i + (i & -i); j < len(b.tree) {
			b.tree[j] += b.tree[i]
		}
	}
}

// 位置 index 增加 delta
//  1<=i<=n
func (b *BITArray) Add(index int, delta int) {
	for ; index < len(b.tree); index += index & -index {
		b.tree[index] += delta
	}
}

// 求前缀和
//  1<=i<=n
func (b *BITArray) Query(index int) (res int) {
	if index > b.n {
		index = b.n
	}
	for ; index > 0; index &= (index - 1) {
		res += b.tree[index]
	}
	return
}

// 1<=left<=right<=n
func (b *BITArray) QueryRange(left, right int) int {
	return b.Query(right) - b.Query(left-1)
}

func (b *BITArray) Len() int {
	return b.n
}

func (b *BITArray) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITArray{")
	for i := 1; i <= b.n; i++ {
		sb.WriteString(fmt.Sprintf("%d", b.QueryRange(i, i)))
		if i != b.n {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

//
//
//
// !Range Add Range Sum, 0-based.
type BITArray2 struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITArray2(n int) *BITArray2 {
	return &BITArray2{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

// 切片内[start, end)的每个元素加上delta.
//  0<=start<=end<=n
func (b *BITArray2) Add(start, end, delta int) {
	end--
	b.add(start, delta)
	b.add(end+1, -delta)
}

// 求切片内[start, end)的和.
//  0<=start<=end<=n
func (b *BITArray2) Query(start, end int) int {
	end--
	return b.query(end) - b.query(start-1)
}

func (b *BITArray2) add(index, delta int) {
	index++
	rawIndex := index
	for index <= b.n {
		b.tree1[index] += delta
		b.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (b *BITArray2) query(index int) (res int) {
	index++
	if index > b.n {
		index = b.n
	}
	rawIndex := index
	for index > 0 {
		res += rawIndex*b.tree1[index] - b.tree2[index]
		index &= (index - 1)
	}
	return
}

//
//
//
// !一维区间查询 区间修改, 0-based.
type BIT interface {
	// 区间 [left, right) 里的每个数增加 delta
	Add(left, right, delta int)

	// 查询区间 [left, right) 的和
	Query(left, right int) int
}

func NewBIT(n int) BIT {
	if n <= int(1e6) {
		return NewBITSlice(n)
	}

	return NewBITMap(n)
}

//  tree: map[int]int or []int
type BITMap struct {
	n     int
	tree1 map[int]int
	tree2 map[int]int
}

type BITSlice struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITMap(n int) *BITMap {
	return &BITMap{
		n:     n,
		tree1: make(map[int]int, 1<<10),
		tree2: make(map[int]int, 1<<10),
	}
}

func NewBITSlice(n int) *BITSlice {
	return &BITSlice{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

func (bit *BITMap) Add(left, right, delta int) {
	right--
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *BITSlice) Add(left, right, delta int) {
	right--
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *BITMap) Query(left, right int) int {
	right--
	return bit.query(right) - bit.query(left-1)
}

func (bit *BITSlice) Query(left, right int) int {
	right--
	return bit.query(right) - bit.query(left-1)
}

func (bit *BITMap) add(index, delta int) {
	index++
	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *BITSlice) add(index, delta int) {
	index++

	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *BITMap) query(index int) int {
	index++
	if index > bit.n {
		index = bit.n
	}
	rawIndex := index
	res := 0
	for index > 0 {
		res += rawIndex*bit.tree1[index] - bit.tree2[index]
		index &= (index - 1)
	}
	return res
}

func (bit *BITSlice) query(index int) int {
	index++
	if index > bit.n {
		index = bit.n
	}
	rawIndex := index
	res := 0
	for index > 0 {
		res += rawIndex*bit.tree1[index] - bit.tree2[index]
		index &= (index - 1)
	}
	return res
}

func (bit *BITMap) String() string {
	return "not implemented"
}

func (bit *BITSlice) String() string {
	nums := make([]int, bit.n+1)
	for i := 0; i < bit.n; i++ {
		nums[i+1] = bit.Query(i+1, i+1)
	}
	return fmt.Sprint(nums)
}

func (bit *BITMap) Len() int {
	return bit.n
}

func (bit *BITSlice) Len() int {
	return bit.n
}

func maximumWhiteTiles(tiles [][]int, carpetLen int) int {
	bit := NewBIT(1e9)
	for _, tile := range tiles {
		bit.Add(int(tile[0]), int(tile[1]), 1)
	}

	res := 0
	for _, tile := range tiles {
		res = max(res, bit.Query(int(tile[0]), int(tile[0]+carpetLen-1)))
	}

	return res
}

//
//
//
// !二维区间查询 区间修改
type BIT2D struct {
	row   int
	col   int
	tree1 [][]int
	tree2 [][]int
	tree3 [][]int
	tree4 [][]int
}

func NewBIT2D(row, col int) *BIT2D {
	t1, t2, t3, t4 := make([][]int, row+1), make([][]int, row+1), make([][]int, row+1), make([][]int, row+1)
	for i := 0; i <= row; i++ {
		t1[i] = make([]int, col+1)
		t2[i] = make([]int, col+1)
		t3[i] = make([]int, col+1)
		t4[i] = make([]int, col+1)
	}
	return &BIT2D{
		row:   row,
		col:   col,
		tree1: t1,
		tree2: t2,
		tree3: t3,
		tree4: t4,
	}
}

//  (row1,col1) 到 (row2,col2) 里的每一个点的值加上delta
//   0<=row1<=row2<=ROW-1, 0<=col1<=col2<=COL-1
func (b *BIT2D) Add(row1 int, col1 int, row2 int, col2 int, delta int) {
	b.add(row1, col1, delta)
	b.add(row2+1, col1, -delta)
	b.add(row1, col2+1, -delta)
	b.add(row2+1, col2+1, delta)
}

// 查询左上角 (row1,col1) 到右下角 (row2,col2) 的和
//  0<=row1<=row2<=ROW-1, 0<=col1<=col2<=COL-1
func (b *BIT2D) Query(row1 int, col1 int, row2 int, col2 int) int {
	return b.query(row2, col2) - b.query(row2, col1-1) - b.query(row1-1, col2) + b.query(row1-1, col1-1)
}

func (b *BIT2D) add(row int, col int, delta int) {
	row, col = row+1, col+1
	preRow, preCol := row, col
	for curRow := row; curRow <= b.row; curRow += curRow & -curRow {
		for curCol := col; curCol <= b.col; curCol += curCol & -curCol {
			b.tree1[curRow][curCol] += delta
			b.tree2[curRow][curCol] += (preRow - 1) * delta
			b.tree3[curRow][curCol] += (preCol - 1) * delta
			b.tree4[curRow][curCol] += (preRow - 1) * (preCol - 1) * delta
		}
	}
}

func (b *BIT2D) query(row, col int) (res int) {
	row, col = row+1, col+1
	if row > b.row {
		row = b.row
	}
	if col > b.col {
		col = b.col
	}

	preR, preC := row, col
	for r := row; r > 0; r -= r & -r {
		for c := col; c > 0; c -= c & -c {
			res += preR*preC*b.tree1[r][c] -
				preC*b.tree2[r][c] -
				preR*b.tree3[r][c] +
				b.tree4[r][c]
		}
	}

	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type NumMatrix struct {
	matrix [][]int
	row    int
	col    int
	bit    BIT2D
}

func Constructor(matrix [][]int) NumMatrix {
	n := len(matrix)
	m := len(matrix[0])
	bit := NewBIT2D(n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			bit.Add(i, j, i, j, matrix[i][j])
		}
	}

	return NumMatrix{matrix, n, m, *bit}
}

func (this *NumMatrix) Update(row, col, val int) {
	delta := val - this.matrix[row][col]
	this.matrix[row][col] = val
	this.bit.Add(row, col, row, col, delta)
}

func (this *NumMatrix) SumRegion(row1, col1, row2, col2 int) int {
	return this.bit.Query(row1, col1, row2, col2)
}

//
//
//
//
// Fenwick Tree Prefix
// https://suisen-cp.github.io/cp-library-cpp/library/datastructure/fenwick_tree/fenwick_tree_prefix.hpp
// 如果每次都是查询前缀，那么可以使用Fenwick Tree Prefix 维护 monoid.
type S = int

func (*FenwickTreePrefix) e() S        { return 0 }
func (*FenwickTreePrefix) op(a, b S) S { return max(a, b) }
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type FenwickTreePrefix struct {
	n    int
	data []S
}

func NewFenwickTreePrefix(n int) *FenwickTreePrefix {
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 0; i < n+1; i++ {
		res.data[i] = res.e()
	}
	return res
}

func NewFenwickTreePrefixWithSlice(nums []S) *FenwickTreePrefix {
	n := len(nums)
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 1; i < n+1; i++ {
		res.data[i] = nums[i-1]
	}
	for i := 1; i < n+1; i++ {
		if j := i + (i & -i); j <= n {
			res.data[j] = res.op(res.data[j], res.data[i])
		}
	}
	return res
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *FenwickTreePrefix) Update(index int, value S) {
	for index++; index <= f.n; index += index & -index {
		f.data[index] = f.op(f.data[index], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= right <= n
func (f *FenwickTreePrefix) Query(right int) S {
	res := f.e()
	if right > f.n {
		right = f.n
	}
	for ; right > 0; right -= right & -right {
		res = f.op(res, f.data[right])
	}
	return res
}

func main() {
	pf := NewFenwickTreePrefixWithSlice([]S{1, 2, 3, 4, 14, 1, 2, 3})
	pf.Update(9, 1)
	pf.Update(1, 2)
	fmt.Println(pf.Query(100))
	pf.Update(9, 10)
	fmt.Println(pf.Query(100))
}

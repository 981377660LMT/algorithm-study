package main

import (
	"fmt"
	"strings"
)

func main() {
	bitArray := NewBITArrayWithIntSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(bitArray)
	bitArray.Add(1, 1)
	fmt.Println(bitArray)
}

// !Point Add Range Sum, 1-based.
type BITArray struct {
	n    int
	tree []int
}

func NewBITArrayWithIntSlice(nums []int) *BITArray {
	bitArray := &BITArray{
		n:    len(nums),
		tree: make([]int, len(nums)+1),
	}
	bitArray.Build(nums)
	return bitArray
}

func NewBITArray(n int) *BITArray {
	return &BITArray{n: n, tree: make([]int, n+1)}
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
	for ; index > 0; index -= index & -index {
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
// !一维区间查询 区间修改
type BIT interface {
	// 区间 [left, right] 里的每个数增加 delta
	Add(left, right, delta int)

	// 查询区间 [left, right] 的和
	Query(left, right int) int
}

func CreateBIT(n int) BIT {
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
		tree1: make([]int, n+10),
		tree2: make([]int, n+10),
	}
}

func (bit *BITMap) Add(left, right, delta int) {
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *BITSlice) Add(left, right, delta int) {
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *BITMap) Query(left, right int) int {
	return bit.query(right) - bit.query(left-1)
}

func (bit *BITSlice) Query(left, right int) int {
	return bit.query(right) - bit.query(left-1)
}

func (bit *BITMap) add(index, delta int) {
	if index <= 0 {
		errorInfo := fmt.Sprintf("index must be greater than 0, but got %d", index)
		panic(errorInfo)
	}

	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *BITSlice) add(index, delta int) {
	if index <= 0 {
		errorInfo := fmt.Sprintf("index must be greater than 0, but got %d", index)
		panic(errorInfo)
	}

	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *BITMap) query(index int) int {
	if index > bit.n {
		index = bit.n
	}

	rawIndex := index
	res := 0
	for index > 0 {
		res += rawIndex*bit.tree1[index] - bit.tree2[index]
		index -= index & -index
	}
	return res
}

func (bit *BITSlice) query(index int) int {
	if index > bit.n {
		index = bit.n
	}

	rawIndex := index
	res := 0
	for index > 0 {
		res += rawIndex*bit.tree1[index] - bit.tree2[index]
		index -= index & -index
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
	bit := CreateBIT(1e9)
	for _, tile := range tiles {
		bit.Add(int(tile[0]), int(tile[1]), 1)
	}

	res := 0
	for _, tile := range tiles {
		res = max(res, bit.Query(int(tile[0]), int(tile[0]+carpetLen-1)))
	}

	return res
}

// !二维区间查询 区间修改
type BIT2D struct {
	row   int
	col   int
	tree1 map[int]map[int]int
	tree2 map[int]map[int]int
	tree3 map[int]map[int]int
	tree4 map[int]map[int]int
}

func NewBIT2D(row, col int) *BIT2D {
	return &BIT2D{
		row:   row,
		col:   col,
		tree1: make(map[int]map[int]int, 1<<4),
		tree2: make(map[int]map[int]int, 1<<4),
		tree3: make(map[int]map[int]int, 1<<4),
		tree4: make(map[int]map[int]int, 1<<4),
	}
}

func (b *BIT2D) UpdateRange(row1 int, col1 int, row2 int, col2 int, delta int) {
	b.update(row1, col1, delta)
	b.update(row2+1, col1, -delta)
	b.update(row1, col2+1, -delta)
	b.update(row2+1, col2+1, delta)
}

func (b *BIT2D) QueryRange(row1 int, col1 int, row2 int, col2 int) int {
	return b.query(row2, col2) - b.query(row2, col1-1) - b.query(row1-1, col2) + b.query(row1-1, col1-1)
}

func (b *BIT2D) update(row int, col int, delta int) {
	row, col = row+1, col+1
	preRow, preCol := row, col

	for curRow := row; curRow <= b.row; curRow += curRow & -curRow {
		for curCol := col; curCol <= b.col; curCol += curCol & -curCol {
			setDeep(b.tree1, curRow, curCol, delta)
			setDeep(b.tree2, curRow, curCol, (preRow-1)*delta)
			setDeep(b.tree3, curRow, curCol, (preCol-1)*delta)
			setDeep(b.tree4, curRow, curCol, (preRow-1)*(preCol-1)*delta)
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
	for curR := row; curR > 0; curR -= curR & -curR {
		for curC := col; curC > 0; curC -= curC & -curC {
			res += preR*preC*getDeep(b.tree1, curR, curC) - preC*getDeep(b.tree2, curR, curC) - preR*getDeep(b.tree3, curR, curC) + getDeep(b.tree4, curR, curC)
		}
	}

	return
}

func setDeep(mp map[int]map[int]int, key1, key2, delta int) {
	if _, ok := mp[key1]; !ok {
		mp[key1] = make(map[int]int, 1<<4)
	}
	mp[key1][key2] += delta
}

func getDeep(mp map[int]map[int]int, key1, key2 int) int {
	if _, ok := mp[key1]; !ok {
		return 0
	}
	return mp[key1][key2]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// `CSRArray`（Compressed Sparse Row Array，压缩稀疏行数组）是一种用于存储稀疏矩阵的数据结构。稀疏矩阵是指矩阵中大部分元素为零的矩阵，而`CSRArray`通过三个主要的数组来高效地表示这种矩阵，从而节省存储空间并加快处理速度。这三个数组分别是：

// 1. **值（Values）数组**：存储矩阵中所有非零元素的值。
// 2. **列索引（Column Indices）数组**：存储每个非零元素在其所在行中的列号。
// 3. **行偏移（Row Pointers）数组**：存储每一行第一个非零元素在值数组中的位置，以及一个额外的结束元素，表示最后一个非零元素之后的位置。

// 通过这种方式，`CSRArray`能够仅存储非零元素及其位置信息，大大减少了存储稀疏矩阵所需的空间。此外，`CSRArray`还支持高效地执行各种矩阵操作，如矩阵-向量乘法等。
// `CSRArray`的数据结构特别适合处理那些非零元素分布不均匀的大型稀疏矩阵，常见于科学计算、工程模拟、图形处理和数据挖掘等领域。

package main

func main() {
	matrix := [][]int{
		{1, 0, 0, 0},
		{0, 0, 2, 0},
		{0, 3, 0, 0},
		{0, 0, 0, 4},
	}
	csr := NewCSRArray(matrix)

	// Set a value
	csr.Set(1, 2, 5)

	// Get a value
	value := csr.Get(1, 2)
	println(value)

	// Enumerate all values in a row
	csr.Enumerate(1, 0, func(value int) bool {
		println(value, 987)
		return false
	})
}

// Compressed Sparse Row (CSR) Array
type CSRArray[E any] struct {
	csr   []E
	start []int32 // 行长度的前缀和
}

func NewCSRArray[E any](matrix [][]E) *CSRArray[E] {
	row := int32(len(matrix))
	start := make([]int32, row+1)
	for i := int32(0); i < row; i++ {
		start[i+1] = start[i] + int32(len(matrix[i]))
	}
	csr := make([]E, start[row])
	for i := int32(0); i < row; i++ {
		copy(csr[start[i]:], matrix[i])
	}
	return &CSRArray[E]{csr: csr, start: start}
}

func (c *CSRArray[E]) Set(row, col int32, value E) {
	c.csr[c.start[row]+col] = value
}

func (c *CSRArray[E]) Get(row, col int32) E {
	return c.csr[c.start[row]+col]
}

func (c *CSRArray[E]) Enumerate(row, col int32, f func(E) (shouldBreak bool)) {
	for i := c.start[row] + col; i < c.start[row+1]; i++ {
		if f(c.csr[i]) {
			break
		}
	}
}

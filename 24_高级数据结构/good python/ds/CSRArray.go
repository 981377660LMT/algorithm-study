// # # from titan_pylib.data_structures.array.csr_array import CSRArray
// # from typing import Generic, TypeVar, List, Iterator
// # from itertools import chain

// # T = TypeVar("T")

// # class CSRArray(Generic[T]):
// #     """CSR形式の配列です"""

// #     def __init__(self, a: List[List[T]]) -> None:
// #         """2次元配列 ``a`` を CSR 形式にします。

// #         Args:
// #           a (List[List[T]]): 変換する2次元配列です。
// #         """
// #         n = len(a)
// #         start = list(map(len, a))
// #         start.insert(0, 0)
// #         for i in range(n):
// #             start[i + 1] += start[i]
// #         self.csr: List[T] = list(chain(*a))
// #         self.start: List[int] = start

// #     def set(self, i: int, j: int, val: T) -> None:
// #         """インデックスを指定して値を更新します。

// #         Args:
// #           i (int): 行のインデックスです。
// #           j (int): 列のインデックスです。
// #           val (T): a[i][j] 要素を更新する値です。
// #         """
// #         self.csr[self.start[i] + j] = val

// #     def iter(self, i: int, j: int = 0) -> Iterator[T]:
// #         """行を指定してイテレートします。

// #         Args:
// #           i (int): 行のインデックスです。
// #           j (int, optional): 列のインデックスです。デフォルトは ``0`` です。
// #         """
// #         csr = self.csr
// #         for ij in range(self.start[i] + j, self.start[i + 1]):
// #             yield csr[ij]

// `CSRArray`（Compressed Sparse Row Array，压缩稀疏行数组）是一种用于存储稀疏矩阵的数据结构。稀疏矩阵是指矩阵中大部分元素为零的矩阵，而`CSRArray`通过三个主要的数组来高效地表示这种矩阵，从而节省存储空间并加快处理速度。这三个数组分别是：

// 1. **值（Values）数组**：存储矩阵中所有非零元素的值。
// 2. **列索引（Column Indices）数组**：存储每个非零元素在其所在行中的列号。
// 3. **行偏移（Row Pointers）数组**：存储每一行第一个非零元素在值数组中的位置，以及一个额外的结束元素，表示最后一个非零元素之后的位置。

// 通过这种方式，`CSRArray`能够仅存储非零元素及其位置信息，大大减少了存储稀疏矩阵所需的空间。此外，`CSRArray`还支持高效地执行各种矩阵操作，如矩阵-向量乘法等。
// `CSRArray`的数据结构特别适合处理那些非零元素分布不均匀的大型稀疏矩阵，常见于科学计算、工程模拟、图形处理和数据挖掘等领域。

package main

func main() {

}

// Compressed Sparse Row (CSR) Array
type CSRArray[E any] struct {
}

// NewCSRArray creates a new CSRArray from a 2D array.
func NewCSRArray[E any](matrix [][]E) *CSRArray[E] {}

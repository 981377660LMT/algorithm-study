// 稀疏矩阵的常用存储格式（COO、CSR、CSC）
// https://zhuanlan.zhihu.com/p/37525925
// https://zhuanlan.zhihu.com/p/188700729
// https://www.cnblogs.com/xbinworld/p/4273506.html
// https://blog.csdn.net/m0_64204369/article/details/123035598
//
// COO, coordinate format.(坐标格式,每一个元素需要用一个三元组来表示)
// CSR, compressed sparse row format.(按行groupby压缩;总体最好.)
// CSC, compressed sparse column format.(按列压缩)
// BSR, block sparse row format.(分块压缩稀疏行)
// DOK, dictionary of keys format.(键值对格式)
// ...

import { enumerateGroup } from '../../0_数组/数组api/groupby'

/** COO格式常用于从文件中进行稀疏矩阵的读写. */
class SparseMatrixCOO<E> {
  private readonly _rowIndices: Uint32Array
  private readonly _colIndices: Uint32Array
  private readonly _data: E[]

  constructor(
    grid: ArrayLike<ArrayLike<E>>,
    isEmptyCell: (row: number, col: number) => boolean = (r, c) => grid[r][c] === 0
  ) {
    const m = grid.length
    const n = grid[0].length
    const rowIndices = []
    const colIndices = []
    const values = []
    for (let r = 0; r < m; r++) {
      const row = grid[r]
      for (let c = 0; c < n; c++) {
        if (!isEmptyCell(r, c)) {
          rowIndices.push(r)
          colIndices.push(c)
          values.push(row[c])
        }
      }
    }

    this._rowIndices = new Uint32Array(rowIndices)
    this._colIndices = new Uint32Array(colIndices)
    this._data = values
  }

  print(): void {
    for (let i = 0; i < this._data.length; i++) {
      console.log(this._rowIndices[i], this._colIndices[i], this._data[i])
    }
  }
}

/** CSR格式常用于读入数据后进行稀疏矩阵计算. */
class SparseMatrixCSR<E> {
  /** 对于第 i 行而言，该行中非零元素的列索引为 indices[indptr[i]:indptr[i+1]]. */
  private readonly _indices: Uint32Array
  /** 下标i处的元素表示第i行元素的列索引. */
  private readonly _indptr: Uint32Array
  private readonly _data: E[]

  constructor(
    grid: ArrayLike<ArrayLike<E>>,
    isEmptyCell: (row: number, col: number) => boolean = (r, c) => grid[r][c] === 0
  ) {
    const m = grid.length
    const n = grid[0].length
    const indices: number[] = []
    const indptr: number[] = []
    const data: E[] = []
    for (let r = 0; r < m; r++) {
      const row = grid[r]
      for (let c = 0; c < n; c++) {
        if (!isEmptyCell(r, c)) {
          indices.push(r)
          indptr.push(c)
          data.push(row[c])
        }
      }
    }

    this._indices = new Uint32Array(grid.length)
    enumerateGroup(indices, (start, end) => {
      const row = indices[start]
      this._indices[row] = end - start
    })
    for (let i = 1; i < this._indices.length; i++) {
      this._indices[i] += this._indices[i - 1]
    }

    this._indptr = new Uint32Array(indptr)
    this._data = data
  }

  print(): void {
    let start = 0
    for (let i = 0; i < this._indices.length; i++) {
      const end = this._indices[i]
      for (let j = start; j < end; j++) {
        console.log(i, this._indptr[j], this._data[j])
      }
      start = end
    }
  }
}

if (require.main === module) {
  const grid = [
    [0, 0, 0, 0, 0, 30, 99, 0, 0],
    [7, 2, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, -5, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, -6, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 2, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, -4]
  ]

  const coo = new SparseMatrixCOO(grid)
  coo.print()

  console.log('-'.repeat(20))

  const csr = new SparseMatrixCSR(grid)
  csr.print()
}

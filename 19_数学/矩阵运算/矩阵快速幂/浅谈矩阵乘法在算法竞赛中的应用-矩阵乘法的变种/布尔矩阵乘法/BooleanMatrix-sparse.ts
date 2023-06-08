/* eslint-disable no-inner-declarations */
/* eslint-disable prefer-destructuring */

// !布尔矩阵乘法(Boolean Matrix Multiplication, BMM)
// 输入和输出矩阵的元素均为布尔值。
// !按矩阵乘法的公式运算时，可以把“乘”看成and，把“加”看成or
// 对矩阵乘法 C[i][j] |= A[i][k] & B[k][j], 它的一个直观意义是把A的行和B的列看成集合，
// A的第i行包含元素k当且仅当A[i][k]=1。
// B的第j列包含元素k当且仅当B[k][j]=1。
// !那么C[i][j]代表A的第i行和B的第j列是否包含公共元素。
//
// 一个应用是传递闭包(Transitive Closure)的加速计算。
//
// https://zhuanlan.zhihu.com/p/631804105
// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/math/matrix_bool.cc#L4
//
//
// BooleanMatrixSparse
//
//
// 这里是bitset的实现.当输入矩阵比较稀疏时可以跑得非常快.
// !!!不建议使用.
// !建议使用`BooleanSquareMatrix-dense`.

import { BitSet } from '../../../../../18_哈希/BitSet/BitSet'

/**
 * 布尔矩阵乘法,单次矩阵乘法的复杂度为 `O(n^3/32)`.
 * !当输入矩阵比较稀疏时速度比较快(5000*5000,560ms).
 * @deprecated 建议使用`BooleanSquareMatrix-dense`.
 */
class BooleanMatrixSparse {
  static eye(n: number): BooleanMatrixSparse {
    const res = new BooleanMatrixSparse(n, n)
    for (let i = 0; i < n; i++) res._bs[i].add(i)
    return res
  }

  static pow(mat: BooleanMatrixSparse, k: number): BooleanMatrixSparse {
    return mat.copy().ipow(k)
  }

  /**
   * 随机矩阵,5000*5000 => 2.4s.
   */
  static mul(mat1: BooleanMatrixSparse, mat2: BooleanMatrixSparse): BooleanMatrixSparse {
    return mat1.copy().imul(mat2)
  }

  static add(mat1: BooleanMatrixSparse, mat2: BooleanMatrixSparse): BooleanMatrixSparse {
    return mat1.copy().iadd(mat2)
  }

  readonly row: number
  readonly col: number
  private _bs: BitSet[]

  constructor(row: number, col: number, bs?: BitSet[]) {
    if (bs === void 0) {
      bs = Array(row)
      for (let i = 0; i < row; i++) bs[i] = new BitSet(col)
    }
    this.row = row
    this.col = col
    this._bs = bs
  }

  ipow(k: number): BooleanMatrixSparse {
    const res = BooleanMatrixSparse.eye(this.row)
    while (k > 0) {
      if (k & 1) res.imul(this)
      this.imul(this)
      k = Math.floor(k / 2) // !注意不能用 `k >>>= 1`, k可能超过uint32
    }
    const tmp = this._bs
    this._bs = res._bs
    res._bs = tmp
    return res
  }

  imul(other: BooleanMatrixSparse): BooleanMatrixSparse {
    const row = this.row
    const col = other.col
    const res = new BooleanMatrixSparse(row, col)
    const otherBs = other._bs
    for (let i = 0; i < row; i++) {
      const rowBs = this._bs[i]
      const resBs = res._bs[i]
      for (let j = 0; j < col; j++) {
        if (rowBs.has(j)) resBs.ior(otherBs[j])
      }
    }
    const tmp = this._bs
    this._bs = res._bs
    res._bs = tmp
    return res
  }

  iadd(mat: BooleanMatrixSparse): BooleanMatrixSparse {
    for (let i = 0; i < this.row; i++) this._bs[i].ior(mat._bs[i])
    return this
  }

  /**
   * 求出邻接矩阵`mat`的传递闭包`(mat+I)^n`.
   * 随机矩阵,2000*2000 => 4.3s.
   * @deprecated 建议使用`O(n^3/32)`的Floyd-Warshall算法.
   */
  transitiveClosure(): BooleanMatrixSparse {
    if (this.row !== this.col) throw new Error('not a square matrix')
    const n = this.row
    const trans = BooleanMatrixSparse.eye(n).iadd(this)
    trans.ipow(n)
    return trans
  }

  copy(): BooleanMatrixSparse {
    const bs = Array(this.row)
    for (let i = 0; i < this.row; i++) bs[i] = this._bs[i].copy()
    return new BooleanMatrixSparse(this.row, this.col, bs)
  }

  get(row: number, col: number): boolean {
    return this._bs[row].has(col)
  }

  set(row: number, col: number, b: boolean): void {
    if (b) {
      this._bs[row].add(col)
    } else {
      this._bs[row].discard(col)
    }
  }

  print(): void {
    const grid: Uint8Array[] = Array(this.row)
    for (let i = 0; i < this.row; i++) {
      grid[i] = new Uint8Array(this.col)
      for (let j = 0; j < this.col; j++) {
        grid[i][j] = this.get(i, j) ? 1 : 0
      }
    }
    // eslint-disable-next-line no-console
    console.table(grid)
  }
}

export { BooleanMatrixSparse }

if (require.main === module) {
  // ====================
  // 测试随机矩阵
  // BooleanMatrixSparse.mul: 2.405s
  // BooleanMatrixSparse.transitiveClosure: 4.334s
  // ====================
  // 测试稀疏矩阵
  // BooleanMatrixSparse.mul: 534.847ms
  // BooleanMatrixSparse.transitiveClosure: 4.109s
  // ====================
  // 测试稠密矩阵
  // BooleanMatrixSparse.mul: 4.150s
  // BooleanMatrixSparse.transitiveClosure: 4.359s

  const mat = new BooleanMatrixSparse(3, 3)
  mat.set(0, 0, true)
  mat.set(0, 1, true)
  mat.set(1, 2, true)
  mat.set(1, 0, true)
  mat.print()

  testRandom()
  testSparse()
  testDense()

  function testRandom(): void {
    console.log('='.repeat(20))
    console.log('测试随机矩阵')
    // !随机01矩阵
    // 5000*5000的矩阵乘法
    const N_5000 = 5000
    const mat = new BooleanMatrixSparse(N_5000, N_5000)
    for (let i = 0; i < N_5000; i++) {
      for (let j = 0; j < N_5000; j++) {
        if (Math.random() < 0.5) mat.set(i, j, true)
      }
    }
    console.time('BooleanMatrixSparse.mul')
    BooleanMatrixSparse.mul(mat, mat)
    console.timeEnd('BooleanMatrixSparse.mul')

    // 2000*2000的传递闭包
    const N_2000 = 2000
    const mat2 = new BooleanMatrixSparse(N_2000, N_2000)
    for (let i = 0; i < N_2000; i++) {
      for (let j = 0; j < N_2000; j++) {
        if (Math.random() < 0.5) mat2.set(i, j, true)
      }
    }
    console.time('BooleanMatrixSparse.transitiveClosure')
    mat2.transitiveClosure()
    console.timeEnd('BooleanMatrixSparse.transitiveClosure')
  }

  function testSparse(): void {
    console.log('='.repeat(20))
    console.log('测试稀疏矩阵')
    // !稀疏矩阵
    // 5000*5000的矩阵乘法
    const N_5000 = 5000
    const mat = new BooleanMatrixSparse(N_5000, N_5000)
    for (let i = 0; i < N_5000; i++) {
      for (let j = 0; j < N_5000; j++) {
        if (Math.random() < 0.1) mat.set(i, j, true)
      }
    }
    console.time('BooleanMatrixSparse.mul')
    BooleanMatrixSparse.mul(mat, mat)
    console.timeEnd('BooleanMatrixSparse.mul')

    // 2000*2000的传递闭包
    const N_2000 = 2000
    const mat2 = new BooleanMatrixSparse(N_2000, N_2000)
    for (let i = 0; i < N_2000; i++) {
      for (let j = 0; j < N_2000; j++) {
        if (Math.random() < 0.1) mat2.set(i, j, true)
      }
    }
    console.time('BooleanMatrixSparse.transitiveClosure')
    mat2.transitiveClosure()
    console.timeEnd('BooleanMatrixSparse.transitiveClosure')
  }

  function testDense(): void {
    console.log('='.repeat(20))
    console.log('测试稠密矩阵')
    // !稠密矩阵
    // 5000*5000的矩阵乘法
    const N_5000 = 5000
    const mat = new BooleanMatrixSparse(N_5000, N_5000)
    for (let i = 0; i < N_5000; i++) {
      for (let j = 0; j < N_5000; j++) {
        mat.set(i, j, true)
      }
    }
    console.time('BooleanMatrixSparse.mul')
    BooleanMatrixSparse.mul(mat, mat)
    console.timeEnd('BooleanMatrixSparse.mul')

    // 2000*2000的传递闭包
    const N_2000 = 2000
    const mat2 = new BooleanMatrixSparse(N_2000, N_2000)
    for (let i = 0; i < N_2000; i++) {
      for (let j = 0; j < N_2000; j++) {
        mat2.set(i, j, true)
      }
    }
    console.time('BooleanMatrixSparse.transitiveClosure')
    mat2.transitiveClosure()
    console.timeEnd('BooleanMatrixSparse.transitiveClosure')
  }
}

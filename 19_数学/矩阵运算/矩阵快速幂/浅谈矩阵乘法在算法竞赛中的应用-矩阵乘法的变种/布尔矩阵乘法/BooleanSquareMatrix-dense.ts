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
// BooleanSquareMatrixDense
//
//
// 这里是bitset+four russians mathod的实现.适用于输入矩阵稠密的情形.
// complex: O(n^3 / wlogn)

import { BitSet } from '../../../../../18_哈希/BitSet/BitSet'

/**
 * 布尔方针.单次矩阵乘法的复杂度为 `O(n^3/32logn)`.这里的logn为分块的大小.
 * 适用于输入矩阵稠密的情形.
 */
class BooleanSquareMatrixDense {
  static eye(n: number): BooleanSquareMatrixDense {
    const res = new BooleanSquareMatrixDense(n)
    for (let i = 0; i < n; i++) res._bs[i].add(i)
    return res
  }

  static pow(mat: BooleanSquareMatrixDense, k: number): BooleanSquareMatrixDense {
    return mat.copy().ipow(k)
  }

  /**
   * 稠密矩阵,5000*5000 => 920ms.
   */
  static mul(
    mat1: BooleanSquareMatrixDense,
    mat2: BooleanSquareMatrixDense
  ): BooleanSquareMatrixDense {
    return mat1.copy().imul(mat2)
  }

  static add(
    mat1: BooleanSquareMatrixDense,
    mat2: BooleanSquareMatrixDense
  ): BooleanSquareMatrixDense {
    return mat1.copy().iadd(mat2)
  }

  private static _trailingZeros32 = BooleanSquareMatrixDense._initTrailingZeros32()

  private static _initTrailingZeros32(n = 1e4 + 10): Uint8Array {
    const res = new Uint8Array(n + 1)
    res[0] = 32
    for (let i = 1; i < res.length; i++) res[i] = 31 - Math.clz32(i & -i)
    return res
  }

  readonly n: number
  private _bs: BitSet[]
  private _dp?: BitSet[] // 在计算矩阵乘法时用到

  /**
   * @param n n<=1e4.
   */
  constructor(n: number, bs?: BitSet[], dp?: BitSet[]) {
    if (n > 1e4) throw new Error('n should be less than 1e4')
    if (bs === void 0) {
      bs = Array(n)
      for (let i = 0; i < n; i++) bs[i] = new BitSet(n)
    }
    this.n = n
    this._bs = bs
    this._dp = dp
  }

  ipow(k: number): BooleanSquareMatrixDense {
    const res = BooleanSquareMatrixDense.eye(this.n)
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

  imul(other: BooleanSquareMatrixDense): BooleanSquareMatrixDense {
    const n = this.n
    const res = new BooleanSquareMatrixDense(n)
    const step = 6 // !理论最优是logn,实际取6效果最好(n为5000时)
    this._initDpIfAbsent(step, n)
    const dp = this._dp!
    const otherBs = other._bs

    for (let left = 0, right = step; left < n; left = right, right += step) {
      if (right > n) right = n
      for (let state = 1; state < 1 << step; state++) {
        const bsf = BooleanSquareMatrixDense._trailingZeros32[state]
        if (left + bsf < n) {
          dp[state] = dp[state ^ (1 << bsf)].or(otherBs[left + bsf]) // !Xor => f2矩阵乘法
        } else {
          dp[state] = dp[state ^ (1 << bsf)]
        }
      }
      for (let i = 0, now = 0; i < n; i++, now = 0) {
        const thisBsI = this._bs[i]
        const resBsI = res._bs[i]
        // !这里是瓶颈 TODO:位运算优化
        for (let j = left; j < right; j++) now ^= +thisBsI.has(j) << (j - left)
        resBsI.ior(dp[now]) // !IXor => f2矩阵乘法
      }
    }

    const tmp = this._bs
    this._bs = res._bs
    res._bs = tmp
    return res
  }

  iadd(mat: BooleanSquareMatrixDense): BooleanSquareMatrixDense {
    for (let i = 0; i < this.n; i++) this._bs[i].ior(mat._bs[i])
    return this
  }

  /**
   * 求出邻接矩阵`mat`的传递闭包`(mat+I)^n`.
   * 稠密矩阵,2000*2000 => 1.4s.
   */
  transitiveClosure(): BooleanSquareMatrixDense {
    const n = this.n
    const trans = BooleanSquareMatrixDense.eye(n).iadd(this)
    trans.ipow(n)
    return trans
  }

  copy(): BooleanSquareMatrixDense {
    const newBs = Array(this.n)
    for (let i = 0; i < this.n; i++) newBs[i] = this._bs[i].copy()
    return new BooleanSquareMatrixDense(this.n, newBs, this._dp)
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
    const grid: Uint8Array[] = Array(this.n)
    for (let i = 0; i < this.n; i++) {
      grid[i] = new Uint8Array(this.n)
      for (let j = 0; j < this.n; j++) {
        grid[i][j] = this.get(i, j) ? 1 : 0
      }
    }
    // eslint-disable-next-line no-console
    console.table(grid)
  }

  private _initDpIfAbsent(step: number, n: number): void {
    if (this._dp) return
    const dp = Array(1 << step)
    for (let i = 0; i < dp.length; i++) dp[i] = new BitSet(n)
    this._dp = dp
  }
}

export { BooleanSquareMatrixDense }

if (require.main === module) {
  // ====================
  // 测试随机矩阵
  // BooleanSquareMatrixDense.mul: 995.97ms
  // BooleanSquareMatrixDense.transitiveClosure: 1.401s
  // ====================
  // 测试稀疏矩阵
  // BooleanSquareMatrixDense.mul: 933.747ms
  // BooleanSquareMatrixDense.transitiveClosure: 1.333s
  // ====================
  // 测试稠密矩阵
  // BooleanSquareMatrixDense.mul: 921.894ms
  // BooleanSquareMatrixDense.transitiveClosure: 1.408s

  const mat = new BooleanSquareMatrixDense(3)
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
    const mat = new BooleanSquareMatrixDense(N_5000)
    for (let i = 0; i < N_5000; i++) {
      for (let j = 0; j < N_5000; j++) {
        if (Math.random() < 0.5) mat.set(i, j, true)
      }
    }
    console.time('BooleanSquareMatrixDense.mul')
    BooleanSquareMatrixDense.mul(mat, mat)
    console.timeEnd('BooleanSquareMatrixDense.mul')

    // 2000*2000的传递闭包
    const N_2000 = 2000
    const mat2 = new BooleanSquareMatrixDense(N_2000)
    for (let i = 0; i < N_2000; i++) {
      for (let j = 0; j < N_2000; j++) {
        if (Math.random() < 0.5) mat2.set(i, j, true)
      }
    }
    console.time('BooleanSquareMatrixDense.transitiveClosure')
    mat2.transitiveClosure()
    console.timeEnd('BooleanSquareMatrixDense.transitiveClosure')
  }

  function testSparse(): void {
    console.log('='.repeat(20))
    console.log('测试稀疏矩阵')
    // !稀疏矩阵
    // 5000*5000的矩阵乘法
    const N_5000 = 5000
    const mat = new BooleanSquareMatrixDense(N_5000)
    for (let i = 0; i < N_5000; i++) {
      for (let j = 0; j < N_5000; j++) {
        if (Math.random() < 0.1) mat.set(i, j, true)
      }
    }
    console.time('BooleanSquareMatrixDense.mul')
    BooleanSquareMatrixDense.mul(mat, mat)
    console.timeEnd('BooleanSquareMatrixDense.mul')

    // 2000*2000的传递闭包
    const N_2000 = 2000
    const mat2 = new BooleanSquareMatrixDense(N_2000)
    for (let i = 0; i < N_2000; i++) {
      for (let j = 0; j < N_2000; j++) {
        if (Math.random() < 0.1) mat2.set(i, j, true)
      }
    }
    console.time('BooleanSquareMatrixDense.transitiveClosure')
    mat2.transitiveClosure()
    console.timeEnd('BooleanSquareMatrixDense.transitiveClosure')
  }

  function testDense(): void {
    console.log('='.repeat(20))
    console.log('测试稠密矩阵')
    // !稠密矩阵
    // 5000*5000的矩阵乘法
    const N_5000 = 5000
    const mat = new BooleanSquareMatrixDense(N_5000)
    for (let i = 0; i < N_5000; i++) {
      for (let j = 0; j < N_5000; j++) {
        mat.set(i, j, true)
      }
    }
    console.time('BooleanSquareMatrixDense.mul')
    BooleanSquareMatrixDense.mul(mat, mat)
    console.timeEnd('BooleanSquareMatrixDense.mul')

    // 2000*2000的传递闭包
    const N_2000 = 2000
    const mat2 = new BooleanSquareMatrixDense(N_2000)
    for (let i = 0; i < N_2000; i++) {
      for (let j = 0; j < N_2000; j++) {
        mat2.set(i, j, true)
      }
    }
    console.time('BooleanSquareMatrixDense.transitiveClosure')
    mat2.transitiveClosure()
    console.timeEnd('BooleanSquareMatrixDense.transitiveClosure')
  }
}

/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */
/* eslint-disable implicit-arrow-linebreak */

// 幂运算预处理.
// !注意js做带模乘法比较慢.
// 取模优化: const add = (x,y)=> (x += y) >= mod ? x - mod : x;

import { mulUint32 } from '../../数论/快速幂/mulUint32'
import { qpow } from '../../数论/快速幂/qpow'

type Public<C> = { [K in keyof C]: C[K] }

/**
 * 带有预处理的幺半群的幂运算.
 */
class PowerQuery<E> {
  private readonly _base: () => E
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _useCache: boolean
  private readonly _cacheLevel: number
  private readonly _data: E[][]

  /**
   * 幂运算预处理.
   * @param base 幂运算的基.
   * @param e 幺元.
   * @param op 结合法则.
   * @param cacheLevel 幂运算的 log 底数.默认为 16.
   */
  constructor(base: () => E, e: () => E, op: (a: E, b: E) => E, cacheLevel = 16) {
    this._base = base
    this._e = e
    this._op = op
    this._useCache = cacheLevel >= 2
    this._cacheLevel = cacheLevel
    this._data = this._useCache ? [this._makePow(base())] : []
  }

  /**
   * 计算 `base` 的 `exp` 次幂.
   */
  pow(exp: number): E {
    if (exp < 0) {
      throw new Error('exp must be non-negative')
    }

    if (!this._useCache) {
      return this.powWithOutCache(exp)
    }

    const logBase = this._cacheLevel
    const data = this._data
    let res = this._e()
    let level = 0
    while (exp) {
      const mod = exp % logBase
      exp = Math.floor(exp / logBase)
      if (data.length === level) {
        const last = data[level - 1]
        data.push(this._makePow(last[last.length - 1]))
      }
      res = this._op(res, data[level][mod])
      level++
    }

    return res
  }

  powWithOutCache(exp: number): E {
    if (exp < 0) {
      throw new Error('exp must be non-negative')
    }

    let res = this._e()
    let base = this._base()
    while (exp) {
      if (exp & 1) {
        res = this._op(res, base)
      }
      base = this._op(base, base)
      exp = Math.floor(exp / 2)
    }

    return res
  }

  private _makePow(e: E): E[] {
    const res = Array<E>(this._cacheLevel + 1)
    res[0] = this._e()
    for (let i = 0; i < this._cacheLevel; i++) {
      res[i + 1] = this._op(res[i], e)
    }
    return res
  }
}

/**
 * 带有预处理的快速幂.
 */
const usePowerQueryQPow = (base: number, mod = 1e9 + 7, cacheLevel = 16) =>
  new PowerQuery(
    () => base,
    () => 1,
    (a, b) => qpow(a, b, mod),
    cacheLevel
  )

/**
 * 带有预处理的矩阵快速幂.
 * 内部使用`uint32`类型数组加速.
 */
const usePowerQueryMatQPow = (base: number[][], mod = 1e9 + 7, cacheLevel = 16): Public<PowerQuery<number[][]>> => {
  if (base.length !== base[0].length) throw new Error('base must be a square matrix')

  const n = base.length
  const pq = new PowerQuery(
    () => compress(base),
    () => eye(n),
    (a, b) => mul(a, b, mod),
    cacheLevel
  )

  return {
    pow(exp: number) {
      const res = pq.pow(exp)
      return decompress(res)
    },
    powWithOutCache(exp: number) {
      const res = pq.powWithOutCache(exp)
      return decompress(res)
    }
  }

  function compress(mat: number[][]): Uint32Array {
    const res = new Uint32Array(n * n)
    for (let i = 0; i < n; i++) {
      const row = mat[i]
      for (let j = 0; j < n; j++) {
        res[i * n + j] = row[j]
      }
    }
    return res
  }

  function decompress(compressed: Uint32Array): number[][] {
    const res: number[][] = Array(n)
    for (let i = 0; i < n; i++) {
      const row = Array(n)
      for (let j = 0; j < n; j++) {
        row[j] = compressed[i * n + j]
      }
      res[i] = row
    }
    return res
  }

  function eye(size: number): Uint32Array {
    const res = new Uint32Array(size * size)
    for (let i = 0; i < size; i++) {
      res[i * size + i] = 1
    }
    return res
  }

  function mul(mat1: Uint32Array, mat2: Uint32Array, mod = 1e9 + 7): Uint32Array {
    const res = new Uint32Array(n * n)
    for (let i = 0; i < n; i++) {
      for (let k = 0; k < n; k++) {
        for (let j = 0; j < n; j++) {
          const pos = i * n + j
          const tmp = res[pos] + mulUint32(mat1[i * n + k], mat2[k * n + j], mod)
          res[pos] = tmp > mod ? tmp - mod : tmp
        }
      }
    }
    return res
  }
}

export { PowerQuery, usePowerQueryQPow, usePowerQueryMatQPow }

if (require.main === module) {
  const MOD = 1e9 + 7
  // const pow2 = usePowerQueryQPow(2, MOD)

  // // 2550. 猴子碰撞的方法数
  // // https://leetcode.cn/problems/count-collisions-of-monkeys-on-a-polygon/
  // function monkeyMove(n: number): number {
  //   const res = pow2.pow(n) - 2
  //   return res < 0 ? res + MOD : res
  // }

  // https://leetcode.cn/contest/ccbft-2021fall/problems/lSjqMF/
  function electricityExperiment(row: number, col: number, position: [x: number, y: number][]) {
    const T = Array.from({ length: row }, () => Array(row).fill(0))
    for (let r = 0; r < row; r++) {
      T[r][r] = 1
      if (r !== 0) T[r][r - 1] = 1
      if (r !== row - 1) T[r][r + 1] = 1
    }

    const Q = usePowerQueryMatQPow(T, MOD)

    const n = position.length
    position.sort((a, b) => a[1] - b[1])
    for (let i = 0; i < n - 1; i++) {
      const [row1, col1] = position[i]
      const [row2, col2] = position[i + 1]
      const rowDiff = Math.abs(row1 - row2)
      const colDiff = Math.abs(col1 - col2)
      if (rowDiff > colDiff) return 0
    }

    let res = 1
    for (let i = 0; i < n - 1; i++) {
      const [row1, col1] = position[i]
      const [row2, col2] = position[i + 1]
      const colDiff = Math.abs(col1 - col2)
      res = mulUint32(res, cal(row1, row2, colDiff), MOD)
    }

    return res

    function cal(row1: number, row2: number, k: number) {
      const resT = Q.pow(k)
      return resT[row1][row2]
    }
  }

  //   5
  // 6
  // [[1, 3], [3, 2], [4, 1]]
  console.log(
    electricityExperiment(5, 6, [
      [1, 3],
      [3, 2],
      [4, 1]
    ])
  )
}

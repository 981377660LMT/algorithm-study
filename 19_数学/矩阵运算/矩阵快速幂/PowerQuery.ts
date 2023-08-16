// 幂运算预处理

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
        data.push(this._makePow(data[level - 1][logBase - 1]))
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
const usePowerQueryQPow = (base: number, mod = 1e9 + 7, cacheLevel = 16) => {
  new PowerQuery<number>(
    () => base,
    () => 1,
    (a, b) => a * b,
    cacheLevel
  )
}

/**
 * 带有预处理的矩阵快速幂.
 * 内部使用`uint32`类型数组加速.
 */
const usePowerQueryMatQPow = (
  base: ArrayLike<ArrayLike<number>>,
  mod = 1e9 + 7,
  cacheLevel = 16
) => {}

export { PowerQuery, usePowerQueryQPow, usePowerQueryMatQPow }

if (require.main === module) {
//   MOD = int(1e9 + 7)


// class Solution:
//     def numTilings(self, n: int) -> int:
//         init = [[5], [2], [1], [0]]
//         if n <= 3:
//             return init[~n][0]
//         T = [[2, 0, 1, 0], [1, 0, 0, 0], [0, 1, 0, 0], [0, 0, 1, 0]]
//         resT = matqpow1(T, n - 3, MOD)
//         res = mul(resT, init, MOD)
//         return res[0][0]

  // 790. 多米诺和托米诺平铺
  // https://leetcode.cn/problems/domino-and-tromino-tiling/
  function numTilings(n: number): number {}
}

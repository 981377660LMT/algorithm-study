// 光速幂
// 已知base和模数mod，求base^n 模mod
// !O(sqrt(maxN))预处理,O(1)查询

class FastPow<E> {
  private readonly _max: number
  private readonly _divData: E[]
  private readonly _modData: E[]
  private readonly _mul: (a: E, b: E) => E

  /**
   * 光速幂.
   * @param base 底数
   * @param maxN 最大幂次
   * @param e 幺元
   * @param op 结合法则
   */
  constructor(base: E, maxN: number, e: () => E, op: (a: E, b: E) => E) {
    const max = Math.ceil(Math.sqrt(maxN))
    const divData = Array<E>(max)
    const modData = Array<E>(max)

    let cur1 = e()
    for (let i = 0; i <= max; i++) {
      modData[i] = cur1
      cur1 = op(cur1, base)
    }

    let cur2 = e()
    const last = modData[max]
    for (let i = 0; i <= max; i++) {
      divData[i] = cur2
      cur2 = op(cur2, last)
    }

    this._max = max
    this._divData = divData
    this._modData = modData
    this._mul = op
  }

  /**
   * 计算 `base` 的 `exp` 次幂.
   * @param exp 0 <= exp <= maxN
   */
  pow(exp: number): E {
    return this._mul(this._divData[Math.floor(exp / this._max)], this._modData[exp % this._max])
  }
}

export { FastPow }

if (require.main === module) {
  const MOD = BigInt(1e9 + 7)
  const pow2 = new FastPow<bigint>(
    2n,
    1e9,
    () => 1n,
    (a, b) => (a * b) % MOD
  )

  // 2550. 猴子碰撞的方法数
  // https://leetcode.cn/problems/count-collisions-of-monkeys-on-a-polygon/
  // eslint-disable-next-line no-inner-declarations
  function monkeyMove(n: number): number {
    const res = pow2.pow(n) - 2n
    return Number(res < 0 ? res + MOD : res)
  }
}

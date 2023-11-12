/* eslint-disable no-inner-declarations */

/**
 * 倍增维护`n`个状态的转移编号和值.
 * 每个状态对应一个`编号(0-n-1)`和`值(幺半群)`.
 * 从状态`a`转移一次到达状态`b`，状态`b`对应的值与`边权`进行结合运算.
 */
class Doubling<E> {
  private readonly _n: number
  private readonly _log: number
  private readonly _e: () => E
  private readonly _op: (e1: E, e2: E) => E
  private readonly _to: Int32Array
  private readonly _dp: E[]
  private _isPrepared = false

  /**
   * @param n 状态的数量.
   * @param maxStep 最大的转移次数.
   * @param e 转移边权的单位元.
   * @param op 转移边权的结合律运算.
   */
  constructor(n: number, maxStep: number, e: () => E, op: (e1: E, e2: E) => E) {
    this._n = n
    this._log = 1
    while (2 ** this._log <= maxStep) this._log++
    this._e = e
    this._op = op

    const size = n * this._log
    this._to = new Int32Array(size)
    this._dp = Array<E>(size)
    for (let i = 0; i < size; i++) {
      this._to[i] = -1
      this._dp[i] = e()
    }
  }

  /**
   * 初始状态(leaves):从 `from` 状态到 `to` 状态，边权为 `weight`.
   * 0 <= from, to < n.
   */
  add(from: number, to: number, weight: E): void {
    if (this._isPrepared) throw new Error('Doubling is prepared')
    if (to < -1 || to >= this._n) throw new RangeError('to is out of range')
    this._to[from] = to
    this._dp[from] = weight
  }

  build(): void {
    if (this._isPrepared) return
    this._isPrepared = true
    const n = this._n
    for (let k = 0; k < this._log - 1; k++) {
      for (let v = 0; v < n; v++) {
        const w = this._to[k * n + v]
        const next = (k + 1) * n + v
        if (w === -1) {
          this._to[next] = -1
          this._dp[next] = this._dp[k * n + v]
          continue
        }
        this._to[next] = this._to[k * n + w]
        this._dp[next] = this._op(this._dp[k * n + v], this._dp[k * n + w])
      }
    }
  }

  /**
   * 从 `from` 状态开始，执行 `step` 次操作，返回最终状态的编号和值.
   * 0 <= from < n.
   * 如果最终状态不存在，返回 [-1, e()].
   */
  jump(from: number, step: number): { to: number; value: E } {
    if (!this._isPrepared) this.build()
    if (step >= 2 ** this._log) throw new RangeError('step is over max step')
    let value = this._e()
    let to = from
    for (let k = 0; k < this._log; k++) {
      if (to === -1) break
      const div = 2 ** k
      if (Math.floor(step / div) & 1) {
        const pos = k * this._n + to
        value = this._op(value, this._dp[pos])
        to = this._to[pos]
      }
    }
    return { to, value }
  }

  /**
   * 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step`以及对应的状态编号和值.
   * 0 <= from < n.
   */
  maxStep(from: number, check: (e: E) => boolean): { step: number; to: number; value: E } {
    if (!this._isPrepared) this.build()
    let res = this._e()
    let step = 0
    for (let k = this._log - 1; ~k; k--) {
      const pos = k * this._n + from
      const to = this._to[pos]
      const next = this._op(res, this._dp[pos])
      if (check(next)) {
        step += 2 ** k
        from = to
        res = next
      }
    }
    return { step, to: from, value: res }
  }
}

export { Doubling }

if (require.main === module) {
  class TreeAncestor {
    private _db: Doubling<number>
    constructor(n: number, parent: number[]) {
      this._db = new Doubling<number>(
        n,
        1e5 + 10,
        () => -1,
        () => 1
      )
      for (let i = 1; i < n; i++) this._db.add(i, parent[i], 1)
      this._db.build()
    }

    getKthAncestor(node: number, k: number): number {
      return this._db.jump(node, k).to
    }
  }

  /**
   * Your TreeAncestor object will be instantiated and called as such:
   * var obj = new TreeAncestor(n, parent)
   * var param_1 = obj.getKthAncestor(node,k)
   */

  // 8027. 在传球游戏中最大化函数值
  // https://leetcode.cn/problems/maximize-value-of-function-in-a-ball-passing-game/
  function getMaxFunctionValue(receiver: number[], k: number): number {
    const n = receiver.length
    const db = new Doubling(
      n,
      k + 10,
      () => 0,
      (a, b) => a + b
    )
    for (let i = 0; i < n; i++) db.add(i, receiver[i], i)
    db.build()

    let res = 0
    for (let i = 0; i < n; i++) {
      const { value } = db.jump(i, k + 1)
      res = Math.max(res, value)
    }
    return res
  }
}

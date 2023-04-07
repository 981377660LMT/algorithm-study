/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * 周期函数的幂(k次转移后的状态).
 */
class PeriodicFunctionPower<State extends number | string> {
  readonly offset: number
  readonly cycle: number
  private readonly _data: State[] = []
  private readonly _used: Map<State, number> = new Map()

  /**
   * @param s0 初始状态.
   * @param next 状态转移函数.
   */
  constructor(s0: State, next: (cur: State) => State) {
    while (!this._used.has(s0)) {
      this._used.set(s0, this._data.length)
      this._data.push(s0)
      s0 = next(s0)
    }
    this.offset = this._used.get(s0)!
    this.cycle = this._data.length - this.offset
  }

  /**
   * 查询k次转移后的状态(第k项).
   * k>=0.
   */
  query(k: number): State {
    const index = k < this.offset ? k : ((k - this.offset) % this.cycle) + this.offset
    return this._data[index]
  }
}

if (require.main === module) {
  // 957. N 天后的牢房
  // eslint-disable-next-line no-inner-declarations
  function prisonAfterNDays(cells: number[], n: number): number[] {
    const s0 = cells.join('')
    const F = new PeriodicFunctionPower(s0, (cur: string) => {
      const next = Array(8).fill('0')
      for (let i = 1; i < 7; i++) {
        if (cur[i - 1] === cur[i + 1]) next[i] = '1'
      }
      return next.join('')
    })

    return F.query(n).split('').map(Number)
  }
}

export { PeriodicFunctionPower }

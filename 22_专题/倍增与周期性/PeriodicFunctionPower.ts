/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * 周期函数的幂(k次转移后的状态).
 */
class PeriodicFunctionPower<State extends number | string> {
  readonly cycleStart: number
  readonly cycleLength: number
  readonly preCycle: State[]
  readonly cycle: State[]

  /**
   * @param s0 初始状态.
   * @param next 状态转移函数.
   */
  constructor(s0: State, next: (cur: State) => State) {
    const visited = new Map<State, number>()
    const history: State[] = []
    let now = s0
    while (!visited.has(now)) {
      visited.set(now, history.length)
      history.push(now)
      now = next(now)
    }
    this.cycleStart = visited.get(now)!
    this.cycleLength = history.length - this.cycleStart
    this.preCycle = history.slice(0, this.cycleStart)
    this.cycle = history.slice(this.cycleStart)
  }

  /**
   * 查询k次转移后的状态(第k项).
   * k>=0.
   */
  query(k: number): State {
    if (k < this.cycleStart) return this.preCycle[k]
    k -= this.cycleStart
    k %= this.cycleLength
    return this.cycle[k]
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

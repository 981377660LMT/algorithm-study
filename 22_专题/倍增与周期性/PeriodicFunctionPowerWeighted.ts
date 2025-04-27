class PeriodicFunctionPowerWeighted<S extends number | string, E> {
  readonly cycleStart: number
  readonly cycleLength: number
  readonly preStates: S[]
  readonly cycleStates: S[]
  readonly preWeights: E[]
  readonly cycleWeights: E[]
  readonly cycleSum: E

  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _pow: (a: E, b: number) => E

  constructor(
    e: () => E,
    op: (a: E, b: E) => E,
    pow: (a: E, b: number) => E,
    s0: S,
    next: (cur: S) => [S, E]
  ) {
    this._e = e
    this._op = op
    this._pow = pow

    const visited = new Map<S, number>()
    const states: S[] = []
    const weights: E[] = []
    let now = s0

    while (!visited.has(now)) {
      visited.set(now, states.length)
      states.push(now)
      const { 0: nxt, 1: w } = next(now)
      weights.push(w)
      now = nxt
    }

    this.cycleStart = visited.get(now)!
    this.cycleLength = states.length - this.cycleStart
    this.preStates = states.slice(0, this.cycleStart)
    this.cycleStates = states.slice(this.cycleStart)

    this.preWeights = Array(this.cycleStart + 1)
    let acc = this._e()
    this.preWeights[0] = acc
    for (let i = 0; i < this.cycleStart; i++) {
      acc = this._op(acc, weights[i])
      this.preWeights[i + 1] = acc
    }

    this.cycleWeights = Array(this.cycleLength + 1)
    acc = this._e()
    this.cycleWeights[0] = acc
    for (let i = 0; i < this.cycleLength; i++) {
      acc = this._op(acc, weights[this.cycleStart + i])
      this.cycleWeights[i + 1] = acc
    }
    this.cycleSum = acc
  }

  /**
   * 查询k次转移后的状态和权重.
   * k>=0.
   */
  query(k: number): [S, E] {
    if (k < this.cycleStart) return [this.preStates[k], this.preWeights[k]]
    const steps = k - this.cycleStart
    const d = Math.floor(steps / this.cycleLength)
    const r = steps % this.cycleLength
    let w = this.preWeights[this.cycleStart]
    if (d > 0) {
      w = this._op(w, this._pow(this.cycleSum, d))
    }
    w = this._op(w, this.cycleWeights[r])
    return [this.cycleStates[r], w]
  }
}

export { PeriodicFunctionPowerWeighted }

type Pair = { value: number; count: number }

class ClampableStack {
  private readonly _clampMin: boolean
  private _total = 0
  private _count = 0
  private _stack: Pair[] = []

  /**
   * @param clampMin 为true时，支持容器内所有数与x取min；为false时，支持容器内所有数与x取max.
   */
  constructor(clampMin: boolean) {
    this._clampMin = clampMin
  }

  addAndClamp(x: number): void {
    let newCount = 1
    if (this._clampMin) {
      while (this._stack.length > 0) {
        const top = this._stack[this._stack.length - 1]
        if (top.value < x) break
        this._stack.pop()
        this._total -= top.value * top.count
        newCount += top.count
      }
    } else {
      while (this._stack.length > 0) {
        const top = this._stack[this._stack.length - 1]
        if (top.value > x) break
        this._stack.pop()
        this._total -= top.value * top.count
        newCount += top.count
      }
    }
    this._total += x * newCount
    this._count++
    this._stack.push({ value: x, count: newCount })
  }

  sum(): number {
    return this._total
  }

  get size(): number {
    return this._count
  }
}

export { ClampableStack }

if (require.main === module) {
  const S1 = new ClampableStack(true)
  S1.addAndClamp(1)
  S1.addAndClamp(2)
  S1.addAndClamp(1)
  console.assert(S1.sum() === 3)
  const S2 = new ClampableStack(false)
  S2.addAndClamp(1)
  S2.addAndClamp(2)
  S2.addAndClamp(1)
  console.assert(S2.sum() === 5)
  console.log('clamped stack passed')
}

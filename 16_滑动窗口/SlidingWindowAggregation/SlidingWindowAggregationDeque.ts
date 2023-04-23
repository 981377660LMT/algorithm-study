// # https://judge.yosupo.jp/submission/118808

class SlidingWindowAggregationDeque<E> {
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _leftVal: E[] = []
  private readonly _leftSum: E[] = []
  private readonly _rightVal: E[] = []
  private readonly _rightSum: E[] = []

  constructor(e: () => E, op: (a: E, b: E) => E) {
    this._e = e
    this._op = op
  }

  query(): E {
    if (!this._leftSum.length || !this._rightSum.length) {
      if (!this._leftSum.length && !this._rightSum.length) return this._e()
      if (!this._leftSum.length) return this._rightSum[this._rightSum.length - 1]
      return this._leftSum[this._leftSum.length - 1]
    }

    return this._op(
      this._leftSum[this._leftSum.length - 1],
      this._rightSum[this._rightSum.length - 1]
    )
  }

  append(val: E) {
    if (!this._rightSum.length) {
      this._rightVal.push(val)
      this._rightSum.push(val)
      return
    }

    this._rightVal.push(val)
    this._rightSum.push(this._op(this._rightSum[this._rightSum.length - 1], val))
  }

  appendLeft(val: E) {
    if (!this._leftSum.length) {
      this._leftVal.push(val)
      this._leftSum.push(val)
      return
    }

    this._leftVal.push(val)
    this._leftSum.push(this._op(val, this._leftSum[this._leftSum.length - 1]))
  }

  pop() {
    if (!this._rightSum.length) {
      const ln = this._leftSum.length >>> 1
      const rn = this._leftSum.length - ln
      const lv = []
      this._leftSum.length = 0
      for (let i = 0; i < ln; i++) {
        lv.push(this._leftVal.pop())
      }

      for (let i = 0; i < rn; i++) {
        const x = this._leftVal.pop()!
        this._rightVal.push(x)
        if (!this._rightSum.length) {
          this._rightSum.push(x)
        } else {
          this._rightSum.push(this._op(this._rightSum[this._rightSum.length - 1], x))
        }
      }

      for (let i = 0; i < ln; i++) {
        const x = lv.pop()!
        this._leftVal.push(x)
        if (!this._leftSum.length) {
          this._leftSum.push(x)
        } else {
          this._leftSum.push(this._op(x, this._leftSum[this._leftSum.length - 1]))
        }
      }
    }

    this._rightVal.pop()
    this._rightSum.pop()
  }

  popLeft() {
    if (!this._leftSum.length) {
      const rn = this._rightSum.length >>> 1
      const ln = this._rightSum.length - rn
      const rv = []
      this._rightSum.length = 0
      for (let i = 0; i < rn; i++) {
        rv.push(this._rightVal.pop())
      }

      for (let i = 0; i < ln; i++) {
        const x = this._rightVal.pop()!
        this._leftVal.push(x)
        if (!this._leftSum.length) {
          this._leftSum.push(x)
        } else {
          this._leftSum.push(this._op(x, this._leftSum[this._leftSum.length - 1]))
        }
      }

      for (let i = 0; i < rn; i++) {
        const x = rv.pop()!
        this._rightVal.push(x)
        if (!this._rightSum.length) {
          this._rightSum.push(x)
        } else {
          this._rightSum.push(this._op(this._rightSum[this._rightSum.length - 1], x))
        }
      }
    }

    this._leftVal.pop()
    this._leftSum.pop()
  }

  get length(): number {
    return this._leftSum.length + this._rightSum.length
  }
}

export { SlidingWindowAggregationDeque }

if (require.main === module) {
  const maxQueue = new SlidingWindowAggregationDeque<number>(() => 0, Math.max)

  maxQueue.append(1)
  maxQueue.append(2)
  maxQueue.append(3)
  maxQueue.append(4)
  console.log(maxQueue.query()) // 4
  maxQueue.pop()
  console.log(maxQueue.query()) // 3
  maxQueue.popLeft()
  console.log(maxQueue.query()) // 2
  maxQueue.appendLeft(5)
  console.log(maxQueue.query()) // 5
}

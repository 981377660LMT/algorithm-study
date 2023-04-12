/* eslint-disable no-inner-declarations */

class QueueAggregation<T> {
  private readonly _front: _StackAggregation<T>
  private readonly _back: _StackAggregation<T>
  private readonly _op: (a: T, b: T) => T

  constructor(e: () => T, op: (a: T, b: T) => T) {
    this._front = new _StackAggregation(e, op)
    this._back = new _StackAggregation(e, op)
    this._op = op
  }

  append(x: T): void {
    if (this.size === 0) {
      this._front.push(x)
    } else {
      this._back.push(x)
    }
  }

  popleft(): void {
    this._front.pop()
    if (this._front.size === 0) {
      while (this._back.size > 0) {
        this._front.push(this._back.top()!)
        this._back.pop()
      }
    }
  }

  query(): T {
    return this._op(this._front.query(), this._back.query())
  }

  get size(): number {
    return this._front.size + this._back.size
  }
}

class _StackAggregation<T> {
  private readonly _values: T[] = []
  private readonly _folds: T[] = []
  private readonly _e: () => T
  private readonly _op: (a: T, b: T) => T

  constructor(e: () => T, op: (a: T, b: T) => T) {
    this._e = e
    this._op = op
  }

  /**
   * 查询栈中幺半群的和.
   */
  query(): T {
    return this._values.length ? this._folds[this._folds.length - 1] : this._e()
  }

  push(x: T): void {
    this._values.push(x)
    this._folds.push(
      this._values.length === 1 ? x : this._op(this._folds[this._folds.length - 1], x)
    )
  }

  pop(): T | undefined {
    this._folds.pop()
    return this._values.pop()
  }

  top(): T | undefined {
    return this._values[this._values.length - 1]
  }

  get size(): number {
    return this._values.length
  }
}

if (require.main === module) {
  const INF = 2e15
  // https://leetcode-cn.com/problems/sliding-window-maximum/
  function maxSlidingWindow(nums: number[], k: number): number[] {
    const n = nums.length
    const ans: number[] = []
    const window = new QueueAggregation(() => -INF, Math.max)
    for (let i = 0; i < n; i++) {
      window.append(nums[i])
      if (i >= k - 1) {
        ans.push(window.query())
        window.popleft()
      }
    }
    return ans
  }
}

export { QueueAggregation }

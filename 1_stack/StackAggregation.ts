// 实现最大栈/最小栈

/**
 * 维护幺半群的栈.
 */
class StackAggregation<T> {
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

export { StackAggregation }

if (require.main === module) {
  const INF = 2e15
  // https://leetcode.cn/problems/min-stack/
  class MinStack {
    private readonly _stack = new StackAggregation(() => INF, Math.min)

    push(val: number): void {
      this._stack.push(val)
    }

    pop(): void {
      this._stack.pop()
    }

    top(): number {
      return this._stack.top()!
    }

    getMin(): number {
      return this._stack.query()
    }
  }
}

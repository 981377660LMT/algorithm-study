/* eslint-disable max-len */

// BITGroup
// BITGroupRangeAddRangeSum

class BITGroup<S> {
  private readonly _n: number
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S
  private readonly _inv: (a: S) => S
  private _total: S
  private readonly _data: S[]

  constructor(
    n: number,
    group: {
      e: () => S
      op: (a: S, b: S) => S
      inv: (a: S) => S
    } & ThisType<void>,
    f?: (index: number) => S
  ) {
    this._n = n
    const { e, op, inv } = group
    this._e = e
    this._op = op
    this._inv = inv
    this._total = e()
    if (!f) {
      this._data = Array(n)
      for (let i = 0; i < n; i++) this._data[i] = e()
    } else {
      this._data = Array(n)
      for (let i = 0; i < n; i++) {
        this._data[i] = f(i)
        this._total = op(this._total, this._data[i])
      }
      for (let i = 1; i <= n; i++) {
        const j = i + (i & -i)
        if (j <= n) this._data[j - 1] = op(this._data[j - 1], this._data[i - 1])
      }
    }
  }

  update(i: number, x: S): void {
    this._total = this._op(this._total, x)
    for (i++; i <= this._n; i += i & -i) this._data[i - 1] = this._op(this._data[i - 1], x)
  }

  queryAll(): S {
    return this._total
  }

  queryPrefix(end: number): S {
    if (end > this._n) end = this._n
    let res = this._e()
    while (end > 0) {
      res = this._op(res, this._data[end - 1])
      end &= end - 1
    }
    return res
  }

  queryRange(start: number, end: number): S {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start === 0) return this.queryPrefix(end)
    if (start > end) return this._e()
    let pos = this._e()
    let neg = this._e()
    while (end > start) {
      pos = this._op(pos, this._data[end - 1])
      end &= end - 1
    }
    while (start > end) {
      neg = this._op(neg, this._data[start - 1])
      start &= start - 1
    }
    return this._op(pos, this._inv!(neg))
  }

  maxRight(predicate: (preSum: S, end: number) => boolean): number {
    let i = 0
    let cur = this._e()
    let k = 1
    while (2 * k <= this._n) k *= 2
    while (k > 0) {
      if (i + k - 1 < this._data.length) {
        const t = this._op(cur, this._data[i + k - 1])
        if (predicate(t, i + k)) {
          i += k
          cur = t
        }
      }
      k >>>= 1
    }
    return i
  }

  toString(): string {
    return `BITGroup: [${Array.from({ length: this._n }, (_, i) => this.queryRange(i, i + 1)).join(', ')}]`
  }
}

class BITGroupRangeAddRangeSum<S> {
  private readonly _n: number
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S
  private readonly _inv: (a: S) => S
  private readonly _mul: (a: S, n: number) => S
  private readonly _bit0: BITGroup<S>
  private readonly _bit1: BITGroup<S>

  constructor(
    n: number,
    group: { e: () => S; op: (a: S, b: S) => S; inv: (a: S) => S; mul: (a: S, n: number) => S } & ThisType<void>,
    f?: (index: number) => S
  ) {
    this._n = n
    const { e, op, inv, mul } = group
    this._e = e
    this._op = op
    this._inv = inv
    this._mul = mul
    this._bit0 = f ? new BITGroup(n, group, f) : new BITGroup(n, group)
    this._bit1 = new BITGroup(n, group)
  }

  update(index: number, x: S): void {
    this._bit0.update(index, x)
  }

  updateRange(start: number, end: number, x: S): void {
    this._bit0.update(start, this._mul(x, -start))
    this._bit0.update(end, this._mul(x, end))
    this._bit1.update(start, x)
    this._bit1.update(end, this._inv(x))
  }

  queryRange(start: number, end: number): S {
    const rightRes = this._op(this._mul(this._bit1.queryPrefix(end), end), this._bit0.queryPrefix(end))
    const leftRes = this._op(this._mul(this._bit1.queryPrefix(start), start), this._bit0.queryPrefix(start))
    return this._op(this._inv(leftRes), rightRes)
  }

  toString(): string {
    return `BITGroupRangeAddRangeSum: [${Array.from({ length: this._n }, (_, i) => this.queryRange(i, i + 1)).join(', ')}]`
  }
}

export { BITGroup, BITGroupRangeAddRangeSum }

if (require.main === module) {
  const bitGroup = new BITGroup(
    10,
    {
      e: () => 0,
      op: (a, b) => a + b,
      inv: a => -a
    },
    i => i
  )

  console.log(bitGroup.toString())

  const bitGroupRangeAddRangeSum = new BITGroupRangeAddRangeSum(
    10,
    {
      e: () => 0,
      op: (a, b) => a + b,
      inv: a => -a,
      mul: (a, n) => a * n
    },
    i => i
  )
  bitGroupRangeAddRangeSum.updateRange(0, 10, 1)
  console.log(bitGroupRangeAddRangeSum.toString())
  console.log(bitGroupRangeAddRangeSum.queryRange(0, 5))
}

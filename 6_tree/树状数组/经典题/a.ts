const INF = 2e15
function beautySum(s: string): number {
  let res = 0
  for (let i = 0; i < s.length; i++) {
    const minSeg = new BITMonoidArray(26, () => INF, Math.min)
    const maxSeg = new BITMonoidArray(26, () => 0, Math.max)
    const counterSeg = new BITMonoidArray(
      26,
      () => 0,
      (a, b) => a + b
    )

    for (let j = i; j < s.length; j++) {
      const c = s.charCodeAt(j) - 97
      minSeg.set(c, counterSeg.get(c) + 1)
      maxSeg.set(c, counterSeg.get(c) + 1)
      counterSeg.update(c, 1)

      const cand1 = minSeg.queryPrefix(26) - maxSeg.queryPrefix(26)
      const cand2 = minSeg.queryRange(0, 26) - maxSeg.queryRange(0, 26)
      if (cand1 !== cand2) throw new Error('error')

      res += cand1
    }
  }
  return res
}

/**
 * 维护幺半群的树状数组.
 * 支持单点更新,单点修改,前缀查询,区间查询.
 */
class BITMonoidArray<E = number> {
  private readonly _n: number
  private readonly _data: E[]
  private readonly _sum: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E

  constructor(nOrArr: number | ArrayLike<E>, e: () => E, op: (a: E, b: E) => E) {
    const n = typeof nOrArr === 'number' ? nOrArr : nOrArr.length
    this._n = n
    this._e = e
    this._op = op
    this._data = Array(n + 1)
    this._sum = Array(n + 1)
    for (let i = 0; i < n + 1; i++) {
      this._data[i] = e()
      this._sum[i] = e()
    }
    if (typeof nOrArr !== 'number') this.build(nOrArr)
  }

  /**
   * 单点修改,时间复杂度O(log^2 n).
   * 0<=index<n.
   */
  set(index: number, value: E): void {
    index++
    this._data[index] = value
    for (; index <= this._n; index += index & -index) {
      this._sum[index] = this._data[index]
      for (let i = 1; i < (index & -index); i <<= 1) {
        this._sum[index] = this._op(this._sum[index], this._sum[index - i])
      }
    }
  }

  get(index: number): E {
    return this.queryRange(index, index + 1)
  }

  /**
   * 单点更新,时间复杂度O(log n).
   * 0<=index<n.
   */
  update(index: number, value: E): void {
    index++
    this._data[index] = this._op(this._data[index], value)
    for (; index <= this._n; index += index & -index) {
      this._sum[index] = this._op(this._sum[index], value)
    }
  }

  /**
   * 查询前缀`[0,right)`聚合值,时间复杂度O(log n).
   * 0<=right<=n.
   */
  queryPrefix(right: number): E {
    if (right > this._n) right = this._n
    let res = this._e()
    for (; right > 0; right &= right - 1) res = this._op(res, this._sum[right])
    return res
  }

  /**
   * 查询区间`[left,right)`聚合值,时间复杂度O(log^2 n).
   * 0<=left<=right<=n.
   */
  queryRange(left: number, right: number): E {
    if (right > this._n) right = this._n
    left++
    let res = this._e()
    while (right >= left) {
      if ((right & (right - 1)) >= left - 1) {
        res = this._op(res, this._sum[right])
        right &= right - 1
      } else {
        res = this._op(res, this._data[right])
        right--
      }
    }
    return res
  }

  build(arr: ArrayLike<E>): void {
    if (arr.length !== this._n) throw new RangeError(`arr length must be equal to ${this._n}`)
    for (let i = 1; i <= this._n; i++) {
      this._data[i] = arr[i - 1]
      this._sum[i] = arr[i - 1]
      for (let j = 1; j < (i & -i); j <<= 1) {
        this._sum[i] = this._op(this._sum[i], this._sum[i - j])
      }
    }
  }

  toString(): string {
    const res: E[] = []
    for (let i = 0; i < this._n; i++) {
      res.push(this.queryRange(i, i + 1))
    }
    return `BITMonoid{${res.join(',')}}`
  }
}

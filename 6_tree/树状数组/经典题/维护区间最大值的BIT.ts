// https://leetcode.cn/circle/discuss/9n7Hnx/

const INF = 2e15

// TODO
// !Not Verified

/**
 * 维护幺半群的树状数组.
 * !1-indexed.
 */
class BITMonoid<E> {
  private readonly _tree: number[]
  private readonly _max: number[]
  private readonly _n: number

  constructor(n: number, e: () => E, op: (a: E, b: E) => E) {
    this._tree = Array(n + 1).fill(0)
    this._max = Array(n + 1).fill(0)
    this._n = n
  }

  // O(log^2 n)
  set(index: number, value: number): void {
    this._max[index] = value
    for (; index <= this._n; index += index & -index) {
      this._tree[index] = value
      for (let i = 1; i < (index & -index); i <<= 1) {
        this._tree[index] = Math.max(this._tree[index], this._tree[index - i])
      }
    }
  }

  // O(logn)
  update(index: number, value: number): void {}

  // O(logn)
  queryPrefix(right: number): number {
    let res = -INF
    for (; right > 0; right -= right & -right) {
      res = Math.max(res, this._tree[right])
    }
    return res
  }

  // O(log^2 n)
  queryRange(left: number, right: number): number {
    let x = right
    let res = -INF
    while (x >= left) {
      if (x - (x & -x) >= left - 1) {
        res = Math.max(res, this._tree[x])
        x -= x & -x
      } else {
        res = Math.max(res, this._max[x])
        x--
      }
    }
    return res
  }
}

export {}

if (require.main === module) {
  const bit = new BITMonoid(10, () => 0, Math.max)

  console.log(bit.toString())
  bit.set(1, 3)
  bit.set(2, 4)
  bit.set(3, 5)
  console.log(bit.toString())

  console.log(bit.toString())
  console.log(bit.queryRange(1, 3))
  console.log(bit.toString())
}

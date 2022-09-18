// 给一个长度为 n 的数列和一个正整数 m，1 < k <= n < 5e5，
// 每次可以在数列中选连续 k 个数减 1，但是不能把任何数减成负数
// 问最多能减多少次
// !贪心
// 枚举左端点i，取i到i+m-1的最小值，然后全都减

import assert from 'assert'

const INF = 2e15

/**
 * 线段树区间叠加最小值RMQ
 *
 * !叠加更新可以省去isLazy数组
 * !如果查询超出范围 返回INF
 */
class MinSegmentTree {
  private readonly _tree: number[]
  private readonly _lazyValue: number[]
  private readonly _size: number

  /**
   * @param sizeOrNums 数组长度或数组
   */
  constructor(sizeOrNums: number | readonly number[]) {
    this._size = typeof sizeOrNums === 'number' ? sizeOrNums : sizeOrNums.length
    this._tree = Array(this._size << 2).fill(0)
    this._lazyValue = Array(this._size << 2).fill(0)
    if (Array.isArray(sizeOrNums)) {
      this._build(1, 1, this._size, sizeOrNums)
    }
  }

  query(left: number, right: number): number {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return INF // !超出范围返回INF
    return this._query(1, left, right, 1, this._size)
  }

  add(left: number, right: number, delta: number): void {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return
    this._add(1, left, right, 1, this._size, delta)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _build(rt: number, l: number, r: number, nums: number[]): void {
    if (l === r) {
      this._tree[rt] = nums[l - 1]
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, nums)
    this._build((rt << 1) | 1, mid + 1, r, nums)
    this._pushUp(rt)
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._tree[rt]

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    let res = INF // !默认返回INF
    if (L <= mid) res = Math.min(res, this._query(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.min(res, this._query((rt << 1) | 1, L, R, mid + 1, r))

    return res
  }

  private _add(rt: number, L: number, R: number, l: number, r: number, delta: number): void {
    if (L <= l && r <= R) {
      this._lazyValue[rt] += delta
      this._tree[rt] += delta
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._add(rt << 1, L, R, l, mid, delta)
    if (mid < R) this._add((rt << 1) | 1, L, R, mid + 1, r, delta)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._tree[rt] = Math.min(this._tree[rt << 1], this._tree[(rt << 1) | 1])
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._lazyValue[rt]) {
      const delta = this._lazyValue[rt]

      this._lazyValue[rt << 1] += delta
      this._lazyValue[(rt << 1) | 1] += delta
      this._tree[rt << 1] += delta
      this._tree[(rt << 1) | 1] += delta

      this._lazyValue[rt] = 0
    }
  }
}

function solve(nums: readonly number[], k: number): number {
  const n = nums.length
  const tree = new MinSegmentTree(nums)

  let res = 0
  for (let left = 0; left + k - 1 < n; left++) {
    const right = left + k - 1
    const min = tree.query(left + 1, right + 1)
    res += min
    tree.add(left + 1, right + 1, -min)
  }

  return res
}

assert.strictEqual(solve([1, 2, 3, 2, 1], 3), 3)
assert.strictEqual(solve([5, 5, 5, 0, 5], 3), 5)

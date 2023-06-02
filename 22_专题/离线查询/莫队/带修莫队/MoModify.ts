/* eslint-disable no-inner-declarations */

import assert from 'assert'

type Query = [
  leftBlock: number,
  rightBlock: number,
  left: number,
  right: number,
  time: number,
  qid: number
]

type Modify<V> = [pos: number, val: V]

/**
 * 支持单点修改的莫队.
 * 时间复杂度: O(n^(5/3)).
 */
class MoModify<V = number> {
  private readonly _chunkSize: number
  private readonly _nums: V[]
  private readonly _queries: Query[] = []
  private readonly _modifies: Modify<V>[] = []

  constructor(nums: V[]) {
    const n = nums.length
    const nums_ = Array(n + 1)
    for (let i = 1; i <= n; i++) nums_[i] = nums[i - 1]
    this._chunkSize = Math.round(n ** (2 / 3))
    this._nums = nums_
  }

  /**
   * 添加一个查询，查询范围为`左闭右开区间` [left, right).
   * 0 <= left <= right <= n.
   */
  addQuery(left: number, right: number): void {
    left++
    this._queries.push([
      0 | (left / this._chunkSize),
      0 | ((right + 1) / this._chunkSize),
      left,
      right + 1,
      this._modifies.length,
      this._queries.length
    ])
  }

  /**
   * 添加一个修改，修改位置为 pos, 修改值为 val.
   * 0 <= pos < n.
   */
  addModify(pos: number, val: V): void {
    pos++
    this._modifies.push([pos, val])
  }

  /**
   * 返回每个查询的结果.
   * add: 将数据添加到窗口.
   * remove: 将数据从窗口移除.
   * query: 查询窗口内的数据.
   */
  run(add: (value: V) => void, remove: (value: V) => void, query: (qid: number) => void): void {
    this._queries.sort((a, b) => {
      const leftA = a[0]
      const leftB = b[0]
      if (leftA !== leftB) return leftA - leftB
      const rightA = a[1]
      const rightB = b[1]
      if (rightA !== rightB) return leftA & 1 ? rightB - rightA : rightA - rightB
      const timeA = a[4]
      const timeB = b[4]
      return rightA & 1 ? timeB - timeA : timeA - timeB
    })

    let left = 1
    let right = 1
    let now = 0

    const nums = this._nums
    const queries = this._queries
    const modifies = this._modifies
    for (let i = 0; i < queries.length; i++) {
      const [, , ql, qr, qt, qid] = queries[i]
      for (; right < qr; right++) add(nums[right])
      for (; left < ql; left++) remove(nums[left])
      for (; left > ql; left--) add(nums[left - 1])
      for (; right > qr; right--) remove(nums[right - 1])
      for (; now < qt; now++) {
        const modify = modifies[now]
        const [p, v] = modify
        if (ql <= p && p < qr) {
          remove(nums[p])
          add(v)
        }
        const old = nums[p]
        nums[p] = v
        modify[1] = old
      }
      for (; now > qt; now--) {
        const modify = modifies[now - 1]
        const [p, v] = modify
        if (ql <= p && p < qr) {
          remove(nums[p])
          add(v)
        }
        const old = nums[p]
        nums[p] = v
        modify[1] = old
      }
      query(qid)
    }
  }
}

export { MoModify }

if (require.main === module) {
  const _pool = new Map<unknown, number>()
  function id(o: unknown): number {
    if (!_pool.has(o)) {
      _pool.set(o, _pool.size)
    }
    return _pool.get(o)!
  }

  // https://www.luogu.com.cn/problem/P1903
  // 6 5
  // 1 2 3 4 5 5
  // Q 1 4
  // Q 2 6
  // R 1 2
  // Q 1 4
  // Q 2 6
  const nums = [1, 2, 3, 4, 5, 5].map(id)
  const mo = new MoModify(nums)
  const ops = [
    ['Q', 1, 4],
    ['Q', 2, 6],
    ['R', 1, 2],
    ['Q', 1, 4],
    ['Q', 2, 6]
  ] as const

  let q = 0
  ops.forEach(([op, l, r]) => {
    l--
    if (op === 'Q') {
      mo.addQuery(l, r)
      q++
    } else {
      mo.addModify(l, id(r))
    }
  })

  const res = Array(q)
  const counter = new Uint32Array(nums.length)
  let kind = 0
  mo.run(add, remove, query)
  assert.deepStrictEqual(res, [4, 4, 3, 4])

  function add(value: number): void {
    if (!counter[value]) kind++
    counter[value]++
  }
  function remove(value: number): void {
    counter[value]--
    if (!counter[value]) kind--
  }
  function query(qid: number): void {
    res[qid] = kind
  }
}

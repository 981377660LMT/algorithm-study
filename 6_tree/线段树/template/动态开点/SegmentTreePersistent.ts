/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable prefer-destructuring */

import { getSizeOf } from './memory'

type PNode<E> = [data: E, left: PNode<E> | undefined, right: PNode<E> | undefined]

class SegmentTreePersistent<E> {
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private _size!: number

  /**
   * 可持久化线段树.支持单点更新,区间查询.
   * @param e 单位元.
   * @param op 结合律的二元操作.
   */
  constructor(e: () => E, op: (a: E, b: E) => E) {
    this._e = e
    this._op = op
  }

  build(leaves: ArrayLike<E>): PNode<E> {
    this._size = leaves.length
    return this._build(0, this._size, leaves)
  }

  update(root: PNode<E>, index: number, value: E): PNode<E> {
    if (index < 0 || index >= this._size) return root
    return this._update(root, index, value, 0, this._size)
  }

  query(root: PNode<E>, start: number, end: number): E {
    if (start < 0) start = 0
    if (end > this._size) end = this._size
    if (start >= end) return this._e()
    return this._query(root, start, end, 0, this._size)
  }

  getAll(root: PNode<E>): E[] {
    const leaves: E[] = Array(this._size)
    let ptr = 0
    dfs(root)
    return leaves

    function dfs(cur: PNode<E> | undefined) {
      if (!cur) return
      if (!cur[1] && !cur[2]) {
        leaves[ptr++] = cur[0]
        return
      }
      dfs(cur[1])
      dfs(cur[2])
    }
  }

  private _build(l: number, r: number, leaves: ArrayLike<E>): PNode<E> {
    if (l + 1 >= r) return [leaves[l], undefined, undefined]
    const mid = (l + r) >> 1
    return this._merge(this._build(l, mid, leaves), this._build(mid, r, leaves))
  }

  private _merge(l: PNode<E>, r: PNode<E>): PNode<E> {
    return [this._op(l[0], r[0]), l, r]
  }

  private _update(root: PNode<E>, index: number, value: E, l: number, r: number): PNode<E> {
    if (r <= index || index + 1 <= l) return root
    if (index <= l && r <= index + 1) return [value, undefined, undefined]
    const mid = (l + r) >> 1
    return this._merge(
      this._update(root[1]!, index, value, l, mid),
      this._update(root[2]!, index, value, mid, r)
    )
  }

  private _query(root: PNode<E>, start: number, end: number, l: number, r: number): E {
    if (r <= start || end <= l) return this._e()
    if (start <= l && r <= end) return root[0]
    const mid = (l + r) >> 1
    return this._op(
      this._query(root[1]!, start, end, l, mid),
      this._query(root[2]!, start, end, mid, r)
    )
  }
}

export { SegmentTreePersistent }

if (require.main === module) {
  const seg = new SegmentTreePersistent(
    () => 0,
    (a, b) => a + b
  )

  let root = seg.build([1, 2, 3, 4, 5])
  console.log(seg.getAll(root))
  root = seg.update(root, 2, 10)
  console.log(seg.getAll(root))
  console.log(seg.query(root, 0, 3))

  console.time('build')
  const size = getSizeOf(() => {
    const n = 2e5
    let newRoot = seg.build(Array(n).fill(1))
    for (let i = 0; i < n; ++i) {
      newRoot = seg.update(newRoot, i, 1)
      seg.query(newRoot, 0, i)
    }
  })
  console.timeEnd('build') // build: 304.167ms
  console.log(size) // 54.44MB
}

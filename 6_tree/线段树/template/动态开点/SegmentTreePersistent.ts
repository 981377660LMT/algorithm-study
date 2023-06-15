/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable prefer-destructuring */

type SegNode<E> = {
  data: E
  left: SegNode<E> | undefined
  right: SegNode<E> | undefined
}

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

  build(leaves: ArrayLike<E>): SegNode<E> {
    this._size = leaves.length
    return this._build(0, this._size, leaves)
  }

  set(root: SegNode<E>, index: number, value: E): SegNode<E> {
    if (index < 0 || index >= this._size) return root
    return this._set(root, index, value, 0, this._size)
  }

  update(root: SegNode<E>, index: number, value: E): SegNode<E> {
    if (index < 0 || index >= this._size) return root
    return this._update(root, index, value, 0, this._size)
  }

  query(root: SegNode<E>, start: number, end: number): E {
    if (start < 0) start = 0
    if (end > this._size) end = this._size
    if (start >= end) return this._e()
    return this._query(root, start, end, 0, this._size)
  }

  getAll(root: SegNode<E>): E[] {
    const leaves: E[] = Array(this._size)
    let ptr = 0
    dfs(root)
    return leaves

    function dfs(cur: SegNode<E> | undefined) {
      if (!cur) return
      if (!cur.left && !cur.right) {
        leaves[ptr++] = cur.data
        return
      }
      dfs(cur.left)
      dfs(cur.right)
    }
  }

  private _build(l: number, r: number, leaves: ArrayLike<E>): SegNode<E> {
    if (l + 1 >= r) return { data: leaves[l], left: undefined, right: undefined }
    const mid = (l + r) >> 1
    return this._merge(this._build(l, mid, leaves), this._build(mid, r, leaves))
  }

  private _merge(l: SegNode<E>, r: SegNode<E>): SegNode<E> {
    return { data: this._op(l.data, r.data), left: l, right: r }
  }

  private _set(root: SegNode<E>, index: number, value: E, l: number, r: number): SegNode<E> {
    if (r <= index || index + 1 <= l) return root
    if (index <= l && r <= index + 1) return { data: value, left: undefined, right: undefined }
    const mid = (l + r) >> 1
    return this._merge(
      this._set(root.left!, index, value, l, mid),
      this._set(root.right!, index, value, mid, r)
    )
  }

  private _update(root: SegNode<E>, index: number, value: E, l: number, r: number): SegNode<E> {
    if (r <= index || index + 1 <= l) return root
    if (index <= l && r <= index + 1) {
      return { data: this._op(root.data, value), left: undefined, right: undefined }
    }
    const mid = (l + r) >> 1
    return this._merge(
      this._update(root.left!, index, value, l, mid),
      this._update(root.right!, index, value, mid, r)
    )
  }

  private _query(root: SegNode<E>, start: number, end: number, l: number, r: number): E {
    if (r <= start || end <= l) return this._e()
    if (start <= l && r <= end) return root.data
    const mid = (l + r) >> 1
    return this._op(
      this._query(root.left!, start, end, l, mid),
      this._query(root.right!, start, end, mid, r)
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

  root = seg.update(root, 2, 123)
  console.log(seg.getAll(root))

  root = seg.set(root, 2, 123)
  console.log(seg.getAll(root))

  console.log(seg.query(root, 0, 3))

  console.time('build')
  const n = 2e5
  let newRoot = seg.build(Array(n).fill(1))
  for (let i = 0; i < n; ++i) {
    newRoot = seg.set(newRoot, i, 1)
    newRoot = seg.update(newRoot, i, 1)
    seg.query(newRoot, 0, i)
  }
  console.timeEnd('build') // build: 391.185ms
}

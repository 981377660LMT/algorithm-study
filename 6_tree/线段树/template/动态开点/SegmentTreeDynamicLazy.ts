/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-param-reassign */
/* eslint-disable prefer-destructuring */

type SegNode<E, Id> = {
  left: SegNode<E, Id> | undefined
  right: SegNode<E, Id> | undefined
  data: E
  id: Id
}

class SegmentTreeDynamicLazy<E = number, Id = number> {
  private static _isPrimitive(
    o: unknown
  ): o is number | string | boolean | symbol | bigint | null | undefined {
    return o === null || (typeof o !== 'object' && typeof o !== 'function')
  }

  private readonly _lower: number
  private readonly _upper: number
  private readonly _e: () => E
  private readonly _eRange: (start: number, end: number) => E
  private readonly _id: () => Id
  private readonly _op: (a: E, b: E) => E
  private readonly _mapping: (id: Id, data: E, size: number) => E
  private readonly _composition: (id1: Id, id2: Id) => Id
  private readonly _equalsToId: (o: Id) => boolean

  /**
   * 区间修改区间查询的动态开点懒标记线段树.线段树维护的值域为`[start, end)`.
   * @param start 值域下界.start>=0.
   * @param end 值域上界.end<=2**31-1.
   * @param operations 线段树的操作.
   * @alias NodeManager
   */
  constructor(
    start: number,
    end: number,
    operations: {
      /**
       * 线段树维护的值的幺元.
       */
      e: () => E

      /**
       * 结点的初始值.用于维护结点的范围.
       */
      eRange?: (start: number, end: number) => E

      /**
       * 更新操作/懒标记的幺元.
       */
      id: () => Id

      /**
       * 合并左右区间的值.
       */
      op: (e1: E, e2: E) => E

      /**
       * 父结点的懒标记更新子结点的值.
       */
      mapping: (lazy: Id, data: E, size: number) => E

      /**
       * 父结点的懒标记更新子结点的懒标记(合并).
       */
      composition: (f: Id, g: Id) => Id

      /**
       * 判断两个懒标记是否相等.默认为`===`.
       */
      equalsId?: (id1: Id, id2: Id) => boolean
    } & ThisType<void>
  ) {
    const { e, eRange, id, op, mapping, composition, equalsId } = operations
    if (!equalsId && !SegmentTreeDynamicLazy._isPrimitive(id())) {
      throw new Error('equalsId must be provided when id() returns an non-primitive value')
    }

    this._lower = start
    this._upper = end
    this._e = e
    this._eRange = eRange || e
    this._id = id
    this._op = op
    this._mapping = mapping
    this._composition = composition
    const identity = id()
    this._equalsToId = equalsId ? (o: Id) => equalsId(o, identity) : (o: Id) => o === identity
  }

  newRoot(): SegNode<E, Id> {
    return {
      left: undefined,
      right: undefined,
      data: this._eRange(this._lower, this._upper),
      id: this._id()
    }
  }

  build(leaves: ArrayLike<E>): SegNode<E, Id> {
    return this._build(0, leaves.length, leaves)!
  }

  get(root: SegNode<E, Id>, index: number): E {
    return this.query(root, index, index + 1)
  }

  set(root: SegNode<E, Id>, index: number, value: E): SegNode<E, Id> {
    if (index < this._lower || index >= this._upper) return root
    return this._set(root, this._lower, this._upper, index, value)
  }

  update(root: SegNode<E, Id>, index: number, value: E): SegNode<E, Id> {
    if (index < this._lower || index >= this._upper) return root
    return this._update(root, this._lower, this._upper, index, value)
  }

  /**
   * 区间`[start,end)`的值与`lazy`进行作用.
   * {@link _lower} <= start <= end <= {@link _upper}.
   */
  updateRange(root: SegNode<E, Id>, start: number, end: number, lazy: Id): SegNode<E, Id> {
    if (start < this._lower) start = this._lower
    if (end > this._upper) end = this._upper
    if (start >= end) return root
    return this._updateRange(root, this._lower, this._upper, start, end, lazy)
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * {@link _lower} <= start <= end <= {@link _upper}.
   */
  query(root: SegNode<E, Id>, start: number, end: number): E {
    if (start < this._lower) start = this._lower
    if (end > this._upper) end = this._upper
    if (start >= end) return this._e()

    let res = this._e()
    const _query = (
      node: SegNode<E, Id> | undefined,
      l: number,
      r: number,
      ql: number,
      qr: number,
      lazy: Id
    ) => {
      ql = l > ql ? l : ql
      qr = r < qr ? r : qr
      if (ql >= qr) return
      if (!node) {
        res = this._op(res, this._mapping(lazy, this._eRange(ql, qr), qr - ql))
        return
      }
      if (l === ql && r === qr) {
        res = this._op(res, this._mapping(lazy, node.data, r - l))
        return
      }
      const mid = (l + r) >>> 1
      lazy = this._composition(lazy, node.id)
      _query(node.left, l, mid, ql, qr, lazy)
      _query(node.right, mid, r, ql, qr, lazy)
    }

    _query(root, this._lower, this._upper, start, end, this._id())
    return res
  }

  queryAll(root: SegNode<E, Id>): E {
    return root.data
  }

  /**
   * 二分查询最大的`end`使得切片`[start:end)`内的聚合值满足`check`.
   * {@link _lower} <= start <= {@link _upper}.
   * @alias findFirst
   */
  maxRight(root: SegNode<E, Id>, start: number, check: (e: E) => boolean): number {
    if (start < this._lower) start = this._lower
    if (start >= this._upper) return this._upper
    let x = this._e()
    const _maxRight = (
      node: SegNode<E, Id> | undefined,
      l: number,
      r: number,
      ql: number
    ): number => {
      if (r <= ql) return r
      if (!node) node = this._newNode(l, r)
      ql = l > ql ? l : ql
      if (l === ql) {
        const tmp = this._op(x, node.data)
        if (check(tmp)) {
          x = tmp
          return r
        }
      }
      if (r === l + 1) return l
      this._pushDown(node, l, r)
      const m = (l + r) >>> 1
      const k = _maxRight(node.left, l, m, ql)
      if (m > k) return k
      return _maxRight(node.right, m, r, ql)
    }
    return _maxRight(root, this._lower, this._upper, start)
  }

  /**
   * 二分查询最小的`start`使得切片`[start:end)`内的聚合值满足`check`.
   * {@link _lower} <= end <= {@link _upper}.
   * @alias findLast
   */
  minLeft(root: SegNode<E, Id>, end: number, check: (e: E) => boolean): number {
    if (end > this._upper) end = this._upper
    if (end <= this._lower) return this._lower
    let x = this._e()
    const _minLeft = (
      node: SegNode<E, Id> | undefined,
      l: number,
      r: number,
      qr: number
    ): number => {
      if (qr <= l) return l
      if (!node) node = this._newNode(l, r)
      qr = r < qr ? r : qr
      if (r === qr) {
        const tmp = this._op(node.data, x)
        if (check(tmp)) {
          x = tmp
          return l
        }
      }
      if (r === l + 1) return r
      this._pushDown(node, l, r)
      const m = (l + r) >>> 1
      const k = _minLeft(node.right, m, r, qr)
      if (m < k) return k
      return _minLeft(node.left, l, m, qr)
    }
    return _minLeft(root, this._lower, this._upper, end)
  }

  /**
   * `破坏性`地合并node1和node2.
   * @warning Not Verified.
   */
  mergeDestructively(node1: SegNode<E, Id>, node2: SegNode<E, Id>): SegNode<E, Id> {
    const newRoot = this._merge(node1, node2)
    if (!newRoot) throw new Error('merge failed')
    return newRoot
  }

  getAll(root: SegNode<E, Id>): E[] {
    if (this._upper - this._lower > 1e7) throw new Error('too large')
    const res: E[] = []
    const _getAll = (node: SegNode<E, Id> | undefined, l: number, r: number, lazy: Id) => {
      if (!node) node = this._newNode(l, r)
      if (r - l === 1) {
        res.push(this._mapping(lazy, node.data, 1))
        return
      }
      const m = (l + r) >>> 1
      lazy = this._composition(lazy, node.id)
      _getAll(node.left, l, m, lazy)
      _getAll(node.right, m, r, lazy)
    }
    _getAll(root, this._lower, this._upper, this._id())
    return res
  }

  private _set(root: SegNode<E, Id>, l: number, r: number, i: number, x: E): SegNode<E, Id> {
    if (l === r - 1) {
      root.data = x
      root.id = this._id()
      return root
    }
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    if (!root.left) root.left = this._newNode(l, mid)
    if (!root.right) root.right = this._newNode(mid, r)

    if (i < mid) {
      root.left = this._set(root.left!, l, mid, i, x)
    } else {
      root.right = this._set(root.right!, mid, r, i, x)
    }
    root.data = this._op(root.left!.data, root.right!.data)
    return root
  }

  private _update(root: SegNode<E, Id>, l: number, r: number, i: number, x: E): SegNode<E, Id> {
    if (l === r - 1) {
      root.data = this._op(root.data, x)
      root.id = this._id()
      return root
    }
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    if (!root.left) root.left = this._newNode(l, mid)
    if (!root.right) root.right = this._newNode(mid, r)

    if (i < mid) {
      root.left = this._update(root.left!, l, mid, i, x)
    } else {
      root.right = this._update(root.right!, mid, r, i, x)
    }
    root.data = this._op(root.left!.data, root.right!.data)
    return root
  }

  private _updateRange(
    root: SegNode<E, Id> | undefined,
    l: number,
    r: number,
    ql: number,
    qr: number,
    lazy: Id
  ): SegNode<E, Id> {
    if (!root) root = this._newNode(l, r)
    ql = l > ql ? l : ql
    qr = r < qr ? r : qr
    if (ql >= qr) return root
    if (l === ql && r === qr) {
      root.data = this._mapping(lazy, root.data, r - l)
      root.id = this._composition(lazy, root.id)
      return root
    }
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1

    root.left = this._updateRange(root.left, l, mid, ql, qr, lazy)
    root.right = this._updateRange(root.right, mid, r, ql, qr, lazy)
    root.data = this._op(root.left!.data, root.right!.data)
    return root
  }

  private _pushDown(node: SegNode<E, Id>, l: number, r: number): void {
    const lazy = node.id
    if (this._equalsToId(lazy)) return
    const mid = (l + r) >>> 1

    if (!node.left) node.left = this._newNode(l, mid)
    const leftChild = node.left!
    leftChild.data = this._mapping(lazy, leftChild.data, mid - l)
    leftChild.id = this._composition(lazy, leftChild.id)

    if (!node.right) node.right = this._newNode(mid, r)
    const rightChild = node.right!
    rightChild.data = this._mapping(lazy, rightChild.data, r - mid)
    rightChild.id = this._composition(lazy, rightChild.id)

    node.id = this._id()
  }

  private _newNode(l: number, r: number): SegNode<E, Id> {
    return { left: undefined, right: undefined, data: this._eRange(l, r), id: this._id() }
  }

  private _build(left: number, right: number, nums: ArrayLike<E>): SegNode<E, Id> | undefined {
    if (left === right) return undefined
    if (right === left + 1) {
      return { left: undefined, right: undefined, data: nums[left], id: this._id() }
    }
    const mid = (left + right) >>> 1
    const lRoot = this._build(left, mid, nums)
    const rRoot = this._build(mid, right, nums)
    return { left: lRoot, right: rRoot, data: this._op(lRoot!.data, rRoot!.data), id: this._id() }
  }

  private _merge(
    node1: SegNode<E, Id> | undefined,
    node2: SegNode<E, Id> | undefined
  ): SegNode<E, Id> | undefined {
    if (!node1 || !node2) return node1 || node2
    node1.left = this._merge(node1.left, node2.left)
    node1.right = this._merge(node1.right, node2.right)
    // pushUp
    const left = node1.left
    const right = node1.right
    node1.data = this._op(left ? left.data : this._e(), right ? right.data : this._e())
    return node1
  }
}

export { SegmentTreeDynamicLazy }

if (require.main === module) {
  const n = 1e9
  // RangeAddRangeSum
  const seg = new SegmentTreeDynamicLazy<number, number>(0, n, {
    e() {
      return 0
    },
    id() {
      return 0
    },
    op(x, y) {
      return x + y
    },
    mapping(f, x) {
      return f + x
    },
    composition(f, g) {
      return f + g
    }
    // equalsId(id1, id2) {}
  })
  let root = seg.newRoot()
  seg.set(root, 0, 1)
  console.log(seg.query(root, 0, 1))
  // console.log(seg.getAll(root))

  const M = 1e5
  const pos = Array.from({ length: M }, () => Math.floor(Math.random() * n))
  console.time('start')

  for (let i = 0; i < M; ++i) {
    root = seg.set(root, pos[i], 1)
    seg.query(root, 0, i)
    root = seg.updateRange(root, 0, i, 1)
    seg.get(root, i)
  }
  console.timeEnd('start') // start: 455.065ms

  // https://leetcode.cn/problems/range-module/
  // RangeAssignRangeSum
  class RangeModule {
    private readonly _seg = new SegmentTreeDynamicLazy<number, number>(0, 1e9, {
      e() {
        return 0
      },
      id() {
        return -1
      },
      op(e1, e2) {
        return e1 + e2
      },
      mapping(lazy, data, size) {
        return lazy === -1 ? data : lazy * size
      },
      composition(f, g) {
        return f === -1 ? g : f
      }
    })
    private readonly _root = this._seg.newRoot()

    addRange(left: number, right: number): void {
      this._seg.updateRange(this._root, left, right, 1)
    }

    queryRange(left: number, right: number): boolean {
      return this._seg.query(this._root, left, right) === right - left
    }

    removeRange(left: number, right: number): void {
      this._seg.updateRange(this._root, left, right, 0)
    }
  }
}

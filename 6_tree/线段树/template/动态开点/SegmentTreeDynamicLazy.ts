/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-param-reassign */

type SegNode<E, Id> = [
  left: SegNode<E, Id> | undefined,
  right: SegNode<E, Id> | undefined,
  data: E,
  lazy: Id
]

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
  private readonly _persistent: boolean

  /**
   * 区间修改区间查询的动态开点懒标记线段树.线段树维护的值域为`[start, end)`.
   * @param start 值域下界.start>=0.
   * @param end 值域上界.end<=2**31-1.
   * @param operations 线段树的操作.
   * @param persistent 是否持久化.持久化后,每次修改都会新建一个结点,否则会复用原来的结点.
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
    } & ThisType<void>,
    persistent = false
  ) {
    const { e, eRange, id, op, mapping, composition, equalsId } = operations
    if (!equalsId && !SegmentTreeDynamicLazy._isPrimitive(id())) {
      throw new Error('equalsId must be provided when id() returns an non-primitive value')
    }
    if (
      persistent &&
      !(SegmentTreeDynamicLazy._isPrimitive(e()) && SegmentTreeDynamicLazy._isPrimitive(id()))
    ) {
      throw new Error('persistent is only supported when e() and id() return primitive values')
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
    this._persistent = persistent
  }

  newRoot(): SegNode<E, Id> {
    return [undefined, undefined, this._eRange(this._lower, this._upper), this._id()]
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
        res = this._op(res, this._mapping(lazy, node[2], r - l))
        return
      }
      const mid = (l + r) >>> 1
      lazy = this._composition(lazy, node[3])
      _query(node[0], l, mid, ql, qr, lazy)
      _query(node[1], mid, r, ql, qr, lazy)
    }

    _query(root, this._lower, this._upper, start, end, this._id())
    return res
  }

  queryAll(root: SegNode<E, Id>): E {
    return root[2]
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
        const tmp = this._op(x, node[2])
        if (check(tmp)) {
          x = tmp
          return r
        }
      }
      if (r === l + 1) return l
      this._pushDown(node, l, r)
      const m = (l + r) >>> 1
      const k = _maxRight(node[0], l, m, ql)
      if (m > k) return k
      return _maxRight(node[1], m, r, ql)
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
        const tmp = this._op(node[2], x)
        if (check(tmp)) {
          x = tmp
          return l
        }
      }
      if (r === l + 1) return r
      this._pushDown(node, l, r)
      const m = (l + r) >>> 1
      const k = _minLeft(node[1], m, r, qr)
      if (m < k) return k
      return _minLeft(node[0], l, m, qr)
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
    const res: E[] = []
    const _getAll = (node: SegNode<E, Id> | undefined, l: number, r: number, lazy: Id) => {
      if (!node) node = this._newNode(l, r)
      if (r - l === 1) {
        res.push(this._mapping(lazy, node[2], 1))
        return
      }
      const m = (l + r) >>> 1
      lazy = this._composition(lazy, node[3])
      _getAll(node[0], l, m, lazy)
      _getAll(node[1], m, r, lazy)
    }
    _getAll(root, this._lower, this._upper, this._id())
    return res
  }

  private _copyNode(node: SegNode<E, Id>): SegNode<E, Id> {
    if (!node || !this._persistent) return node
    return [node[0], node[1], node[2], node[3]] // TODO: 如果是引用类型, 持久化时需要深拷贝
  }

  private _set(root: SegNode<E, Id>, l: number, r: number, i: number, x: E): SegNode<E, Id> {
    if (l === r - 1) {
      root = this._copyNode(root)
      root[2] = x
      root[3] = this._id()
      return root
    }
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    if (!root[0]) root[0] = this._newNode(l, mid)
    if (!root[1]) root[1] = this._newNode(mid, r)
    root = this._copyNode(root)
    if (i < mid) {
      root[0] = this._set(root[0]!, l, mid, i, x)
    } else {
      root[1] = this._set(root[1]!, mid, r, i, x)
    }
    root[2] = this._op(root[0]![2], root[1]![2])
    return root
  }

  private _update(root: SegNode<E, Id>, l: number, r: number, i: number, x: E): SegNode<E, Id> {
    if (l === r - 1) {
      root = this._copyNode(root)
      root[2] = this._op(root[2], x)
      root[3] = this._id()
      return root
    }
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    if (!root[0]) root[0] = this._newNode(l, mid)
    if (!root[1]) root[1] = this._newNode(mid, r)
    root = this._copyNode(root)
    if (i < mid) {
      root[0] = this._update(root[0]!, l, mid, i, x)
    } else {
      root[1] = this._update(root[1]!, mid, r, i, x)
    }
    root[2] = this._op(root[0]![2], root[1]![2])
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
      root = this._copyNode(root)
      root[2] = this._mapping(lazy, root[2], r - l)
      root[3] = this._composition(lazy, root[3])
      return root
    }
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    root = this._copyNode(root)
    root[0] = this._updateRange(root[0], l, mid, ql, qr, lazy)
    root[1] = this._updateRange(root[1], mid, r, ql, qr, lazy)
    root[2] = this._op(root[0]![2], root[1]![2])
    return root
  }

  private _pushDown(node: SegNode<E, Id>, l: number, r: number): void {
    const lazy = node[3]
    if (this._equalsToId(lazy)) return
    const mid = (l + r) >>> 1
    if (!node[0]) {
      node[0] = this._newNode(l, mid)
    } else {
      node[0] = this._copyNode(node[0])
    }
    const leftChild = node[0]!
    leftChild[2] = this._mapping(lazy, leftChild[2], mid - l)
    leftChild[3] = this._composition(lazy, leftChild[3])
    if (!node[1]) {
      node[1] = this._newNode(mid, r)
    } else {
      node[1] = this._copyNode(node[1])
    }
    const rightChild = node[1]!
    rightChild[2] = this._mapping(lazy, rightChild[2], r - mid)
    rightChild[3] = this._composition(lazy, rightChild[3])
    node[3] = this._id()
  }

  private _newNode(l: number, r: number): SegNode<E, Id> {
    return [undefined, undefined, this._eRange(l, r), this._id()]
  }

  private _build(left: number, right: number, nums: ArrayLike<E>): SegNode<E, Id> | undefined {
    if (left === right) return undefined
    if (right === left + 1) return [undefined, undefined, nums[left], this._id()]
    const mid = (left + right) >>> 1
    const lRoot = this._build(left, mid, nums)
    const rRoot = this._build(mid, right, nums)
    return [lRoot, rRoot, this._op(lRoot![2], rRoot![2]), this._id()]
  }

  private _merge(
    node1: SegNode<E, Id> | undefined,
    node2: SegNode<E, Id> | undefined
  ): SegNode<E, Id> | undefined {
    if (!node1 || !node2) return node1 || node2
    node1[0] = this._merge(node1[0], node2[0])
    node1[1] = this._merge(node1[1], node2[1])
    // pushUp
    const left = node1[0]
    const right = node1[1]
    node1[2] = this._op(left ? left[2] : this._e(), right ? right[2] : this._e())
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

  const M = 1e5
  const pos = Array.from({ length: M }, () => Math.floor(Math.random() * n))
  console.time('start')

  for (let i = 0; i < M; ++i) {
    root = seg.set(root, pos[i], 1)
    seg.query(root, 0, i)
    root = seg.updateRange(root, 0, i, 1)
    seg.get(root, i)
  }
  console.timeEnd('start') // start: 681.428ms

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
      op(x, y) {
        return x + y
      },
      mapping(f, x, size) {
        return f === -1 ? x : f * size
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

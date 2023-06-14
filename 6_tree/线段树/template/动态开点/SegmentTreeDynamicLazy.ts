/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-param-reassign */
/* eslint-disable prefer-destructuring */

// API:
//  new SegmentTreeDynamicLazy(start, end, operations, persistent)
//  newRoot()
//  build(leaves)
//  get(root, index)
//  set(root, index, value)
//  update(root, index, value)
//  updateRange(root, start, end, lazy)
//  query(root, start, end)
//  queryAll(root)
//  maxRight(root, start, predicate)
//  minLeft(root, end, predicate)
//  getAll(root)
//  mergeDestructively(root1, root2)

/**
 * !不要用数组`[]`代替对象`{}`.数组会导致性能下降.
 */
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
  private readonly _persistent: boolean
  private _root: SegNode<E, Id>

  /**
   * 区间修改区间查询的动态开点懒标记线段树.线段树维护的值域为`[start, end)`.
   * @param start 值域下界.start>=0.
   * @param end 值域上界.
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
       * 判断两个懒标记是否相等.比较方式默认为`===`.
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
    this._upper = end + 5
    this._e = e
    this._eRange = eRange || e
    this._id = id
    this._op = op
    this._mapping = mapping
    this._composition = composition
    const identity = id()
    this._equalsToId = equalsId ? (o: Id) => equalsId(o, identity) : (o: Id) => o === identity
    this._persistent = persistent

    this._root = this.newRoot()
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
    const newRoot = this._build(0, leaves.length, leaves)!
    this._root = newRoot
    return newRoot
  }

  get(index: number, root: SegNode<E, Id> = this._root): E {
    return this.query(index, index + 1, root)
  }

  set(index: number, value: E, root: SegNode<E, Id> = this._root): SegNode<E, Id> {
    if (index < this._lower || index >= this._upper) return root
    const newRoot = this._set(root, this._lower, this._upper, index, value)
    this._root = newRoot
    return newRoot
  }

  update(index: number, value: E, root: SegNode<E, Id> = this._root): SegNode<E, Id> {
    if (index < this._lower || index >= this._upper) return root
    const newRoot = this._update(root, this._lower, this._upper, index, value)
    this._root = newRoot
    return newRoot
  }

  /**
   * 区间`[start,end)`的值与`lazy`进行作用.
   * {@link _lower} <= start <= end <= {@link _upper}.
   */
  updateRange(
    start: number,
    end: number,
    lazy: Id,
    root: SegNode<E, Id> = this._root
  ): SegNode<E, Id> {
    if (start < this._lower) start = this._lower
    if (end > this._upper) end = this._upper
    if (start >= end) return root
    const newRoot = this._updateRange(root, this._lower, this._upper, start, end, lazy)
    this._root = newRoot
    return newRoot
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * {@link _lower} <= start <= end <= {@link _upper}.
   */
  query(start: number, end: number, root: SegNode<E, Id> = this._root): E {
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
      const mid = Math.floor(l + (r - l) / 2)
      lazy = this._composition(lazy, node.id)
      _query(node.left, l, mid, ql, qr, lazy)
      _query(node.right, mid, r, ql, qr, lazy)
    }

    _query(root, this._lower, this._upper, start, end, this._id())
    return res
  }

  queryAll(root: SegNode<E, Id> = this._root): E {
    return root.data
  }

  /**
   * 二分查询最大的`end`使得切片`[start:end)`内的聚合值满足`check`.
   * {@link _lower} <= start <= {@link _upper}.
   * @alias findFirst
   */
  maxRight(start: number, check: (e: E) => boolean, root: SegNode<E, Id> = this._root): number {
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
      const m = Math.floor(l + (r - l) / 2)
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
  minLeft(end: number, check: (e: E) => boolean, root: SegNode<E, Id> = this._root): number {
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
      const m = Math.floor(l + (r - l) / 2)
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

  getAll(root: SegNode<E, Id> = this._root): E[] {
    if (this._upper - this._lower > 1e7) throw new Error('too large')
    const res: E[] = []
    const _getAll = (node: SegNode<E, Id> | undefined, l: number, r: number, lazy: Id) => {
      if (!node) node = this._newNode(l, r)
      if (r - l === 1) {
        res.push(this._mapping(lazy, node.data, 1))
        return
      }
      const m = Math.floor(l + (r - l) / 2)
      lazy = this._composition(lazy, node.id)
      _getAll(node.left, l, m, lazy)
      _getAll(node.right, m, r, lazy)
    }
    _getAll(root, this._lower, this._upper, this._id())
    return res
  }

  private _copyNode(node: SegNode<E, Id>): SegNode<E, Id> {
    if (!node || !this._persistent) return node
    // TODO: 如果是引用类型, 持久化时需要深拷贝
    // !不要使用`...`,很慢
    return { left: node.left, right: node.right, data: node.data, id: node.id }
  }

  private _set(root: SegNode<E, Id>, l: number, r: number, i: number, x: E): SegNode<E, Id> {
    if (l === r - 1) {
      root = this._copyNode(root)
      root.data = x
      root.id = this._id()
      return root
    }
    this._pushDown(root, l, r)
    const m = Math.floor(l + (r - l) / 2)
    if (!root.left) root.left = this._newNode(l, m)
    if (!root.right) root.right = this._newNode(m, r)
    root = this._copyNode(root)
    if (i < m) {
      root.left = this._set(root.left!, l, m, i, x)
    } else {
      root.right = this._set(root.right!, m, r, i, x)
    }
    root.data = this._op(root.left!.data, root.right!.data)
    return root
  }

  private _update(root: SegNode<E, Id>, l: number, r: number, i: number, x: E): SegNode<E, Id> {
    if (l === r - 1) {
      root = this._copyNode(root)
      root.data = this._op(root.data, x)
      root.id = this._id()
      return root
    }
    this._pushDown(root, l, r)
    const m = Math.floor(l + (r - l) / 2)
    if (!root.left) root.left = this._newNode(l, m)
    if (!root.right) root.right = this._newNode(m, r)
    root = this._copyNode(root)
    if (i < m) {
      root.left = this._update(root.left!, l, m, i, x)
    } else {
      root.right = this._update(root.right!, m, r, i, x)
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
      root = this._copyNode(root)
      root.data = this._mapping(lazy, root.data, r - l)
      root.id = this._composition(lazy, root.id)
      return root
    }
    this._pushDown(root, l, r)
    const m = Math.floor(l + (r - l) / 2)
    root = this._copyNode(root)
    root.left = this._updateRange(root.left, l, m, ql, qr, lazy)
    root.right = this._updateRange(root.right, m, r, ql, qr, lazy)
    root.data = this._op(root.left!.data, root.right!.data)
    return root
  }

  private _pushDown(node: SegNode<E, Id>, l: number, r: number): void {
    const lazy = node.id
    if (this._equalsToId(lazy)) return
    const m = Math.floor(l + (r - l) / 2)

    if (!node.left) {
      node.left = this._newNode(l, m)
    } else {
      node.left = this._copyNode(node.left)
    }
    const leftChild = node.left!
    leftChild.data = this._mapping(lazy, leftChild.data, m - l)
    leftChild.id = this._composition(lazy, leftChild.id)

    if (!node.right) {
      node.right = this._newNode(m, r)
    } else {
      node.right = this._copyNode(node.right)
    }
    const rightChild = node.right!
    rightChild.data = this._mapping(lazy, rightChild.data, r - m)
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
    const m = (left + right) >>> 1
    const lRoot = this._build(left, m, nums)
    const rRoot = this._build(m, right, nums)
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
  root = seg.set(0, 1, root)
  console.log(seg.query(0, 1, root))
  // console.log(seg.getAll(root))

  const M = 1e5
  const pos = Array.from({ length: M }, () => Math.floor(Math.random() * n))
  console.time('start')
  for (let i = 0; i < M; ++i) {
    root = seg.set(pos[i], 1, root)
    seg.query(0, i, root)
    root = seg.updateRange(0, i, 1, root)
    seg.get(i, root)
  }
  console.timeEnd('start') // start: 474.62ms

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

    addRange(left: number, right: number): void {
      this._seg.updateRange(left, right, 1)
    }

    queryRange(left: number, right: number): boolean {
      return this._seg.query(left, right) === right - left
    }

    removeRange(left: number, right: number): void {
      this._seg.updateRange(left, right, 0)
    }
  }

  class CountIntervals {
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

    add(left: number, right: number): void {
      this._seg.updateRange(left, right + 1, 1)
    }

    count(): number {
      return this._seg.queryAll()
    }
  }

  /**
   * Your CountIntervals object will be instantiated and called as such:
   * var obj = new CountIntervals()
   * obj.add(left,right)
   * var param_2 = obj.count()
   */
}

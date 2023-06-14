type DNode<E, Id> = [
  left: DNode<E, Id> | undefined,
  right: DNode<E, Id> | undefined,
  data: E,
  lazy: Id
]

class SegmentTreeDynamicLazy<E = number, Id = number> {
  private static _isPrimitive(
    o: unknown
  ): o is number | string | boolean | symbol | bigint | null | undefined {
    return o === null || (typeof o !== 'object' && typeof o !== 'function')
  }

  private readonly _n: number
  private readonly _e: () => E
  private readonly _eRange: (start: number, end: number) => E
  private readonly _id: () => Id
  private readonly _op: (a: E, b: E) => E
  private readonly _mapping: (id: Id, data: E, size: number) => E
  private readonly _composition: (id1: Id, id2: Id) => Id
  private readonly _equalsToId: (o: Id) => boolean
  private _root: DNode<E, Id>

  /**
   * 区间修改区间查询的动态开点懒标记线段树.
   * @param n 值域上界.线段树维护的值域为`[0, n)`.
   * @param operations 线段树的操作.
   */
  constructor(
    n: number,
    operations: {
      /**
       * 线段树维护的值的幺元.
       */
      e: () => E

      /**
       * 结点的初始值.
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

    let size = 2
    while (size <= n) size <<= 1
    this._n = size
    this._e = e
    this._eRange = eRange || e
    this._id = id
    this._op = op
    this._mapping = mapping
    this._composition = composition
    const identity = id()
    this._equalsToId = equalsId ? (o: Id) => equalsId(o, identity) : (o: Id) => o === identity

    this._root = this._newNode(0, this._n)
  }

  set(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    this._set(this._root, 0, this._n, index, value)
  }

  get(index: number): E {
    if (index < 0 || index >= this._n) return this._e()
    let left = 0
    let right = this._n
    let root = this._root
    while (left + 1 < right) {
      this._pushDown(root, left, right)
      const mid = (left + right) >>> 1
      if (index < mid) {
        const leftChild = root[0]
        if (!leftChild) return this._e()
        root = leftChild
        right = mid
      } else {
        const rightChild = root[1]
        if (!rightChild) return this._e()
        root = rightChild
        left = mid
      }
    }
    return root[2]
  }

  /**
   * 区间`[start,end)`的值与`lazy`进行作用.
   * 0 <= start <= end <= n.
   */
  update(start: number, end: number, lazy: Id): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    this._update(this._root, 0, this._n, start, end, lazy)
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): E {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return this._eRange(start, start)
    return this._query(this._root, 0, this._n, start, end)
  }

  queryAll(): E {
    return this._root[2]
  }

  mergeDestructively(other: SegmentTreeDynamicLazy<E, Id>): void {
    const newRoot = this._merge(this._root, other._root)
    if (!newRoot) throw new Error('merge failed')
    this._root = newRoot
  }

  private _newNode(l: number, r: number): DNode<E, Id> {
    return [undefined, undefined, this._eRange(l, r), this._id()]
  }

  private _pushUp(node: DNode<E, Id>): void {
    const left = node[0]
    const right = node[1]
    node[2] = this._op(left ? left[2] : this._e(), right ? right[2] : this._e())
  }

  private _pushDown(node: DNode<E, Id>, l: number, r: number): void {
    const lazy = node[3]
    if (this._equalsToId(lazy)) return
    const mid = (l + r) >>> 1
    let leftChild = node[0]
    if (!leftChild) {
      leftChild = this._newNode(l, mid)
      node[0] = leftChild
    }
    leftChild[2] = this._mapping(lazy, leftChild[2], mid - l)
    leftChild[3] = this._composition(lazy, leftChild[3])
    let rightChild = node[1]
    if (!rightChild) {
      rightChild = this._newNode(mid, r)
      node[1] = rightChild
    }
    rightChild[2] = this._mapping(lazy, rightChild[2], r - mid)
    rightChild[3] = this._composition(lazy, rightChild[3])
    node[3] = this._id()
  }

  private _propagate(node: DNode<E, Id>, lazy: Id, size: number): void {
    node[2] = this._mapping(lazy, node[2], size)
    node[3] = this._composition(lazy, node[3])
  }

  private _set(node: DNode<E, Id>, l: number, r: number, i: number, x: E): void {
    if (l + 1 === r) {
      node[2] = x
      return
    }
    const mid = (l + r) >>> 1
    this._pushDown(node, l, r)
    if (i < mid) {
      if (!node[0]) node[0] = this._newNode(l, mid)
      this._set(node[0], l, mid, i, x)
    } else {
      if (!node[1]) node[1] = this._newNode(mid, r)
      this._set(node[1], mid, r, i, x)
    }
    this._pushUp(node)
  }

  private _update(node: DNode<E, Id>, l: number, r: number, a: number, b: number, x: Id): void {
    if (l === a && r === b) {
      this._propagate(node, x, r - l)
      if (l + 1 === r) node[3] = this._id()
      return
    }
    const mid = (l + r) >>> 1
    this._pushDown(node, l, r)
    if (a < mid) {
      let leftChild = node[0]
      if (!leftChild) {
        leftChild = this._newNode(l, mid)
        node[0] = leftChild
      }
      this._update(leftChild, l, mid, a, b < mid ? b : mid, x)
    }
    if (mid < b) {
      let rightChild = node[1]
      if (!rightChild) {
        rightChild = this._newNode(mid, r)
        node[1] = rightChild
      }
      this._update(rightChild, mid, r, a > mid ? a : mid, b, x)
    }
  }

  private _query(node: DNode<E, Id>, l: number, r: number, a: number, b: number): E {
    if (l === a && r === b) return node[2]
    const mid = (l + r) >>> 1
    this._pushDown(node, l, r)
    let res = this._eRange(l, r)
    if (a < mid && node[0]) res = this._op(res, this._query(node[0], l, mid, a, b < mid ? b : mid))
    if (mid < b && node[1]) res = this._op(res, this._query(node[1], mid, r, a > mid ? a : mid, b))
    return res
  }

  private _merge(
    node1: DNode<E, Id> | undefined,
    node2: DNode<E, Id> | undefined
  ): DNode<E, Id> | undefined {
    if (!node1 || !node2) return node1 || node2
    node1[0] = this._merge(node1[0], node2[0])
    node1[1] = this._merge(node1[1], node2[1])
    this._pushUp(node1)
    return node1
  }
}

export { SegmentTreeDynamicLazy }

if (require.main === module) {
  const n = 1e9
  const seg = new SegmentTreeDynamicLazy<number, number>(n, {
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
  seg.set(0, 1)
  console.log(seg.query(0, 1))

  const M = 1e5
  const pos = Array.from({ length: M }, () => Math.floor(Math.random() * n))
  console.time('start')
  for (let i = 0; i < M; ++i) {
    seg.set(pos[i], 1)
    seg.query(0, i)
    seg.update(0, i, 1)
    seg.get(i)
  }
  console.timeEnd('start') // start: 681.428ms

  // https://leetcode.cn/problems/range-module/
  // RangeAssignRangeSum

  class RangeModule {
    private readonly _seg = new SegmentTreeDynamicLazy<[sum: number, size: number], number>(1e9, {
      e() {
        return [0, 0]
      },
      eRange(start, end) {
        return [0, end - start]
      },
      id() {
        return -1
      },
      op(x, y) {
        return [x[0] + y[0], x[1] + y[1]]
      },
      mapping(f, x) {
        return f === -1 ? x : [f * x[1], x[1]]
      },
      composition(f, g) {
        return f === -1 ? g : f
      }
    })

    addRange(left: number, right: number): void {
      this._seg.update(left, right, 1)
    }

    queryRange(left: number, right: number): boolean {
      return this._seg.query(left, right)[0] === right - left
    }

    removeRange(left: number, right: number): void {
      this._seg.update(left, right, 0)
    }
  }

  /**
   * Your RangeModule object will be instantiated and called as such:
   * var obj = new RangeModule()
   * obj.addRange(left,right)
   * var param_2 = obj.queryRange(left,right)
   * obj.removeRange(left,right)
   */
  //   ["RangeModule", "addRange", "removeRange", "queryRange", "queryRange", "queryRange"]
  // [[], [10, 20], [14, 16], [10, 14], [13, 15], [16, 17]]

  // 来源：力扣（LeetCode）
  // 链接：https://leetcode.cn/problems/range-module
  // 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
  const RM = new RangeModule()
  RM.addRange(10, 20)
  RM.removeRange(14, 16)
  console.log(RM.queryRange(10, 14))
  console.log(RM.queryRange(13, 15))
  console.log(RM.queryRange(16, 17))
}

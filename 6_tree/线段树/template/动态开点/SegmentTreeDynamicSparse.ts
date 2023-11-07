/* eslint-disable no-inner-declarations */
// 单点修改, 区间查询
// 大多数位置的元素始终是单位元素的动态开点线段树.(非常稀疏)
// !其优点是不使用持久化时, 节点数可以保持在 O(N) 左右.
// 当持久化时时，可能会变得更慢。

// API
//  new SegmentTreeDynamicSparse(start, end, e, op, persistent)
//  newRoot()
//  get(index, root)
//  set(index, value, root)
//  update(index, value, root)
//  query(start, end, root)
//  queryAll(root)
//  maxRight(start, check, root)
//  minLeft(end, check, root)

/**
 * !不要用数组`[]`代替对象`{}`.数组会导致性能下降.
 */
type SegNode<E> = {
  left: SegNode<E> | undefined
  right: SegNode<E> | undefined
  index: number
  data: E
  sum: E
}

class SegmentTreeDynamic<E = number> {
  private static _isPrimitive(o: unknown): o is number | string | boolean | symbol | bigint | null | undefined {
    return o === null || (typeof o !== 'object' && typeof o !== 'function')
  }

  private readonly _lower: number
  private readonly _upper: number
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _persistent: boolean
  private _root: SegNode<E>

  /**
   * 单点修改区间查询的动态开点线段树.线段树维护的值域为`[start, end)`.
   * @param start 值域下界.start>=0.
   * @param end 值域上界.
   * @param e 幺元.
   * @param op 结合律的二元操作.
   * @param persistent 是否持久化.持久化后,每次修改都会新建一个结点,否则会复用原来的结点.
   */
  constructor(start: number, end: number, e: () => E, op: (a: E, b: E) => E, persistent = false) {
    if (persistent && !SegmentTreeDynamic._isPrimitive(e())) {
      throw new Error('persistent is only supported when e() return primitive values')
    }
    this._lower = start
    this._upper = end + 5
    this._e = e
    this._op = op
    this._persistent = persistent
    this._root = this.newRoot()
  }

  newRoot(): SegNode<E> {
    return undefined as any // nil
  }

  get(index: number, root: SegNode<E> = this._root): E {
    if (index < this._lower || index >= this._upper) return this._e()
    return this._get(root, index)
  }

  set(index: number, value: E, root: SegNode<E> = this._root): SegNode<E> {
    if (index < this._lower || index >= this._upper) return root
    const newRoot = this._set(root, this._lower, this._upper, index, value)
    this._root = newRoot
    return newRoot
  }

  update(index: number, value: E, root: SegNode<E> = this._root): SegNode<E> {
    if (index < this._lower || index >= this._upper) return root
    const newRoot = this._update(root, this._lower, this._upper, index, value)
    this._root = newRoot
    return newRoot
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * {@link _lower} <= start <= end <= {@link _upper}.
   */
  query(start: number, end: number, root: SegNode<E> = this._root): E {
    if (start < this._lower) start = this._lower
    if (end > this._upper) end = this._upper
    if (start >= end) return this._e()

    let res = this._e()
    const _query = (node: SegNode<E> | undefined, l: number, r: number, ql: number, qr: number) => {
      if (!node) return
      ql = l > ql ? l : ql
      qr = r < qr ? r : qr
      if (ql >= qr) return
      if (l === ql && r === qr) {
        res = this._op(res, node.sum)
        return
      }
      const m = Math.floor(l + (r - l) / 2)
      _query(node.left, l, m, ql, qr)
      if (ql <= node.index && node.index < qr) {
        res = this._op(res, node.data)
      }
      _query(node.right, m, r, ql, qr)
    }

    _query(root, this._lower, this._upper, start, end)
    return res
  }

  queryAll(root: SegNode<E> = this._root): E {
    return this.query(this._lower, this._upper, root)
  }

  /**
   * 二分查询最大的`end`使得切片`[start:end)`内的聚合值满足`check`.
   * {@link _lower} <= start <= {@link _upper}.
   * @alias findFirst
   */
  maxRight(start: number, check: (e: E) => boolean, root: SegNode<E> = this._root): number {
    if (start < this._lower) start = this._lower
    if (start >= this._upper) return this._upper

    let x = this._e()
    const _maxRight = (node: SegNode<E> | undefined, l: number, r: number, ql: number): number => {
      if (!node || r <= ql) return this._upper
      const tmp = this._op(x, node.sum)
      if (check(tmp)) {
        x = tmp
        return this._upper
      }
      const m = Math.floor(l + (r - l) / 2)
      const k = _maxRight(node.left, l, m, ql)
      if (k !== this._upper) return k
      if (ql <= node.index) {
        x = this._op(x, node.data)
        if (!check(x)) {
          return node.index
        }
      }
      return _maxRight(node.right, m, r, ql)
    }

    return _maxRight(root, this._lower, this._upper, start)
  }

  /**
   * 二分查询最小的`start`使得切片`[start:end)`内的聚合值满足`check`.
   * {@link _lower} <= end <= {@link _upper}.
   * @alias findLast
   */
  minLeft(end: number, check: (e: E) => boolean, root: SegNode<E> = this._root): number {
    if (end > this._upper) end = this._upper
    if (end <= this._lower) return this._lower

    let x = this._e()
    const _minLeft = (node: SegNode<E> | undefined, l: number, r: number, qr: number): number => {
      if (!node || qr <= l) return this._lower
      const tmp = this._op(node.sum, x)
      if (check(tmp)) {
        x = tmp
        return this._lower
      }
      const m = Math.floor(l + (r - l) / 2)
      const k = _minLeft(node.right, m, r, qr)
      if (k !== this._lower) return k
      if (node.index < qr) {
        x = this._op(node.data, x)
        if (!check(x)) {
          return node.index + 1
        }
      }
      return _minLeft(node.left, l, m, qr)
    }

    return _minLeft(root, this._lower, this._upper, end)
  }

  getAll(root: SegNode<E> = this._root): [index: number, value: E][] {
    const res: [number, E][] = []
    const _getAll = (node: SegNode<E> | undefined) => {
      if (!node) return
      _getAll(node.left)
      res.push([node.index, node.data])
      _getAll(node.right)
    }
    _getAll(root)
    return res
  }

  copy(node: SegNode<E>): SegNode<E> {
    if (!node || !this._persistent) return node
    return { left: node.left, right: node.right, index: node.index, data: node.data, sum: node.sum }
  }

  private _get(root: SegNode<E> | undefined, index: number): E {
    if (!root) return this._e()
    if (index === root.index) return root.data
    if (index < root.index) return this._get(root.left, index)
    return this._get(root.right, index)
  }

  private _set(root: SegNode<E> | undefined, l: number, r: number, i: number, x: E): SegNode<E> {
    if (!root) return SegmentTreeDynamic._newNode(i, x)
    root = this.copy(root)
    if (root.index === i) {
      root.data = x
      this._pushUp(root)
      return root
    }
    const m = Math.floor(l + (r - l) / 2)
    if (i < m) {
      if (root.index < i) {
        const tmp1 = root.index
        root.index = i
        i = tmp1
        const tmp2 = root.data
        root.data = x
        x = tmp2
      }
      root.left = this._set(root.left, l, m, i, x)
    } else {
      if (i < root.index) {
        const tmp1 = root.index
        root.index = i
        i = tmp1
        const tmp2 = root.data
        root.data = x
        x = tmp2
      }
      root.right = this._set(root.right, m, r, i, x)
    }
    this._pushUp(root)
    return root
  }

  private _pushUp(root: SegNode<E>): void {
    root.sum = root.data
    if (root.left) root.sum = this._op(root.left.sum, root.sum)
    if (root.right) root.sum = this._op(root.sum, root.right.sum)
  }

  private _update(root: SegNode<E> | undefined, l: number, r: number, i: number, x: E): SegNode<E> {
    if (!root) return SegmentTreeDynamic._newNode(i, x)
    root = this.copy(root)
    if (root.index === i) {
      root.data = this._op(root.data, x)
      this._pushUp(root)
      return root
    }
    const m = Math.floor(l + (r - l) / 2)
    if (i < m) {
      if (root.index < i) {
        const tmp1 = root.index
        root.index = i
        i = tmp1
        const tmp2 = root.data
        root.data = x
        x = tmp2
      }
      root.left = this._update(root.left, l, m, i, x)
    } else {
      if (i < root.index) {
        const tmp1 = root.index
        root.index = i
        i = tmp1
        const tmp2 = root.data
        root.data = x
        x = tmp2
      }
      root.right = this._update(root.right, m, r, i, x)
    }
    this._pushUp(root)
    return root
  }

  private static _newNode<V>(index: number, value: V): SegNode<V> {
    return {
      index,
      left: undefined,
      right: undefined,
      data: value,
      sum: value
    }
  }
}

export { SegmentTreeDynamic }

if (require.main === module) {
  const seg = new SegmentTreeDynamic<number>(
    0,
    10,
    () => 0,
    (a, b) => a + b
  )

  seg.update(0, 1)
  seg.set(0, 23)
  console.log(seg.getAll(), seg.queryAll())
  seg.update(0, 1)
  console.log(seg.getAll(), seg.queryAll())
  seg.set(0, 1)
  seg.update(0, 1)
  console.log(seg.getAll(), seg.queryAll())

  // check with queryAll
  for (let i = 0; i <= 10; i++) {
    seg.set(i, i)
    if (seg.queryAll() !== seg.query(0, 100)) {
      console.log(seg.queryAll(), seg.query(0, 100))
      throw new Error('queryAll failed')
    }
  }

  // https://leetcode.cn/problems/maximum-number-of-jumps-to-reach-the-last-index/
  // 6899. 达到末尾下标所需的最大跳跃次数
  // 给你一个下标从 0 开始、由 n 个整数组成的数组 nums 和一个整数 target 。
  // 你的初始位置在下标 0 。在一步操作中，你可以从下标 i 跳跃到任意满足下述条件的下标 j ：
  // 0 <= i < j < n
  // -target <= nums[j] - nums[i] <= target
  // 返回到达下标 n - 1 处所需的 最大跳跃次数 。
  // 如果无法到达下标 n - 1 ，返回 -1 。
  //
  // !O(nlog(U)), 线段树维护`值域最大值`
  function maximumJumps(nums: number[], target: number): number {
    const INF = 2e15
    const n = nums.length
    const seg = new SegmentTreeDynamic<number>(-3e9, 3e9, () => -INF, Math.max)
    seg.update(nums[0], 0)
    for (let i = 1; i < n; i++) {
      const preMax = seg.query(nums[i] - target, nums[i] + target + 1)
      seg.update(nums[i], preMax + 1)
    }
    const res = seg.get(nums[n - 1])
    return res >= 0 ? res : -1
  }
}

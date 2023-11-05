export {}

// 给你一个下标从 0 开始的整数数组 nums 。

// nums 一个长度为 k 的 子序列 指的是选出 k 个 下标 i0 < i1 < ... < ik-1 ，如果这个子序列满足以下条件，我们说它是 平衡的 ：

// 对于范围 [1, k - 1] 内的所有 j ，nums[ij] - nums[ij-1] >= ij - ij-1 都成立。
// nums 长度为 1 的 子序列 是平衡的。

// 请你返回一个整数，表示 nums 平衡 子序列里面的 最大元素和 。

// 一个数组的 子序列 指的是从原数组中删除一些元素（也可能一个元素也不删除）后，剩余元素保持相对顺序得到的 非空 新数组。

function maxBalancedSubsequenceSum(nums: number[]): number {
  if (nums.length === 1) return nums[0]
  const max = Math.max(...nums)
  if (max <= 0) return max

  const n = nums.length
  const newNums = nums.map((v, i) => v - i)

  const dp = new SegmentTreeDynamic(0, max + n + 10, () => 0, Math.max)
  const res = Array(n).fill(0)

  for (let i = 0; i < n; i++) {
    if (nums[i] <= 0) continue
    const key = newNums[i] + n
    const num = nums[i]
    const preMax = dp.query(0, key + 1)
    res[i] = preMax + num
    dp.update(key, res[i])
  }

  return Math.max(...res)
}

type SegNode<E> = {
  left: SegNode<E> | undefined
  right: SegNode<E> | undefined
  index: number
  data: E
  sum: E
}

class SegmentTreeDynamic<E = number> {
  private static _isPrimitive(
    o: unknown
  ): o is number | string | boolean | symbol | bigint | null | undefined {
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

  private _copyNode(node: SegNode<E>): SegNode<E> {
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
    root = this._copyNode(root)
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
    root = this._copyNode(root)
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

if (require.main === module) {
  // nums = [5,-1,-3,8]
  // console.log(maxBalancedSubsequenceSum([5, -1, -3, 8]))
  // nums = [3,3,5,6]
  // console.log(maxBalancedSubsequenceSum([3, 3, 5, 6]))
  // [3,5,-5,38,15,43,39,45,-21,-31,-17,-14,0,46,-32,-10,47,43,31,43,35,31,7,14,-33,-31]

  // 134
  console.log(
    maxBalancedSubsequenceSum([
      3, 5, -5, 38, 15, 43, 39, 45, -21, -31, -17, -14, 0, 46, -32, -10, 47, 43, 31, 43, 35, 31, 7,
      14, -33, -31
    ])
  )
}

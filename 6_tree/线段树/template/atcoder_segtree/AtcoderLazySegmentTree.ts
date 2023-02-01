/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */
// !由于lazy模板通用性 效率不如自己维护数组的线段树
// !注意如果是单点查询,可以去掉所有pushUp函数逻辑(js使用bigint会比较慢)
// !如果是单点修改,可以去掉所有懒标记逻辑

const INF = 2e15
const n = 10

if (require.main === module) {
  // Range Add Range Max
  const rangeAddRangeMax = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => 0,
    id: () => 0,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => f + x,
    composition: (f, g) => f + g
  })

  // Range Add Range Min
  const rangeAddRangeMin = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => INF,
    id: () => 0,
    op: (a, b) => Math.min(a, b),
    mapping: (f, x) => f + x,
    composition: (f, g) => f + g
  })

  // Range Chmax Range Max
  const rangeChmaxRangeMax = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => 0,
    id: () => -INF,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => (f === -INF ? x : Math.max(f, x)),
    composition: (f, g) => (f === -INF ? g : Math.max(f, g))
  })

  // Range Chmin Range Min
  const rangeChminRangeMin = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => INF,
    id: () => INF,
    op: (a, b) => Math.min(a, b),
    mapping: (f, x) => (f === INF ? x : Math.min(f, x)),
    composition: (f, g) => (f === INF ? g : Math.min(f, g))
  })

  // Range Assign (0/1) Range Sum
  const rangeAssignRangeSum = useAtcoderLazySegmentTree<[sum: number, size: number], number>(10, {
    e: () => [0, 1],
    id: () => -1,
    op: ([sum1, size1], [sum2, size2]) => [sum1 + sum2, size1 + size2],
    mapping: (f, [sum, size]) => (f === -1 ? [sum, size] : [f * size, size]),
    composition: (f, g) => (f === -1 ? g : f)
  })

  // Range Assign Range Max
  const rangeAssignRangeMax = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => 0,
    id: () => -INF,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => (f === -INF ? x : f),
    composition: (f, g) => (f === -INF ? g : f)
  })

  // Range Assign Range Min
  const rangeAssignRangeMin = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => INF,
    id: () => INF,
    op: (a, b) => Math.min(a, b),
    mapping: (f, x) => (f === INF ? x : f),
    composition: (f, g) => (f === INF ? g : f)
  })

  // Range Assign Point Get
  const rangeAssignPointGet = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => INF,
    id: () => INF,
    op: (a, b) => (a === INF ? b : a),
    mapping: (f, x) => (f === INF ? x : f),
    composition: (f, g) => (f === INF ? g : f)
  })
}

/**
 * E 线段树维护的值的类型
 * Id 更新操作的值的类型/懒标记的值的类型
 */
interface Operation<E, Id> {
  /**
   * 线段树维护的值的幺元
   */
  e: (this: void) => E

  /**
   * 更新操作/懒标记的幺元
   */
  id: (this: void) => Id

  /**
   * 合并左右区间的值
   */
  op: (this: void, data1: E, data2: E) => E

  /**
   * 父结点的懒标记更新子结点的值
   */
  mapping: (this: void, parentLazy: Id, childData: E) => E

  /**
   * 父结点的懒标记更新子结点的懒标记(合并)
   */
  composition: (this: void, parentLazy: Id, childLazy: Id) => Id
}

interface AtcoderSegmentTree<E, Id> {
  /**
   * 查询切片 `[left:right]` 内的值
   *
   * `0 <= left <= right <= n`
   */
  query: (left: number, right: number) => E

  queryAll: () => E

  /**
   * 更新切片 `[left:right]` 内的值
   *
   * `0 <= left <= right <= n`
   */
  update: (left: number, right: number, value: Id) => void

  /**
   * 树上二分查询最大的 `right` 使得切片 `[left:right]` 内的值满足 `predicate`
   */
  maxRight: (left: number, predicate: (value: E) => boolean) => number

  /**
   * 树上二分查询最小的 `left` 使得切片 `[left:right]` 内的值满足 `predicate`
   */
  minLeft: (right: number, predicate: (value: E) => boolean) => number
}

/**
 * @see {@link https://betrue12.hateblo.jp/entry/2020/09/22/194541}
 */
function useAtcoderLazySegmentTree<E = number, Id = number>(
  sizeOrArray: number | ArrayLike<E>,
  operation: Operation<E, Id>
): AtcoderSegmentTree<E, Id> {
  // !bind会导致性能损失 因此这里统一不使用bind 接口里声明this为void
  const _e = operation.e
  const _id = operation.id
  const _op = operation.op
  const _mapping = operation.mapping
  const _composition = operation.composition

  const _n = typeof sizeOrArray === 'number' ? sizeOrArray : sizeOrArray.length
  const _log = 32 - Math.clz32(_n - 1)
  const _size = 1 << _log
  const _data = Array<E>(_size * 2).fill(_e())
  const _lazy = Array<Id>(_size * 2).fill(_id())

  if (Array.isArray(sizeOrArray)) {
    for (let i = 0; i < _n; i++) {
      _data[_size + i] = sizeOrArray[i]
    }
  }

  for (let i = _size - 1; i > 0; i--) {
    _pushUp(i)
  }

  function query(left: number, right: number): E {
    if (left < 0) left = 0
    if (right > _n) right = _n
    if (left >= right) return _e()

    left += _size
    right += _size
    for (let i = _log; i > 0; i--) {
      if ((left >> i) << i !== left) _pushDown(left >> i)
      if ((right >> i) << i !== right) _pushDown(right >> i)
    }

    let leftRes = _e()
    let rightRes = _e()
    while (left < right) {
      if (left & 1) {
        leftRes = _op(leftRes, _data[left])
        left++
      }

      if (right & 1) {
        right--
        rightRes = _op(_data[right], rightRes)
      }

      left >>= 1
      right >>= 1
    }

    return _op(leftRes, rightRes)
  }

  function queryAll(): E {
    return _data[1]
  }

  function update(left: number, right: number, value: Id): void {
    if (left < 0) left = 0
    if (right > _n) right = _n
    if (left >= right) return

    left += _size
    right += _size
    for (let i = _log; i > 0; i--) {
      if ((left >> i) << i !== left) _pushDown(left >> i)
      if ((right >> i) << i !== right) _pushDown((right - 1) >> i)
    }

    let preLeft = left
    let preRight = right
    while (left < right) {
      if (left & 1) {
        _propagate(left, value)
        left++
      }

      if (right & 1) {
        right--
        _propagate(right, value)
      }

      left >>= 1
      right >>= 1
    }

    left = preLeft
    right = preRight
    for (let i = 1; i < _log + 1; i++) {
      if ((left >> i) << i !== left) _pushUp(left >> i)
      if ((right >> i) << i !== right) _pushUp((right - 1) >> i)
    }
  }

  function maxRight(left: number, predicate: (value: E) => boolean): number {
    _checkBoundsBeginEnd(left, left)
    if (left === _n) return _n

    left += _size
    for (let i = _log; i > 0; i--) {
      _pushDown(left >> i)
    }

    let res = _e()
    while (true) {
      while (left % 2 === 0) left >>= 1
      if (!predicate(_op(res, _data[left]))) {
        while (left < _size) {
          _pushDown(left)
          left *= 2
          if (predicate(_op(res, _data[left]))) {
            res = _op(res, _data[left])
            left++
          }
        }

        return left - _size
      }

      res = _op(res, _data[left])
      left++
      if ((left & -left) === left) break
    }

    return _n
  }

  function minLeft(right: number, predicate: (value: E) => boolean): number {
    _checkBoundsBeginEnd(right, right)
    if (right === 0) return 0

    right += _size
    for (let i = _log; i > 0; i--) {
      _pushDown((right - 1) >> i)
    }

    let res = _e()
    while (true) {
      right--
      while (right > 1 && right % 2) right >>= 1
      if (!predicate(_op(_data[right], res))) {
        while (right < _size) {
          _pushDown(right)
          right = right * 2 + 1
          if (predicate(_op(_data[right], res))) {
            res = _op(_data[right], res)
            right--
          }
        }

        return right + 1 - _size
      }

      res = _op(_data[right], res)
      if ((right & -right) === right) break
    }

    return 0
  }

  return {
    query,
    queryAll,
    update,
    maxRight,
    minLeft
  }

  function _pushUp(root: number): void {
    _data[root] = _op(_data[root * 2], _data[root * 2 + 1])
  }

  function _pushDown(root: number): void {
    _propagate(2 * root, _lazy[root])
    _propagate(2 * root + 1, _lazy[root])
    _lazy[root] = _id()
  }

  // 更新子结点的lazy和data
  function _propagate(root: number, parentLazy: Id): void {
    _data[root] = _mapping(parentLazy, _data[root])
    // !叶子结点不需要更新lazy
    if (root < _size) {
      _lazy[root] = _composition(parentLazy, _lazy[root])
    }
  }

  function _checkBoundsBeginEnd(begin: number, end: number): void {
    if (begin < 0 || begin > end || end > _n) {
      throw new RangeError(`Invalid range: [${begin}, ${end}) out of [0, ${_n})`)
    }
  }
}

export { useAtcoderLazySegmentTree, Operation, AtcoderSegmentTree } // base api

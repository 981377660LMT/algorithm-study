/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */

/**
 * S 线段树维护的值的类型
 *
 * F 操作更新的值的类型/懒标记的值的类型
 */
interface Operation<S, F> {
  /**
   * 线段树维护的值的幺元
   * @alias e
   */
  dataUnit: () => S

  /**
   * 懒标记的幺元
   * @alias id
   */
  lazyUnit: () => F

  /**
   * 合并左右区间的值
   * @alias op
   */
  mergeChildren: (data1: S, data2: S) => S

  /**
   * 父结点的懒标记更新子结点的值
   * @alias mapping
   */
  updateData: (parentLazy: F, childData: S) => S

  /**
   * 父结点的懒标记更新子结点的懒标记
   * @alias composition
   */
  updateLazy: (parentLazy: F, childLazy: F) => F
}

interface AtcoderSegmentTree<S, F> {
  /**
   * 查询切片 `[left:right]` 内的值
   *
   * `0 <= left <= right <= n`
   * @alias prod
   */
  query: (left: number, right: number) => S

  /**
   * @alias all_prod
   */
  queryAll: () => S

  /**
   * 更新切片 `[left:right]` 内的值
   *
   * `0 <= left <= right <= n`
   * @alias apply
   */
  update: (left: number, right: number, value: F) => void

  /**
   * 树上二分查询最大的 `right` 使得切片 `[left:right]` 内的值满足 `predicate`
   */
  maxRight: (left: number, predicate: (value: S) => boolean) => number

  /**
   * 树上二分查询最小的 `left` 使得切片 `[left:right]` 内的值满足 `predicate`
   */
  minLeft: (right: number, predicate: (value: S) => boolean) => number
}

/**
 * @see {@link https://betrue12.hateblo.jp/entry/2020/09/22/194541}
 *      {@link https://github/shakayami/ACL-for-python/lazysegtree}
 */
function useAtcoderSegmentTree<S, F>(
  sizeOrArray: number | ArrayLike<S>,
  operation: Operation<S, F>
): AtcoderSegmentTree<S, F> {
  const _n = typeof sizeOrArray === 'number' ? sizeOrArray : sizeOrArray.length
  const _log = Math.ceil(Math.log2(_n))
  const _size = 1 << _log
  const _data = Array<S>(_size * 2).fill(operation.dataUnit())
  const _lazy = Array<F>(_size * 2).fill(operation.lazyUnit())
  const _e = operation.dataUnit
  const _id = operation.lazyUnit
  const _op = operation.mergeChildren
  const _mapping = operation.updateData
  const _composition = operation.updateLazy

  if (Array.isArray(sizeOrArray)) {
    for (let i = 0; i < _n; i++) {
      _data[_size + i] = sizeOrArray[i]
    }
  }

  for (let i = _size - 1; i > 0; i--) {
    _pushUp(i)
  }

  function query(left: number, right: number): S {
    _checkRange(left, right)
    if (left === right) return _e()

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

  function queryAll(): S {
    return _data[1]
  }

  function update(left: number, right: number, value: F): void {
    _checkRange(left, right)
    if (left === right) return

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
        _allApply(left, value)
        left++
      }

      if (right & 1) {
        right--
        _allApply(right, value)
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

  function maxRight(left: number, predicate: (value: S) => boolean): number {
    _checkRange(left, left)
    if (!predicate(_e())) return left
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

  function minLeft(right: number, predicate: (value: S) => boolean): number {
    _checkRange(right, right)
    if (!predicate(_e())) return right
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
    _allApply(2 * root, _lazy[root])
    _allApply(2 * root + 1, _lazy[root])
    _lazy[root] = _id()
  }

  // pushDown辅助函数 更新子结点的lazy和data
  function _allApply(root: number, parentLazy: F): void {
    _data[root] = _mapping(parentLazy, _data[root])
    if (root < _size) {
      _lazy[root] = _composition(parentLazy, _lazy[root])
    }
  }

  function _checkRange(left: number, right: number): void {
    if (left >= 0 && left <= right && right <= _n) return
    throw new RangeError(`Invalid range: [${left}, ${right}) out of [0, ${_n})`)
  }
}

export { useAtcoderSegmentTree, Operation, AtcoderSegmentTree }

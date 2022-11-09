/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

/**
 * 权值树状数组模拟名次树
 *
 * @param size 最大值 没有使用离散化，不能超出3e7
 * @param offset 偏移量(存在负数时，所有数据都要加上偏移量)
 */
function useBITArray(size: number, offset = 0) {
  const _max = size + offset
  const _log = 32 - Math.clz32(_max) // 存储的值域需要为二的幂次
  const _upper = 1 << _log
  const _tree = new Uint32Array(_upper + 1)
  const _counter = new Uint32Array(_upper + 1)
  let _size = 0

  // 树状数组上二分询问第k大
  function at(index: number): number {
    index++
    let left = 1
    let right = _upper
    while (left ^ right) {
      const mid = (left + right) >>> 1
      if (_tree[mid] < index) {
        index -= _tree[mid]
        left = mid + 1
      } else {
        right = mid
      }
    }

    return left - offset
  }

  function add(value: number): void {
    const index = value + offset
    _add(index, 1)
    _counter[index]++
    _size++
  }

  function remove(value: number): void {
    const index = value + offset
    if (_counter[index] === 0) {
      return
    }

    _add(index, -1)
    _counter[index]--
    _size--
  }

  function bisectLeft(value: number): number {
    return _query(value + offset - 1)
  }

  function bisectRight(value: number): number {
    return _query(value + offset)
  }

  return {
    add,
    remove,
    bisectLeft,
    bisectRight,
    at,
    get size() {
      return _size
    }
  }

  function _add(index: number, delta: number): void {
    while (index <= _upper) {
      _tree[index] += delta
      index += index & -index
    }
  }

  function _query(index: number): number {
    let res = 0
    while (index > 0) {
      res += _tree[index]
      index -= index & -index
    }
    return res
  }
}

export {}

if (require.main === module) {
  const foo = useBITArray(1e6)
  foo.add(1)
  foo.add(1)
  foo.add(2)
  foo.add(2)
  foo.add(3)
  console.log(foo.bisectLeft(1))
  console.log(foo.bisectLeft(2))
  console.log(foo.bisectRight(1))
  console.log(foo.at(0))
  console.log(foo.at(1))
  console.log(foo.at(2))
  console.log(foo.at(3))
  console.log(foo.at(4))
}

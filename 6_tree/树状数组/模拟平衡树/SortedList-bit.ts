/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
// SortedList

/**
 * 权值树状数组模拟名次树,需要预先离散化.
 *
 * @param max 最大值 没有使用离散化时不能超出1e6
 */
function useSortedList(max: number) {
  const _log = 32 - Math.clz32(max) // 存储的值域需要为二的幂次
  const _upper = 1 << _log
  const _tree = new Uint32Array(_upper + 1)
  const _counter = new Uint32Array(_upper + 1)
  let _size = 0

  function at(index: number): number | undefined {
    index++
    if (index < 1) index += _size
    if (index < 1 || index > _size) return undefined
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

    return left - 1
  }

  function add(value: number): void {
    _add(value, 1)
    _counter[value]++
    _size++
  }

  function discard(value: number): void {
    if (_counter[value] === 0) {
      return
    }
    _counter[value]--
    _add(value, -1)
    _size--
  }

  // < value
  function bisectLeft(value: number): number {
    return _query(value - 1)
  }

  // <= value
  function bisectRight(value: number): number {
    return _query(value)
  }

  return {
    add,
    discard,
    bisectLeft,
    bisectRight,
    at,
    get size() {
      return _size
    },
    toString() {
      const res: number[] = []
      _counter.forEach((count, index) => {
        if (count > 0) {
          for (let i = 0; i < count; i++) {
            res.push(index)
          }
        }
      })
      return `SortedList{${res.toString()}}`
    }
  }

  function _add(index: number, delta: number): void {
    index++
    for (; index <= _upper; index += index & -index) {
      _tree[index] += delta
    }
  }

  function _query(index: number): number {
    index++
    let res = 0
    for (; index > 0; index -= index & -index) {
      res += _tree[index]
    }
    return res
  }
}

/**
 * (松)离散化.
 * @returns
 * rank: 给定一个数,返回它的排名`(0-count)`.
 * count: 离散化(去重)后的元素个数.
 */
function discretize(nums: number[]): [rank: (num: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const rank = (num: number) => {
    let left = 0
    let right = allNums.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (allNums[mid] >= num) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }
  return [rank, allNums.length]
}

export {}

if (require.main === module) {
  // https://leetcode.cn/problems/count-the-number-of-fair-pairs/
  // lower <= nums[i] + nums[j] <= upper
  function countFairPairs(nums: number[], lower: number, upper: number): number {
    // 离散化
    const allNums = nums.slice()
    nums.forEach(num => {
      allNums.push(lower - num, upper - num)
    })
    const [rank, n] = discretize(allNums)

    const sl = useSortedList(n)
    let res = 0
    nums.forEach(num => {
      const right = sl.bisectRight(rank(upper - num))
      const left = sl.bisectLeft(rank(lower - num))
      res += right - left
      sl.add(rank(num))
    })
    return res
  }
}

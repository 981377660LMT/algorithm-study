/* eslint-disable max-len */

/**
 * (松)离散化.
 * @param offset 离散化后的排名偏移量.
 * @returns
 * - getRank: 给定一个数,返回它的排名`(offset ~ offset + count)`.
 * - count: 离散化(去重)后的元素个数.
 */
function discretizeSparse(nums: number[], offset = 0): [getRank: (num: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)

  // bisect_left
  const getRank = (num: number): number => {
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
    return left + offset
  }
  return [getRank, allNums.length]
}

/**
 * (紧)离散化.
 * @param offset 离散化后的排名偏移量.
 * @returns
 * - getRank: 给定一个数,返回它的排名`(offset ~ offset + count)`.
 * - count: 离散化(去重)后的元素个数.
 */
function discretizeCompressed(nums: number[], offset = 0): [getRank: (num: number) => number, getValue: (rank: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const mp = new Map<number, number>()
  for (let index = 0; index < allNums.length; index++) mp.set(allNums[index], index + offset)
  const getRank = (num: number) => mp.get(num)!
  const getValue = (rank: number) => allNums[rank - offset]
  return [getRank, getValue, allNums.length]
}

/**
 * 不带相同值的离散化,转换为`0-n-1`.
 * @returns
 * - rank: 离散化后的排名.
 * - keys: keys[ranks[i]] = nums[i].
 */
function discretizeUnique(nums: number[]): [rank: Uint32Array, keys: number[]] {
  let rank = argSort(nums)
  const keys = reArrage(nums, rank)
  rank = argSort(rank)
  return [rank, keys]
}

/**
 * 返回数组的排序索引.
 */
function argSort<T>(arr: ArrayLike<T>): Uint32Array {
  const n = arr.length
  const order = new Uint32Array(n)
  for (let i = 0; i < n; i++) order[i] = i
  order.sort((a, b) => (arr[a] < arr[b] ? -1 : 1))
  return order
}

/**
 * 按照索引数组重新排列数组.
 */
function reArrage<T>(arr: ArrayLike<T>, order: ArrayLike<number>): T[] {
  const n = arr.length
  const res = Array(n)
  for (let i = 0; i < n; i++) res[i] = arr[order[i]]
  return res
}

export { discretizeCompressed, discretizeSparse, discretizeUnique, argSort, reArrage }

if (require.main === module) {
  const nums = [3, 3, 1, 2, 2, 0, 4]
  console.log(discretizeUnique(nums))
}

/**
 * (紧)离散化.
 * @returns
 * - rank: 给定一个在 nums 中的值,返回它的排名`(0~rank.size-1)`.
 * - newNums: 离散化后的数组.
 */
function discretizeCompressed(nums: number[]): [rank: Map<number, number>, newNums: Uint32Array] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const rank = new Map<number, number>()
  for (let index = 0; index < allNums.length; index++) rank.set(allNums[index], index)
  const newNums = new Uint32Array(nums.length)
  for (let index = 0; index < nums.length; index++) newNums[index] = rank.get(nums[index])!
  return [rank, newNums]
}

/**
 * (松)离散化.
 * @returns
 * - rank: 给定一个数,返回它的排名`(0-count)`.
 * - count: 离散化(去重)后的元素个数.
 */
function discretizeSparse(nums: number[]): [rank: (num: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  // bisect_left
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

export { discretizeCompressed, discretizeSparse }

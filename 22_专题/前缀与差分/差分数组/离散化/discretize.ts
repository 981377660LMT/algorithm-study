/* eslint-disable max-len */

/**
 * (松)离散化.
 * @param offset 离散化后的排名偏移量.
 * @returns
 * - getRank: 给定一个数,返回它的排名`(offset ~ offset + count)`.
 * - count: 离散化(去重)后的元素个数.
 */
function discretizeSparse(
  nums: number[],
  offset = 0
): [getRank: (num: number) => number, count: number] {
  const allNums = nums.slice().sort((a, b) => a - b)
  // uniqueInplace
  let slow = 0
  for (let fast = 0; fast < allNums.length; fast++) {
    if (allNums[fast] !== allNums[slow]) allNums[++slow] = allNums[fast]
  }
  allNums.length = slow + 1

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
function discretizeCompressed(
  nums: number[],
  offset = 0
): [getRank: (num: number) => number, count: number] {
  const allNums = nums.slice().sort((a, b) => a - b)
  // uniqueInplace
  let slow = 0
  for (let fast = 0; fast < allNums.length; fast++) {
    if (allNums[fast] !== allNums[slow]) allNums[++slow] = allNums[fast]
  }
  allNums.length = slow + 1

  const mp = new Map<number, number>()
  for (let index = 0; index < allNums.length; index++) mp.set(allNums[index], index + offset)
  const getRank = (num: number) => mp.get(num)!
  return [getRank, allNums.length]
}

export { discretizeCompressed, discretizeSparse }

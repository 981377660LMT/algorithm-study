/**
 *
 * @param nums
 * @description
 * 给定一个非空整数数组，找到使所有数组元素相等所需的最小移动数，
 * 其中每次移动可将选定的一个元素加1或减1。
 * 您可以假设数组的长度最多为10000。
 */
function minMoves2(nums: number[]): number {
  if (nums.length <= 1) return 0
  const sorted = nums.slice().sort((a, b) => a - b)
  const mid = sorted[sorted.length >> 1]
  return sorted.reduce((pre, cur) => pre + Math.abs(cur - mid), 0)
}

console.log(minMoves2([1, 2, 3]))
console.log(minMoves2([1, 3, 2]))

import { smallestRange } from './632. 最小区间'

/**
 * @param {number[]} nums
 * @return {number}
 * @description 如果元素是 偶数 ，可以除以 2 如果元素是 奇数 ，可以乘上 2
 * 求任意两个元素间最大差值的最小值
 * 这题和632最小区间等价
 */
function minimumDeviation(nums: number[]): number {
  const grid = Array.from<number, number[]>({ length: nums.length }, () => [])

  for (let i = 0; i < nums.length; i++) {
    if (nums[i] % 2 === 1) {
      grid[i].push(nums[i], nums[i] * 2)
    } else {
      while (nums[i] % 2 === 0) {
        grid[i].push(nums[i])
        nums[i] /= 2
      }
      grid[i].push(nums[i])
      grid[i].reverse()
    }
  }

  const res = smallestRange(grid)
  return res[1] - res[0]
}

console.log(minimumDeviation([1, 2, 3, 4]))

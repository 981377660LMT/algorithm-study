/**
 * @param {number[]} nums  一个二进制数组 nums
 * @return {number}  返回该子数组的长度
 * 找到含有相同数量的 0 和 1 的最长连续子数组，并返回该子数组的长度。
 *
 * 等价于前缀和相等
 */
const findMaxLength = function (nums: number[]): number {
  const pre = new Map<number, number>([[0, -1]])
  let sum = 0
  let max = 0

  for (let i = 0; i < nums.length; i++) {
    const cur = nums[i]
    sum += cur === 0 ? -1 : 1
    if (pre.has(sum)) max = Math.max(max, i - pre.get(sum)!)
    else pre.set(sum, i)
  }

  return max
}

console.log(findMaxLength([0, 1, 0]))

export default 1

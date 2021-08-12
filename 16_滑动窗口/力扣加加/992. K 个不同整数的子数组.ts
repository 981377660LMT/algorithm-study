/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * @summary 题目转换成求：不超过k种元素的子数组个数 - 不超过k-1种元素的子数组个数(转化为水果成蓝问题)
 */
const subarraysWithKDistinct = function (nums: number[], k: number): number {
  const helper = (n: number): number => {
    let l = 0
    let res = 0
    const map = new Map<number, number>()

    for (let r = 0; r < nums.length; r++) {
      const cur = nums[r]
      map.set(cur, map.get(cur)! + 1 || 1)
      while (map.size > n) {
        l++
        const pre = nums[l - 1]
        const count = map.get(pre)!
        if (count === 1) map.delete(pre)
        else map.set(pre, count - 1)
      }
      res += r - l + 1
    }

    return res
  }

  return helper(k) - helper(k - 1)
}

console.log(subarraysWithKDistinct([1, 2, 1, 2, 3], 2))

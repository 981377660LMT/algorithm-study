/**
 * @param {number[]} nums nums[i] 不是 0 就是 1
 * @param {number} goal
 * @return {number}
 * @description 请你统计并返回有多少个和为 goal 的 非空 子数组。子数组 是数组的一段连续部分。
 * @summary 前缀和相差为goal 子数组和就为goal
 */
const numSubarraysWithSum = function (nums: number[], goal: number): number {
  // 此时的前缀和,出现了几个
  const map = new Map<number, number>([[0, 1]])
  let sum = 0
  let res = 0

  for (let i = 0; i < nums.length; i++) {
    sum += nums[i]
    const pre = sum - goal
    if (map.has(pre)) res += map.get(pre)!
    map.set(sum, (map.get(sum) ?? 0) + 1)
  }

  return res
}

console.log(numSubarraysWithSum([1, 0, 1, 0, 1], 2))
console.log(numSubarraysWithSum([0, 0, 0, 0, 0], 0))

export {}

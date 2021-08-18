/**
 * @param {number[]} nums 1 <= nums.length <= 105
 * @param {number} k k 的取值范围是 [1, 数组中不相同的元素的个数]
 * @return {number[]}
 * 给你一个整数数组 nums 和一个整数 k ，请你返回其中出现频率前 k 高的元素。
 */
const topKFrequent = function (nums: number[], k: number): number[] {
  const res = []
  const map = new Map()
  const bucket: Set<number>[] = []

  nums.forEach(num => map.set(num, map.get(num) + 1 || 1))

  for (const [num, count] of map) {
    bucket[count] = (bucket[count] || new Set()).add(num)
  }

  for (let i = bucket.length - 1; i >= 0; i--) {
    if (bucket[i]) res.push(...bucket[i])
    if (res.length === k) break
  }

  return res
}

export default 1

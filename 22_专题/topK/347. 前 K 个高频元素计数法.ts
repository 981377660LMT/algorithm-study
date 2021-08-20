/**
 * @param {number[]} nums 1 <= nums.length <= 105
 * @param {number} k k 的取值范围是 [1, 数组中不相同的元素的个数]
 * @return {number[]}
 * 给你一个整数数组 nums 和一个整数 k ，请你返回其中出现频率前 k 高的元素。
 */
const topKFrequent = function (nums: number[], k: number): number[] {
  const res = []
  const counter = new Map<number, number>()
  // 将频率作为数组下标
  const bucket: Set<number>[] = []

  nums.forEach(num => counter.set(num, (counter.get(num) || 0) + 1))

  for (const [num, freq] of counter) {
    bucket[freq] = (bucket[freq] || new Set()).add(num)
  }

  for (let i = bucket.length - 1; i >= 0; i--) {
    if (bucket[i]) res.push(...bucket[i])
    if (res.length === k) break
  }

  return res
}

export default 1

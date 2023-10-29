import { Heap } from '../../../../8_heap/Heap'

/**
 * 第k小的子集和.k从1开始.
 * n,k<=1e5.
 */
function kthSubsetSum(nums: number[], k: number): number {
  const n = nums.length
  let sum = 0
  for (let i = 0; i < n; i++) {
    const x = nums[i]
    if (x >= 0) {
      sum += x
    } else {
      nums[i] = -x
    }
  }
  nums.sort((a, b) => a - b)
  const pq = new Heap<{ sum: number; i: number }>((a, b) => b.sum - a.sum, [{ sum, i: 0 }])
  for (let _ = 0; _ < k - 1; _++) {
    const p = pq.pop()!
    if (p.i < n) {
      pq.push({ sum: p.sum - nums[p.i], i: p.i + 1 }) // 保留 nums[p.i-1]
      if (p.i > 0) {
        pq.push({ sum: p.sum - nums[p.i] + nums[p.i - 1], i: p.i + 1 }) // 不保留 nums[p.i-1]，把之前减去的加回来
      }
    }
  }
  return pq.peek()!.sum
}

export { kthSubsetSum }

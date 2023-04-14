/**
 * @param {number[]} nums
 * @return {boolean}
 * @description
 * 所有 nums[seq[j]] 应当不是 全正 就是 全负
 * k > 1
 */
function circularArrayLoop(nums: number[]): boolean {
  const n = nums.length
  // 注意这里索引模n的处理
  const next = (index: number) => (((index + nums[index]) % n) + n) % n

  const visited = new Uint8Array(n)
  for (let i = 0; i < n; i++) {
    if (visited[i]) continue
    let slow = i
    let fast = next(slow)
    visited[slow] = 1
    visited[fast] = 1
    while (nums[fast] * nums[i] > 0 && nums[next(fast)] * nums[i] > 0) {
      if (fast === slow) {
        // 只有一个点的情况
        if (slow === next(slow)) break
        else return true
      }
      slow = next(slow)
      fast = next(next(fast))
      visited[slow] = 1
      visited[fast] = 1
    }
  }

  return false
}

console.log(circularArrayLoop([2, -1, 1, 2, 2]))
console.log(circularArrayLoop([-2, 1, -1, -2, -2]))

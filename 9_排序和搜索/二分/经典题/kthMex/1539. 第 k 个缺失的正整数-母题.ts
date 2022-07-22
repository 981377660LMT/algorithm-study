// 给你一个 严格升序排列 的正整数数组 arr 和一个整数 k 。
// 请你找到这个数组里第 k 个缺失的正整数。
function findKthMex(nums: number[], k: number): number {
  let left = 0
  let right = nums.length - 1
  while (left <= right) {
    const mid = (left + right) >> 1
    const missing = nums[mid] - (mid + 1)
    if (missing >= k) right = mid - 1
    else left = mid + 1
  }

  return k + left
}

if (require.main === module) {
  console.log(
    findKthMex(
      [
        96, 44, 99, 25, 61, 84, 88, 18, 19, 33, 60, 86, 52, 19, 32, 47, 35, 50, 94, 17, 29, 98, 22,
        21, 72, 100, 40, 84,
      ].sort((a, b) => a - b),
      35
    )
  )
}

export { findKthMex }
// 解释：缺失的正整数包括 [1,5,6,8,9,10,12,13,...] 。第 5 个缺失的正整数为 9 。

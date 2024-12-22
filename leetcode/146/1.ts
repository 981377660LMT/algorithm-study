export {}
function countSubarrays(nums: number[]): number {
  const n = nums.length
  let res = 0
  for (let i = 0; i <= n - 3; i++) {
    if (nums[i] + nums[i + 2] === nums[i + 1] / 2) {
      res++
    }
  }
  return res
}

export {}
function constructTransformedArray(nums: number[]): number[] {
  const n = nums.length
  const res = Array(n)
  for (let i = 0; i < n; i++) {
    const move = nums[i]
    if (move === 0) {
      res[i] = nums[i]
    } else {
      const index = (((i + move) % n) + n) % n
      res[i] = nums[index]
    }
  }
  return res
}

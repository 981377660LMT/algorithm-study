function kLengthApart(nums: number[], k: number): boolean {
  let preIndex = -1

  for (let i = 0; i < nums.length; i++) {
    if (nums[i] !== 1) continue
    if (preIndex !== -1 && i - preIndex <= k) return false
    preIndex = i
  }

  return true
}

console.log(kLengthApart([1, 0, 0, 0, 1, 0, 0, 1], 2))

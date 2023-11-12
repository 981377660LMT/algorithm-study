const maximumStrongPairXor = function (nums) {
  let res = 0
  nums.sort((a, b) => a - b)
  for (let i = 0; i < nums.length; i++) {
    for (let j = i + 1; j < nums.length; j++) {
      if (nums[j] <= 2 * nums[i]) {
        res = Math.max(res, nums[i] ^ nums[j])
      }
    }
  }
  return res
}

// 轻量操作 js 可过 1e9， 3000ms+

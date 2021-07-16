const canJump = (nums: number[]) => {
  if (nums.length === 1) return 0
  let maxCanJumpIndex = 0
  let maxCanJumpCount = 0

  for (let curIndex = 0; curIndex < nums.length - 1; curIndex++) {
    if (maxCanJumpIndex >= nums.length - 1) return maxCanJumpCount
    const curValue = nums[curIndex]

    maxCanJumpIndex = Math.max(maxCanJumpIndex, curIndex + curValue)
  }
}

console.log(canJump([2, 3, 1, 1, 4]))
console.log(canJump([2, 3, 0, 1, 4]))

export {}

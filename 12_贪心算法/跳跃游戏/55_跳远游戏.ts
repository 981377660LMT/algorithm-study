// 判断你是否能够到达最后一个位置。
const canJump = (nums: number[]): boolean => {
  if (nums.length === 1) return true
  let maxCanJumpIndex = 0

  // 求最大覆盖范围
  for (let curIndex = 0; curIndex < nums.length - 1; curIndex++) {
    const curValue = nums[curIndex]
    maxCanJumpIndex = Math.max(maxCanJumpIndex, curIndex + curValue)

    // 肯定要到这个0的点，跳不出这个点则结束
    if (curValue === 0 && maxCanJumpIndex <= curIndex) return false
  }

  return true
}

console.log(canJump([2, 3, 1, 1, 4]))
console.log(canJump([3, 2, 1, 0, 4]))
console.log(canJump([2, 0, 0]))

export {}

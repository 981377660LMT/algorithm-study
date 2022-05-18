// 判断你是否能够到达最后一个位置。
const canJump = (nums: number[]): boolean => {
  if (nums.length === 1) return true
  let next = 0

  // 求最大覆盖范围
  for (let cur = 0; cur < nums.length - 1; cur++) {
    next = Math.max(next, cur + nums[cur])

    // 肯定要到这个0的点，跳不出这个点则结束
    if (nums[cur] === 0 && next <= cur) return false
  }

  return true
}

console.log(canJump([2, 3, 1, 1, 4]))
console.log(canJump([3, 2, 1, 0, 4]))
console.log(canJump([2, 0, 0]))

export {}

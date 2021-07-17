// O(n)
// https://blog.csdn.net/zemprogram/article/details/104071994
const canJump = (nums: number[]) => {
  const target = nums.length - 1
  let cur = 0
  let next = 0
  let step = 0

  for (let i = 0; i < target; i++) {
    if (cur >= target) break
    // 看哪一步最远
    next = Math.max(next, i + nums[i])
    // 需要做出决定，pre跳到最远的cur
    if (cur === i) {
      cur = next
      step++
    }
  }

  return step
}

console.log(canJump([2, 3, 1, 1, 4]))
console.log(canJump([2, 3, 0, 1, 4]))

export {}

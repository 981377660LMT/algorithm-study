// O(n)
// https://blog.csdn.net/zemprogram/article/details/104071994
// 你的目标是使用最少的跳跃次数到达数组的最后一个位置。
// 假设你总是可以到达数组的最后一个位置。

// 以最小的步数增加覆盖范围，覆盖范围一旦覆盖了终点，得到的就是最小步数！
const canJump = (nums: number[]) => {
  const target = nums.length - 1
  let cur = 0
  let next = 0
  let step = 0

  for (let i = 0; i < target; i++) {
    if (cur >= target) break
    // 看之前所有点中最大的覆盖范围
    next = Math.max(next, i + nums[i])
    // 需要做出决定，pre跳到最远的cur
    if (cur === i) {
      if (cur >= next) return -1 // 无法到达
      cur = next
      step++
    }
  }

  return step
}

if (require.main === module) {
  console.log(canJump([2, 3, 1, 1, 4]))
  console.log(canJump([2, 3, 0, 1, 4]))
  console.log(canJump([1, 0, 0]))
}

export {}

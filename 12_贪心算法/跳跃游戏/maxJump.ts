/**
 *
 * @param jumps  每个点处可达到的最远坐标
 * @param target 需要到达的点
 * @returns  需要的最小步数
 */
function maxJump(jumps: number[], target: number): number {
  let cur = 0
  let next = 0
  let step = 0

  for (let i = 0; i < target; i++) {
    if (cur >= target) break
    // 看之前所有点中最大的覆盖范围
    next = Math.max(next, jumps[i])
    // 需要做出决定，pre跳到最远的cur
    if (i === cur) {
      if (cur >= next) return -1 // 无法到达
      cur = next
      step++
    }
  }

  return step
}

export { maxJump }

if (require.main === module) {
  console.log(maxJump([1, 2, 0, 0, 0], 5))
  // console.log(maxJump([2, 3, 0, 1, 4]))
  // console.log(maxJump([1, 0, 0]))
}

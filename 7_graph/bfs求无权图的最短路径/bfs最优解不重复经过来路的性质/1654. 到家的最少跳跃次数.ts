import { ArrayDeque } from '../../../2_queue/Deque/ArrayDeque'

/**
 * 
 * @param bad 
 * @param rightJump  rightJump <= 2000
 * @param leftJump  leftJump <= 2000
 * @param target  家 0 <= target <= 2000
 * 它可以 往左 跳恰好 a 个位置。
   它可以 往右 跳恰好 b 个位置。
   它不能 连续 往左跳 2 次。
   它不能跳到任何 forbidden 数组中的位置。
   跳蚤可以往前跳 超过 它的家的位置，但是它 不能跳到负整数 的位置。
   @description 注意bfs范围剪枝
 */
function minimumJumps(bad: number[], rightJump: number, leftJump: number, target: number): number {
  const visited = new Set([...bad, 0])
  const upper = target + leftJump + rightJump // 上界 可以取大一点6000
  const queue = new ArrayDeque<[cur: number, step: number, canBackJump: boolean]>()
  queue.push([0, 0, false])

  while (queue.length) {
    const [cur, step, canBackJump] = queue.shift()!
    if (cur === target) return step

    if (canBackJump) {
      const next1 = cur - leftJump
      if (next1 >= 0 && !visited.has(next1)) {
        queue.push([next1, step + 1, false])
        visited.add(next1)
      }
    }

    const next2 = cur + rightJump
    // 剪枝 :出右边界后没必要继续走
    if (next2 <= upper && !visited.has(next2)) {
      queue.push([next2, step + 1, true])
      visited.add(next2)
    }
  }

  return -1
}

console.log(minimumJumps([14, 4, 18, 1, 15], 3, 15, 9))

export {}

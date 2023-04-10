import { ArrayDeque } from '../../../../2_queue/Deque/ArrayDeque'

function minimumJumps(bad: number[], rightJump: number, leftJump: number, target: number): number {
  const visited = new Set([...bad, 0])
  const upper = 6000 // 上界 可以取大一点6000
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

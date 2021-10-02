/**
 * 
 * @param forbidden 
 * @param a 
 * @param b 
 * @param x  家
 * 它可以 往左 跳恰好 a 个位置。
   它可以 往右 跳恰好 b 个位置。
   它不能 连续 往左跳 2 次。
   它不能跳到任何 forbidden 数组中的位置。
 */
function minimumJumps(forbidden: number[], a: number, b: number, x: number): number {
  const visited = new Set(forbidden)
  const limit = 2000 + a + b
  const queue: [number, number, boolean][] = [[0, 0, true]]

  while (queue.length) {
    const [current, jumps, canBackJump] = queue.shift()!
    if (current == x) return jumps

    // 剪枝1 :bfs最优解不走回头路
    if (visited.has(current)) continue
    visited.add(current)

    let nextJump: number
    if (canBackJump) {
      nextJump = current - b
      if (nextJump >= 0) queue.push([nextJump, jumps + 1, false])
    }

    nextJump = current + a
    // 剪枝2 :出界后没必要继续走
    if (nextJump <= limit) queue.push([nextJump, jumps + 1, true])
  }

  return -1
}

console.log(minimumJumps([14, 4, 18, 1, 15], 3, 15, 9))

export {}

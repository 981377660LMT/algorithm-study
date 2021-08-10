/**
 * @param {number[]} forbidden
 * @param {number} a
 * @param {number} b
 * @param {number} x
 * @return {number}
 * @description
 * 它可以 往前 跳恰好 a 个位置（即往右跳）。
   它可以 往后 跳恰好 b 个位置（即往左跳）。
   它不能 连续 往后跳 2 次。
   它不能跳到任何 forbidden 数组中的位置。
   1 <= a, b, forbidden[i] <= 2000
 */
var minimumJumps = function (forbidden, a, b, x) {
  const visited = new Set(forbidden)
  const limit = 2000 + a + b
  const queue = [[0, 0, true]]

  while (queue.length) {
    const [current, jumps, backJump] = queue.shift()
    if (current == x) return jumps

    if (visited.has(current)) continue

    visited.add(current)
    let nextJump
    if (backJump) {
      nextJump = current - b
      if (nextJump >= 0) queue.push([nextJump, jumps + 1, false])
    }

    nextJump = current + a
    if (nextJump <= limit) queue.push([nextJump, jumps + 1, true])
  }
  return -1
}

console.log(minimumJumps([14, 4, 18, 1, 15], 3, 15, 9))
